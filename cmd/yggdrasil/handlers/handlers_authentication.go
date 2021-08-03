package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/subzerobo/ratatoskr/pkg/errors"
	"github.com/subzerobo/ratatoskr/pkg/rest"
	"net/http"
)

// HandleSignup godoc
// @Summary Create new account using basic account information
// @Description Creates new account using email, password and company name
// @ID handle_signup
// @Tags Authentication
// @Accept	json
// @Produce	json
// @Param EventMessage body SignupRequest true "Signup Request Payload"
// @Success 200 {object} rest.StandardResponse "Success Result"
// @Failure 409 {object} rest.StandardResponse "Email is already in-use"
// @Failure 400 {object} rest.StandardResponse "Validation error"
// @Failure 500 {object} rest.StandardResponse
// @Router /v1/auth/signup [post]
func (h *YggdrasilHandler) HandleSignup(c *gin.Context) {
	req := SignupRequest{}
	
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, rest.GetFailValidationResponse(err))
		return
	}
	
	err := h.accountSvc.Signup(req.Email, req.Password, req.Company)
	if err != nil {
		kind, _ := errors.AsKindContext(err)
		c.JSON(kind.GetHttpStatus(), rest.GetFailMessageResponse(err.Error()))
		return
	}
	
	c.JSON(http.StatusOK, rest.GetSuccessResponse(nil))
}

// HandleLogin godoc
// @Summary Login user using email and password combination
// @Description Login user using email and password combination
// @ID handle_login
// @Tags Authentication
// @Accept	json
// @Produce	json
// @Param EventMessage body LoginRequest true "Login Request Payload"
// @Success 200 {object} rest.StandardResponse "success Result"
// @Failure 400 {object} rest.StandardResponse "validation error"
// @Failure 401 {object} rest.StandardResponse "invalid credentials"
// @Failure 403 {object} rest.StandardResponse "account is not active or email is not confirmed"
// @Failure 500 {object} rest.StandardResponse
// @Router /v1/auth/login [post]
func (h *YggdrasilHandler) HandleLogin(c *gin.Context) {
	req := LoginRequest{}
	
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, rest.GetFailValidationResponse(err))
		return
	}
	
	uInfo, token, err := h.accountSvc.Login(req.Email, req.Password)
	if err != nil {
		kind, _ := errors.AsKindContext(err)
		c.JSON(kind.GetHttpStatus(), rest.GetFailMessageResponse(err.Error()))
		return
	}
	
	c.JSON(http.StatusOK, rest.GetSuccessResponse(LoginResponse{
		Token: token,
		UserInfo: userInfoResponse{
			Company:      uInfo.CompanyName,
			Email:        uInfo.Email,
			ProfilePhoto: uInfo.Picture,
		},
	}))
}

// HandleOAuthLoginURL godoc
// @SummaryGenerate the OAuth URL for the specified provider
// @Description This endpoint generate an URL for the specified provider name on the url
// @ID handle_oauth_step1
// @Tags Authentication
// @Produce	json
// @Param provider path string true "OAuth provider name (google,..)"
// @Success 200 {object} rest.StandardResponse{data=OauthRedirectResponse} "Success Result"
// @Failure 500 {object} rest.StandardResponse{data=[]string} "Validation error(s)"
// @Router /v1/auth/oauth/{provider} [get]
func (h *YggdrasilHandler) HandleOAuthLoginURL(c *gin.Context) {
	provider := c.Param("provider")
	res, err := h.accountSvc.OAuthAuthenticate(provider)
	if err != nil {
		kind, _ := errors.AsKindContext(err)
		c.JSON(kind.GetHttpStatus(), rest.GetFailMessageResponse(err.Error()))
		return
	}
	
	c.JSON(http.StatusOK, rest.GetSuccessResponse(OauthRedirectResponse{URL: res}))
}

func (h *YggdrasilHandler) HandleOAuthCallback(c *gin.Context) {
	provider := c.Param("provider")
	state := c.Query("state")
	code := c.Query("code")
	uInfo, token, err := h.accountSvc.OAuthCallBack(state, code, provider)
	
	if err != nil {
		kind, _ := errors.AsKindContext(err)
		c.JSON(kind.GetHttpStatus(), rest.GetFailMessageResponse(err.Error()))
		return
	}
	
	c.JSON(http.StatusOK, rest.GetSuccessResponse(LoginResponse{
		Token: token,
		UserInfo: userInfoResponse{
			Company:      uInfo.CompanyName,
			Email:        uInfo.Email,
			ProfilePhoto: uInfo.Picture,
		},
	}))
	
}

type SignupRequest struct {
	Email    string `json:"email" binding:"required,email" example:"ali.kaviani@gmail.com"`
	Password string `json:"password" binding:"required,min=8" example:"testpassword"`
	Company  string `json:"company" binding:"required" example:"ACME"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"ali.kaviani@gmail.com"`
	Password string `json:"password" binding:"required,min=8" example:"testpassword"`
}

type LoginResponse struct {
	Token    string           `json:"token"`
	UserInfo userInfoResponse `json:"user_info,omitempty"`
}

type userInfoResponse struct {
	Company      string `json:"company"`
	Email        string `json:"email"`
	ProfilePhoto string `json:"profile_photo"`
}

type OauthRedirectResponse struct {
	URL string `json:"url"`
}
