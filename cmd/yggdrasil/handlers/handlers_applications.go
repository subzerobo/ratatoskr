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
// @Summary Gets application details
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

// HandleGetAndroidGroups godoc
// @Summary Android groups and channels
// @Description Gets a list of Android groups and child categories for the given Ratatoskr App
// @ID handle_get_android_groups
// @Tags AndroidGroups
// @Security BearerToken
// @Produce	json
// @Param app_uuid path string true "UUID of user-owned application"
// @Success 200 {object} rest.StandardResponse{data=[]applications.AndroidGroupModel} "Success Result"
// @Failure 500 {object} rest.StandardResponse
// @Router /v1/application/:app_uuid/android_groups [get]
func (h *YggdrasilHandler) HandleGetAndroidGroups(c *gin.Context) {
	aUUID := c.Param("app_uuid")
	res, err := h.applicationSvc.GetAndroidGroups(aUUID)
	if err != nil {
		kind, _ := errors.AsKindContext(err)
		c.JSON(kind.GetHttpStatus(), rest.GetFailMessageResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, rest.GetSuccessResponse(res))
}

// HandleCreateAndroidGroup godoc
// @Summary Creates android group
// @Description Creates a new Android group for the given Ratatoskr App
// @ID handle_create_android_group
// @Tags AndroidGroups
// @Security BearerToken
// @Accept	json
// @Produce	json
// @Param Application body AndroidGroupRequest true "Create Android Group"
// @Param app_uuid path string true "UUID of user-owned application"
// @Success 200 {object} rest.StandardResponse{} "Success Result"
// @Failure 500 {object} rest.StandardResponse
// @Router /v1/application/:app_uuid/android_groups [post]
func (h *YggdrasilHandler) HandleCreateAndroidGroup(c *gin.Context) {
	req := AndroidGroupRequest{}
	aUUID := c.Param("app_uuid")
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, rest.GetFailValidationResponse(err))
		return
	}

	claims := getClaims(c)

	err := h.applicationSvc.CreateAndroidGroup(claims.UserID, aUUID, req.Name)
	if err != nil {
		kind, _ := errors.AsKindContext(err)
		c.JSON(kind.GetHttpStatus(), rest.GetFailMessageResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, rest.GetSuccessResponse(nil))
}

// HandleUpdateAndroidGroup godoc
// @Summary Update android group
// @Description  Updates an Android group name for the given Ratatoskr App
// @ID handle_update_android_group
// @Tags AndroidGroups
// @Security BearerToken
// @Accept	json
// @Produce	json
// @Param Application body AndroidGroupRequest true "Update Android Group"
// @Param app_uuid path string true "UUID of user-owned application"
// @Param uuid path string true "UUID of android group"
// @Success 200 {object} rest.StandardResponse{} "Success Result"
// @Failure 500 {object} rest.StandardResponse
// @Router /v1/application/:app_uuid/android_groups/:uuid [put]
func (h *YggdrasilHandler) HandleUpdateAndroidGroup(c *gin.Context) {
	req := AndroidGroupRequest{}
	aUUID := c.Param("app_uuid")
	gUUID := c.Param("uuid")
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, rest.GetFailValidationResponse(err))
		return
	}
	claims := getClaims(c)

	err := h.applicationSvc.UpdateAndroidGroup(claims.UserID, aUUID, gUUID, req.Name)
	if err != nil {
		kind, _ := errors.AsKindContext(err)
		c.JSON(kind.GetHttpStatus(), rest.GetFailMessageResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, rest.GetSuccessResponse(nil))
}

