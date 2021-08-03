package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/subzerobo/ratatoskr/internal/services/applications"
	"github.com/subzerobo/ratatoskr/pkg/errors"
	"github.com/subzerobo/ratatoskr/pkg/rest"
	"net/http"
	"strconv"
)

// HandleCreateApplication godoc
// @Summary Create new application for logged-in account
// @Description Create new application for logged-in account
// @ID handle_create_application
// @Tags Applications
// @Security BearerToken
// @Accept	json
// @Produce	json
// @Param Application body ApplicationRequest true "Create Application Request"
// @Success 200 {object} rest.StandardResponse{data=ApplicationResponse} "Success Result"
// @Failure 400 {object} rest.StandardResponse "Validation error"
// @Failure 500 {object} rest.StandardResponse
// @Router /v1/applications [post]
func (h *YggdrasilHandler) HandleCreateApplication(c *gin.Context) {
	req := ApplicationRequest{}
	
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, rest.GetFailValidationResponse(err))
		return
	}
	
	claims := getClaims(c)
	
	res, err := h.applicationSvc.Create(applications.ApplicationModel{
		Name:         req.Name,
		FCMSenderID:  req.FMCSenderID,
		FCMAdminJSON: req.FCMAdminJson,
		URL:          req.URL,
		AccountID:    claims.UserID,
	})
	if err != nil {
		kind, _ := errors.AsKindContext(err)
		c.JSON(kind.GetHttpStatus(), rest.GetFailMessageResponse(err.Error()))
		return
	}
	
	c.JSON(http.StatusOK, rest.GetSuccessResponse(toSingle(res)))
}

// HandleListMyApplication godoc
// @Summary Lists all of account registered applications
// @Description Lists all of account registered applications
// @ID handle_list_my_applications
// @Tags Applications
// @Security BearerToken
// @Accept	json
// @Produce	json
// @Success 200 {object} rest.StandardResponse{data=[]ApplicationResponse} "Success Result"
// @Failure 500 {object} rest.StandardResponse
// @Router /v1/applications [get]
func (h *YggdrasilHandler) HandleListMyApplication(c *gin.Context) {
	claims := getClaims(c)
	
	res, err := h.applicationSvc.List(claims.UserID)
	if err != nil {
		kind, _ := errors.AsKindContext(err)
		c.JSON(kind.GetHttpStatus(), rest.GetFailMessageResponse(err.Error()))
		return
	}
	
	c.JSON(http.StatusOK, rest.GetSuccessResponse(toList(res)))
}

// HandleApplicationDetail godoc
// @Summary Gets the general info about specific application uuid
// @Description Gets the general info about specific application uuid
// @ID handle_get_application_detail
// @Tags Applications
// @Security BearerToken
// @Produce	json
// @Param uuid path string true "UUID of user-owned application"
// @Success 200 {object} rest.StandardResponse{data=ApplicationResponse} "Success Result"
// @Failure 404 {object} rest.StandardResponse
// @Failure 500 {object} rest.StandardResponse
// @Router /v1/applications/{uuid} [get]
func (h *YggdrasilHandler) HandleApplicationDetail(c *gin.Context) {
	uuid := c.Param("uuid")
	claims := getClaims(c)
	
	res, err := h.applicationSvc.Details(claims.UserID, uuid)
	if err != nil {
		kind, _ := errors.AsKindContext(err)
		c.JSON(kind.GetHttpStatus(), rest.GetFailMessageResponse(err.Error()))
		return
	}
	
	c.JSON(http.StatusOK, rest.GetSuccessResponse(toSingle(res)))
}

// HandleResetAuthToken godoc
// @Summary Updates the Auth token with new one
// @Description Updates the Auth token with new one
// @ID handle_reset_auth_token
// @Tags Applications
// @Security BearerToken
// @Produce	json
// @Param uuid path string true "UUID of user-owned application"
// @Success 200 {object} rest.StandardResponse{data=AuthKeyResponse} "Success Result"
// @Failure 500 {object} rest.StandardResponse
// @Router /v1/applications/{uuid} [patch]
func (h *YggdrasilHandler) HandleResetAuthToken(c *gin.Context) {
	uuid := c.Param("uuid")
	claims := getClaims(c)
	
	res, err := h.applicationSvc.UpdateAuthKey(claims.UserID, uuid)
	if err != nil {
		kind, _ := errors.AsKindContext(err)
		c.JSON(kind.GetHttpStatus(), rest.GetFailMessageResponse(err.Error()))
		return
	}
	
	c.JSON(http.StatusOK, rest.GetSuccessResponse(AuthKeyResponse{res}))
}

// HandleUpdateIdentityVerification godoc
// @Summary Updates Identity verification
// @Description Updates Identity verification
// @ID handle_update_identity_verification
// @Tags Applications
// @Security BearerToken
// @Accept	json
// @Produce	json
// @Param status path string true "UUID of user-owned application"
// @Param uuid path string true "Status of Identity verification"
// @Success 200 {object} rest.StandardResponse{} "Success Result"
// @Failure 500 {object} rest.StandardResponse
// @Router /v1/applications/{uuid}/{status} [put]
func (h *YggdrasilHandler) HandleUpdateIdentityVerification(c *gin.Context) {
	uuid := c.Param("uuid")
	status := c.Param("status")
	idVerificationStatus, _ := strconv.ParseBool(status)
	claims := getClaims(c)
	
	err := h.applicationSvc.UpdateIdentityVerification(claims.UserID, uuid, idVerificationStatus)
	if err != nil {
		kind, _ := errors.AsKindContext(err)
		c.JSON(kind.GetHttpStatus(), rest.GetFailMessageResponse(err.Error()))
		return
	}
	
	c.JSON(http.StatusOK, rest.GetSuccessResponse(nil))
}

type ApplicationRequest struct {
	Name         string `json:"name" binding:"required" example:"My Fancy Application"`
	FMCSenderID  string `json:"fmc_sender_id" binding:"required" example:"123456789"`
	FCMAdminJson string `json:"fcm_admin_json" binding:"required" example:"{....}"`
	URL          string `json:"url" binding:"required" example:"https://myfancywebsite.com"`
}

type ApplicationResponse struct {
	ID           uint   `json:"id" example:"1"`
	UUID         string `json:"uuid" example:"2550a565-98b4-47ce-9529-ab5c0da51556"`
	Name         string `json:"name"  example:"My Fancy Application"`
	FMCSenderID  string `json:"fmc_sender_id" example:"123456789"`
	FCMAdminJson string `json:"fcm_admin_json" example:"{....}"`
	AuthKey      string `json:"auth_key" example:"E4YfpiZLajkjtOO8BbOlNK5Skbs2Ez63EdrFBE7xdiruInuB7geHYlHpkr5rPHSy"`
	URL          string `json:"url" example:"https://myfancywebsite.com"`
}

type AuthKeyResponse struct {
	AuthKey      string `json:"auth_key" example:"E4YfpiZLajkjtOO8BbOlNK5Skbs2Ez63EdrFBE7xdiruInuB7geHYlHpkr5rPHSy"`
}

func toSingle(item *applications.ApplicationModel) *ApplicationResponse {
	return &ApplicationResponse{
		ID:           item.ID,
		UUID:         item.UUID,
		Name:         item.Name,
		FMCSenderID:  item.FCMSenderID,
		FCMAdminJson: item.FCMAdminJSON,
		URL:          item.URL,
	}
}

func toList(list []*applications.ApplicationModel) []*ApplicationResponse {
	results := make([]*ApplicationResponse, 0)
	for _, res := range list {
		results = append(results, toSingle(res))
	}
	return results
}
