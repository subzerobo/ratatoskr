package authentication

import (
	"database/sql"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type AccountModel struct {
	ID                 uint
	UUID               string
	Email              string
	EncryptedPassword  string
	OAuthProvider      string
	OAuthUID           string
	Picture            string
	CompanyName        string
	IsSuperUser        bool
	LastLoginDate      sql.NullTime
	Active             bool
	Confirmed          bool
	ConfirmationToken  string
	ConfirmationSentAt sql.NullTime
	CreatedAt          time.Time `gorm:"default:current_timestamp"`
	UpdatedAt          time.Time `gorm:"default:current_timestamp"`
}

// GoogleResponse is model to unmarshal the google oauth login response
type GoogleResponse struct {
	ID        string `json:"sub"`
	Email     string `json:"email"`
	Picture   string `json:"picture"`
	Name      string `json:"name"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
}

type AppClaims struct {
	jwt.StandardClaims
	UserID  uint   `json:"user_id"`
	Version string   `json:"version"`
	Email   string   `json:"email"`
	Roles   []string `json:"roles"`
}