// HandleDeleteAndroidGroup godoc
// @Summary Delete android group
// @Description  Deletes an Android group for the given Ratatoskr App with it's child channels
// @ID handle_delete_android_group
// @Tags AndroidGroups
// @Security BearerToken
// @Produce	json
// @Param app_uuid path string true "UUID of user-owned application"
// @Param uuid path string true "UUID of android group"
// @Success 200 {object} rest.StandardResponse{} "Success Result"
// @Failure 500 {object} rest.StandardResponse
// @Router /v1/application/:app_uuid/android_groups/:uuid [delete]
func (h *YggdrasilHandler) HandleDeleteAndroidGroup(c *gin.Context) {
	aUUID := c.Param("app_uuid")
	gUUID := c.Param("uuid")
	claims := getClaims(c)

	err := h.applicationSvc.DeleteAndroidGroup(claims.UserID, aUUID, gUUID)
	if err != nil {
		kind, _ := errors.AsKindContext(err)
		c.JSON(kind.GetHttpStatus(), rest.GetFailMessageResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, rest.GetSuccessResponse(nil))
}

// HandleCreateAndroidCategory godoc
// @Summary Create android category
// @Description  Create an Android category for the given Ratatoskr App / Android group
// @ID handle_create_android_category
// @Tags AndroidGroups
// @Security BearerToken
// @Accept	json
// @Produce	json
// @Param Application body AndroidCategoryRequest true "Update Android Group"
// @Param app_uuid path string true "UUID of user-owned application"
// @Param g_uuid path string true "UUID of android group"
// @Success 200 {object} rest.StandardResponse{} "Success Result"
// @Failure 500 {object} rest.StandardResponse
// @Router /v1/application/{app_uuid}/android_categories/{g_uuid} [post]
func (h *YggdrasilHandler) HandleCreateAndroidCategory(c *gin.Context) {
	req := AndroidCategoryRequest{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, rest.GetFailValidationResponse(err))
		return
	}

	aUUID := c.Param("app_uuid")
	gUUID := c.Param("g_uuid")
	claims := getClaims(c)

	err := h.applicationSvc.CreateAndroidCategory(claims.UserID, aUUID, gUUID, applications.AndroidGroupCategoryModel{
		CategoryName:        req.CategoryName,
		CategoryDescription: req.CategoryDescription,
		Priority:            req.Priority,
		Sound:               req.Sound,
		SoundName:           req.SoundName,
		Vibration:           req.Vibration,
		VibrationPattern:    req.VibrationPattern,
		Led:                 req.Led,
		LedColor:            req.LedColor,
		EnableBadge:         req.EnableBadge,
		LockScreen:          req.LockScreen,
	})
	if err != nil {
		kind, _ := errors.AsKindContext(err)
		c.JSON(kind.GetHttpStatus(), rest.GetFailMessageResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, rest.GetSuccessResponse(nil))
}

// HandleUpdateAndroidCategory godoc
// @Summary Updates android category
// @Description  Updates an Android category for the given Ratatoskr App / Android group
// @ID handle_update_android_category
// @Tags AndroidGroups
// @Security BearerToken
// @Produce	json
// @Param app_uuid path string true "UUID of user-owned application"
// @Param g_uuid path string true "UUID of android group"
// @Param c_uuid path string true "UUID of android category"
// @Success 200 {object} rest.StandardResponse{} "Success Result"
// @Failure 500 {object} rest.StandardResponse
// @Router /v1/application/{app_uuid}/android_categories/{g_uuid}/{c_uuid} [put]
func (h *YggdrasilHandler) HandleUpdateAndroidCategory(c *gin.Context) {
	req := AndroidCategoryRequest{}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, rest.GetFailValidationResponse(err))
		return
	}

	aUUID := c.Param("app_uuid")
	gUUID := c.Param("g_uuid")
	cUUID := c.Param("c_uuid")
	claims := getClaims(c)

	err := h.applicationSvc.UpdateAndroidCategory(claims.UserID, aUUID, gUUID, applications.AndroidGroupCategoryModel{
		CategoryUUID:        cUUID,
		CategoryName:        req.CategoryName,
		CategoryDescription: req.CategoryDescription,
		Priority:            req.Priority,
		Sound:               req.Sound,
		SoundName:           req.SoundName,
		Vibration:           req.Vibration,
		VibrationPattern:    req.VibrationPattern,
		Led:                 req.Led,
		LedColor:            req.LedColor,
		EnableBadge:         req.EnableBadge,
		LockScreen:          req.LockScreen,
	})
	if err != nil {
		kind, _ := errors.AsKindContext(err)
		c.JSON(kind.GetHttpStatus(), rest.GetFailMessageResponse(err.Error()))
		return
	}
	c.JSON(http.StatusOK, rest.GetSuccessResponse(nil))
}

// HandleDeleteAndroidCategory godoc
// @Summary Delete android category
// @Description  Deletes an Android category for the given Ratatoskr App / Android group
// @ID handle_delete_android_category
// @Tags AndroidGroups
// @Security BearerToken
// @Produce	json
// @Param app_uuid path string true "UUID of user-owned application"
// @Param g_uuid path string true "UUID of android group"
// @Param c_uuid path string true "UUID of android category"
// @Success 200 {object} rest.StandardResponse{} "Success Result"
// @Failure 500 {object} rest.StandardResponse
// @Router /v1/application/{app_uuid}/android_categories/{g_uuid}/{c_uuid} [delete]
func (h *YggdrasilHandler) HandleDeleteAndroidCategory(c *gin.Context) {
	aUUID := c.Param("app_uuid")
	gUUID := c.Param("g_uuid")
	cUUID := c.Param("c_uuid")
	claims := getClaims(c)

	err := h.applicationSvc.DeleteAndroidCategory(claims.UserID, aUUID, gUUID, cUUID)
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
	AuthKey string `json:"auth_key" example:"E4YfpiZLajkjtOO8BbOlNK5Skbs2Ez63EdrFBE7xdiruInuB7geHYlHpkr5rPHSy"`
}

type AndroidGroupRequest struct {
	Name string `json:"name" binding:"required"`
}

type AndroidCategoryRequest struct {
	CategoryName        string `json:"name" binding:"required" example:"Test Group"`
	CategoryDescription string `json:"description" example:"Test Description"`
	Priority            string `json:"priority" example:"3"`
	Sound               int    `json:"sound" example:"1"`
	SoundName           string `json:"sound_name" example:"resource_name"`
	Vibration           int    `json:"vibration" example:"1"`
	VibrationPattern    string `json:"vibration_pattern" example:"xxxx"`
	Led                 int    `json:"led" example:"1"`
	LedColor            string `json:"led_color" example:"#3300ccc"`
	EnableBadge         int    `json:"enable_badge" example:"0"`
	LockScreen          int    `json:"lock_screen" example:"1"`
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
