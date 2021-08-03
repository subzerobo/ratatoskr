package authentication

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"
	"github.com/subzerobo/ratatoskr/internal/constants"
	"github.com/subzerobo/ratatoskr/pkg/errors"
	"github.com/subzerobo/ratatoskr/pkg/mailer"
	"github.com/subzerobo/ratatoskr/pkg/utils"
	"io/ioutil"
	"time"
)

const DefaultStateValue = "1"

var (
	ErrAlreadyRegisteredEmail  = errors.New("email is already in use")
	ErrInvalidOauthProvider    = errors.New("oauth provider is not found")
	ErrAccountIsNotActive      = errors.New("account has been disabled by support")
	ErrUserIsNotConfirmed      = errors.New("user email is not confirmed yet")
	ErrInvalidLoginCredentials = errors.New("invalid login credentials")
)

type Service interface {
	Signup(email string, password string, company string) error
	Login(email string, password string) (*AccountModel, string, error)
	OAuthAuthenticate(provider string) (string, error)
	OAuthCallBack(state, code, provider string) (*AccountModel, string, error)
}

type service struct {
	BasePath   string
	Mailer     mailer.Service
	Config     Config
	repository Repository
	stateStore StateStore
}

func CreateService(r Repository, stateStore StateStore, config Config, mailer mailer.Service, basePath string) Service {
	return &service{
		repository: r,
		stateStore: stateStore,
		Config:     config,
		Mailer:     mailer,
		BasePath:   basePath,
	}
}

func (s service) Signup(email string, password string, company string) error {
	
	// Check Email is available
	facc, err := s.repository.GetAccountByEmail(email)
	if err != nil && !errors.HasKind(err, errors.NotFound) {
		return err
	}
	
	// Email is already in use!
	if facc != nil {
		return errors.WithKindCtx(ErrAlreadyRegisteredEmail, "", errors.Conflict, nil)
	}
	
	// We can create the account
	hashedPass, _ := utils.HashPassword(password)
	confirmationToken := utils.RandomString(64)
	model := AccountModel{
		Email:             email,
		EncryptedPassword: hashedPass,
		CompanyName:       company,
		IsSuperUser:       false,
		Active:            true,
		Confirmed:         false,
		ConfirmationToken: confirmationToken,
		ConfirmationSentAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}
	_, err = s.repository.CreateAccount(model)
	if err != nil {
		return errors.Wrap(err, "failed to create the account")
	}
	
	// Send Confirmation Email
	subject := "Ratatoskr Signup Confirmation"
	confirmUrl := fmt.Sprintf("%s/v1/auth/verify/%s", s.BasePath, confirmationToken)
	plainTextContent := fmt.Sprintf("Thanks for verifying you %s account!, Please visit this link to verify your registration: %s", email, confirmUrl)
	htmlContent := fmt.Sprintf(`Thanks for verifying you %[1]s account <br> Please visit this <a href="%[2]s">link</a> to verify your registration or copy this url and visit it using your browser <br> URL: %[2]s`, email, confirmUrl)
	err = s.Mailer.SendEmail(email, "", subject, plainTextContent, htmlContent)
	return err
}

func (s service) Login(email string, password string) (*AccountModel, string, error) {
	account, err := s.repository.GetAccountByEmail(email)
	if err != nil && !errors.HasKind(err, errors.NotFound) {
		return nil, "", err
	}
	
	if !utils.CheckPasswordHash(password, account.EncryptedPassword) {
		return nil, "", errors.WithKindCtx(ErrInvalidLoginCredentials, "", errors.Unauthorized, nil)
	}
	
	if !account.Active {
		return nil, "", errors.WithKindCtx(ErrAccountIsNotActive, "", errors.Forbidden, nil)
	}
	
	if !account.Confirmed {
		return nil, "", errors.WithKindCtx(ErrUserIsNotConfirmed, "", errors.Forbidden, nil)
	}
	
	var roles []constants.JWTRole
	switch {
	case account.IsSuperUser:
		roles = []constants.JWTRole{constants.Admin, constants.Free}
	default:
		roles = []constants.JWTRole{constants.Free}
	}
	
	token, err := s.generateJWTToken(account, roles)
	return account, token, err
}

func (s service) OAuthAuthenticate(provider string) (string, error) {
	providers := s.Config.GetOAuthProviders()
	providerConfig, ok := providers[provider]
	if !ok {
		return "", ErrInvalidOauthProvider
	}
	
	// Create State
	b := make([]byte, 32)
	rand.Read(b)
	state := base64.StdEncoding.EncodeToString(b)
	err := s.stateStore.SetState(state, DefaultStateValue)
	if err != nil {
		return "", errors.Wrap(err, "failed to store state")
	}
	
	u := providerConfig.AuthCodeURL(state)
	return u, nil
}

func (s service) OAuthCallBack(state, code, provider string) (*AccountModel, string, error) {
	providers := s.Config.GetOAuthProviders()
	providerConfig, ok := providers[provider]
	if !ok {
		return nil, "", ErrInvalidOauthProvider
	}
	
	var OpenID struct {
		ID           string
		ProfilePhoto string
		Email        string
		FirstName    string
		LastName     string
		FullName     string
	}
	
	// Check State
	state, err := s.stateStore.GetState(state)
	if err != nil || state != DefaultStateValue {
		return nil, "", errors.Wrap(err, "failed to validate the state")
	}
	
	token, err := providerConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, "", errors.WithKindCtx(err, "code exchange error", errors.Unauthorized, nil)
	}
	
	// TODO: Replace hard-code url with config : https://www.googleapis.com/oauth2/v3/userinfo
	client := providerConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to call googleapis.com")
	}
	
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, "", errors.Wrap(err, "failed to read googleapis.com response body")
	}
	var googleDTO GoogleResponse
	if err = json.Unmarshal(content, &googleDTO); err != nil {
		return nil, "", errors.Wrap(err, "failed to unmarshal googleapis.com response")
	}
	
	OpenID.ID = googleDTO.ID
	OpenID.FirstName = googleDTO.FirstName
	OpenID.LastName = googleDTO.LastName
	OpenID.Email = googleDTO.Email
	OpenID.FullName = googleDTO.Name
	OpenID.ProfilePhoto = googleDTO.Picture
	
	user, err := s.repository.GetAccountByEmail(OpenID.Email)
	if err != nil {
		if errors.HasKind(err, errors.NotFound) { // Register User
			// Generate Random Password
			hashedPassword, _ := utils.HashPassword(utils.RandomString(15))
			newUser := AccountModel{
				Email:             OpenID.Email,
				EncryptedPassword: hashedPassword,
				OAuthProvider:     provider,
				OAuthUID:          OpenID.ID,
				Picture:           OpenID.ProfilePhoto,
				CompanyName:       OpenID.FullName,
				Confirmed:         true,
				Active:            true,
			}
			
			user, err := s.repository.CreateAccount(newUser)
			if err != nil {
				return nil, "", errors.Wrap(err, "failed to retrieve the registered user")
			}
			var roles []constants.JWTRole
			switch {
			case user.IsSuperUser:
				roles = []constants.JWTRole{constants.Admin, constants.Free}
			default:
				roles = []constants.JWTRole{constants.Free}
			}
			token, err := s.generateJWTToken(user, roles)
			return user, token, err
		} else {
			return nil, "", errors.Wrapf(err, "Cannot retrieve user with email: %s", OpenID.Email)
		}
	} else {
		// Found User AppLogin the User
		if !user.Active {
			return nil, "", errors.WithKindCtx(ErrAccountIsNotActive, "", errors.Forbidden, nil)
		}
		var roles []constants.JWTRole
		switch {
		case user.IsSuperUser:
			roles = []constants.JWTRole{constants.Admin, constants.Free}
		default:
			roles = []constants.JWTRole{constants.Free}
		}
		token, err := s.generateJWTToken(user, roles)
		return user, token, err
	}
}

// generateJWTToken creates new JWTToken using user object
func (s service) generateJWTToken(user *AccountModel, Roles []constants.JWTRole) (string, error) {
	// Create JWT Token
	roles := make([]string, 0)
	for _, r := range Roles {
		roles = append(roles, r.String())
	}
	claims := AppClaims{
		StandardClaims: jwt.StandardClaims{
			Audience:  s.Config.JWT.Audience,
			ExpiresAt: time.Now().Add(time.Duration(24*s.Config.JWT.ExpireDays) * time.Hour).Unix(),
			Id:        uuid.NewV4().String(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    s.Config.JWT.Issuer,
			Subject:   user.UUID,
		},
		UserID:  user.ID,
		Version: s.Config.JWT.Version,
		Email:   user.Email,
		Roles:   roles,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(s.Config.JWT.Secret))
	return signedToken, nil
}
