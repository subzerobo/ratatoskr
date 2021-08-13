package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/subzerobo/ratatoskr/internal/services/devices"
	"github.com/subzerobo/ratatoskr/pkg/errors"
	"github.com/subzerobo/ratatoskr/pkg/rest"
	"net/http"
	"time"
)

// HandleAddDevice godoc
// @Summary Register a new device to one of your Ratatoskr apps
// @Description Register a new device to one of your Ratatoskr apps
// @ID handle_add_device
// @Tags Devices,SDK
// @Accept	json
// @Produce	json
// @Param Device body DeviceRequest true "Create Device Request"
// @Success 200 {object} rest.StandardResponse "Success Result"
// @Failure 400 {object} rest.StandardResponse "Validation error"
// @Failure 500 {object} rest.StandardResponse
// @Router /v1/devices [post]
func (h *BifrostHandler) HandleAddDevice(c *gin.Context) {
	req := DeviceRequest{}

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, rest.GetFailValidationResponse(err))
		return
	}

	res, err := h.deviceSvc.Upsert(devices.DeviceModel{
		DeviceType:         &req.DeviceType,
		Identifier:         &req.Identifier,
		Language:           &req.Language,
		Timezone:           req.Timezone,
		AppVersion:         &req.AppVersion,
		DeviceVendor:       &req.DeviceVendor,
		DeviceModel:        &req.DeviceModel,
		DeviceOS:           &req.DeviceOS,
		DeviceOSVersion:    &req.DeviceOSVersion,
		ADID:               &req.ADID,
		SDK:                &req.SDK,
		NotificationTypes:  &req.NotificationTypes,
		Long:               &req.Long,
		Lat:                &req.Lat,
		Country:            &req.Country,
		ExternalUserID:     &req.ExternalUserID,
		ExternalUserIDHash: &req.ExternalUserIdHash,
		Tags:               req.Tags,
	}, req.AppId)
	if err != nil {
		kind, _ := errors.AsKindContext(err)
		c.JSON(kind.GetHttpStatus(), rest.GetFailMessageResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, rest.GetSuccessResponse(DeviceSuccessResponse{*res.UUID}))
}

// HandleViewDevice godoc
// @Summary View the details of an existing device in one of your Ratatoskr apps
// @Description View the details of an existing device in one of your Ratatoskr apps
// @ID handle_view_device
// @Tags Devices,SDK
// @Produce	json
// @Param uuid path string true "UUID of application"
// @Param app_uuid path string true "UUID of device"
// @Success 200 {object} rest.StandardResponse{data=DeviceViewResponse} "Success Result"
// @Failure 500 {object} rest.StandardResponse
// @Router /v1/devices/{uuid}/{app_uuid} [get]
func (h *BifrostHandler) HandleViewDevice(c *gin.Context) {
	UUID := c.Param("uuid")
	appUUID := c.Param("app_uuid")

	res, err := h.deviceSvc.Get(UUID, appUUID)
	if err != nil {
		kind, _ := errors.AsKindContext(err)
		c.JSON(kind.GetHttpStatus(), rest.GetFailMessageResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, rest.GetSuccessResponse(toSingle(res)))
}

// HandleViewDevices godoc
// @Summary View the details of multiple devices in one of your Ratatoskr apps
// @Description View the details of multiple devices in one of your Ratatoskr apps
// @ID handle_view_devices
// @Tags Devices,SDK
// @Produce	json
// @Param app_uuid query string true "UUID of device"
// @Param limit query int false "How many devices to return. Max is 300. Default is 300"
// @Param last_device_id query int false "Previous max record id.Default is 0. Results are sorted by id;"
// @Security APIKey
// @Success 200 {object} rest.StandardResponse{data=[]DeviceViewResponse} "Success Result"
// @Failure 500 {object} rest.StandardResponse
// @Router /v1/devices [get]
func (h *BifrostHandler) HandleViewDevices(c *gin.Context) {
	appUUID := c.Query("app_uuid")

	authToken := c.GetHeader("Authorization")
	err := h.applicationSvc.CheckApplicationToken(authToken, appUUID)
	if err != nil {
		kind, _ := errors.AsKindContext(err)
		c.JSON(kind.GetHttpStatus(), rest.GetFailMessageResponse(err.Error()))
		return
	}

	mPaging, err := getMorePagination(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, rest.GetFailMessageResponse(err.Error()))
	}

	res, err := h.deviceSvc.GetList(appUUID, mPaging)
	if err != nil {
		kind, _ := errors.AsKindContext(err)
		c.JSON(kind.GetHttpStatus(), rest.GetFailMessageResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, rest.GetSuccessResponse(toList(res)))
}

// HandleEditDevice godoc
// @Summary Update an existing device in one of your Ratatoskr apps
// @Description Update an existing device in one of your Ratatoskr apps
// @ID handle_edit_device
// @Tags Devices,SDK
// @Accept	json
// @Produce	json
// @Param UUID path string true "Device Unique Identifier"
// @Param Device body DeviceRequest true "Create Device Request"
// @Success 200 {object} rest.StandardResponse "Success Result"
// @Failure 400 {object} rest.StandardResponse "Validation error"
// @Failure 500 {object} rest.StandardResponse
// @Router /v1/devices/{UUID} [put]
func (h *BifrostHandler) HandleEditDevice(c *gin.Context) {
	req := DeviceEditRequest{}

	uuid := c.Param("uuid")

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, rest.GetFailValidationResponse(err))
		return
	}

	_, err := h.deviceSvc.Update(devices.DeviceModel{
		UUID:              &uuid,
		DeviceType:        req.DeviceType,
		Identifier:        req.Identifier,
		Language:          req.Language,
		Timezone:          req.Timezone,
		AppVersion:        req.AppVersion,
		DeviceVendor:      req.DeviceVendor,
		DeviceModel:       req.DeviceModel,
		DeviceOS:          req.DeviceOS,
		DeviceOSVersion:   req.DeviceOSVersion,
		ADID:              req.ADID,
		SDK:               req.SDK,
		SessionCount:      req.SessionCount,
		NotificationTypes: req.NotificationTypes,
		Long:              req.Long,
		Lat:               req.Lat,
		Country:           req.Country,
		ExternalUserID:    req.ExternalUserID,
		Tags:              req.Tags,
		BadgeCount:        req.BadgeCount,
		AmountSpent:       req.AmountSpent,
	})
	if err != nil {
		kind, _ := errors.AsKindContext(err)
		c.JSON(kind.GetHttpStatus(), rest.GetFailMessageResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, rest.GetSuccessResponse(nil))
}

// HandleEditUserTags godoc
// @Summary Update an existing device's tags in one of your Ratatoskr apps using the External User ID.
// @Description Update an existing device's tags in one of your Ratatoskr apps using the External User ID.
// @ID handle_edit_user_tags
// @Tags Devices,SDK
// @Accept	json
// @Produce	json
// @Param APP_UUID path string true "App Unique Identifier UUID"
// @Param EXTERNAL_USER_ID path string true "External User ID"
// @Param UserTags body UserTagRequest true "List of User Tags"
// @Success 200 {object} rest.StandardResponse "Success Result"
// @Failure 400 {object} rest.StandardResponse "Validation error"
// @Failure 500 {object} rest.StandardResponse
// @Router /v1/applications/{APP_UUID}/users/{EXTERNAL_USER_ID} [put]
func (h *BifrostHandler) HandleEditUserTags(c *gin.Context) {
	req := UserTagRequest{}

	appUUID := c.Param("app_uuid")
	externalUserID := c.Param("external_user_id")

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, rest.GetFailValidationResponse(err))
		return
	}

	err := h.deviceSvc.UpdateUserTags(appUUID, externalUserID, req.Tags)
	if err != nil {
		kind, _ := errors.AsKindContext(err)
		c.JSON(kind.GetHttpStatus(), rest.GetFailMessageResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, rest.GetSuccessResponse(nil))
}

type DeviceRequest struct {
	AppId              string            `json:"app_id" binding:"required" example:"407f8f90-d83b-4ad5-912c-556a27c8f249"`
	DeviceType         string            `json:"device_type" binding:"required" example:"android | ios | web"`
	Identifier         string            `json:"identifier" binding:"required" example:"APA91bHbYHk7aq-Uam_2pyJ2qbZvqllyyh2wjfPRaw5gLEX2SUlQBRvOc6sck1sa7H7nGeLNlDco8lXj83HWWwzV..."`
	Language           string            `json:"language" example:"fa"`
	Timezone           int               `json:"timezone" example:"12600"`
	AppVersion         string            `json:"app_version" example:"2.1.1"`
	DeviceVendor       string            `json:"device_vendor" example:"Samsung"`
	DeviceModel        string            `json:"device_model" example:"SM-989F"`
	DeviceOS           string            `json:"device_os" example:"Android"`
	DeviceOSVersion    string            `json:"device_os_version" example:"8.0"`
	ADID               string            `json:"adid" example:"dbdf14cc-a5e7-445f-a972-2112ab335b14"`
	SDK                string            `json:"sdk" example:"1.0"`
	SessionCount       int               `json:"session_count" example:"1"`
	NotificationTypes  int               `json:"notification_types" example:"1"`
	Long               float32           `json:"long" example:"35.123456"`
	Lat                float32           `json:"lat" example:"54.123456"`
	Country            string            `json:"country" example:"IR"`
	ExternalUserID     string            `json:"external_user_id" example:"u-12"`
	ExternalUserIdHash string            `json:"external_user_id_hash" example:"xxxxxxxx" `
	Tags               map[string]string `json:"tags"`
}

type DeviceEditRequest struct {
	DeviceType        *string           `json:"device_type" example:"android | ios | web"`
	Identifier        *string           `json:"identifier" example:"APA91bHbYHk7aq-Uam_2pyJ2qbZvqllyyh2wjfPRaw5gLEX2SUlQBRvOc6sck1sa7H7nGeLNlDco8lXj83HWWwzV..."`
	Language          *string           `json:"language" example:"fa"`
	Timezone          int               `json:"timezone" example:"12600"`
	AppVersion        *string           `json:"app_version" example:"2.1.1"`
	DeviceVendor      *string           `json:"device_vendor" example:"Samsung"`
	DeviceModel       *string           `json:"device_model" example:"SM-989F"`
	DeviceOS          *string           `json:"device_os" example:"Android"`
	DeviceOSVersion   *string           `json:"device_os_version" example:"8.0"`
	ADID              *string           `json:"adid" example:"dbdf14cc-a5e7-445f-a972-2112ab335b14"`
	SDK               *string           `json:"sdk" example:"1.0"`
	SessionCount      *int              `json:"session_count" example:"1"`
	NotificationTypes *int              `json:"notification_types" example:"1"`
	Long              *float32          `json:"long" example:"35.123456"`
	Lat               *float32          `json:"lat" example:"54.123456"`
	Country           *string           `json:"country" example:"IR"`
	ExternalUserID    *string           `json:"external_user_id" example:"u-12"`
	Tags              map[string]string `json:"tags"`
	BadgeCount        *int              `json:"badge_count" example:"1"`
	AmountSpent       *float32          `json:"amount_spent" example:"29.99"`
}

type DeviceViewResponse struct {
	Identifier      string            `json:"identifier" example:"APA91bHbYHk7aq-Uam_2pyJ2qbZvqllyyh2wjfPRaw5gLEX2SUlQBRvOc6sck1sa7H7nGeLNlDco8lXj83HWWwzV..."`
	DeviceType      string            `json:"device_type" example:"android | ios | web"`
	Language        string            `json:"language" example:"fa"`
	Timezone        int               `json:"timezone" example:"12600"`
	AppVersion      string            `json:"app_version" example:"2.1.1"`
	DeviceVendor    string            `json:"device_vendor" example:"Samsung"`
	DeviceModel     string            `json:"device_model" example:"SM-989F"`
	DeviceOS        string            `json:"device_os" example:"Android"`
	DeviceOSVersion string            `json:"device_os_version" example:"8.0"`
	ADID            string            `json:"adid" example:"dbdf14cc-a5e7-445f-a972-2112ab335b14"`
	Tags            map[string]string `json:"tags"`
	CreatedAt       time.Time         `json:"created_at"`
	LastActiveAt    time.Time         `json:"last_active_at"`
	ExternalUserID  string            `json:"external_user_id" example:"u-12"`
	BadgeCount      int               `json:"badge_count" example:"1"`
}

type DeviceSuccessResponse struct {
	UUID string `json:"uuid"`
}

type UserTagRequest struct {
	Tags map[string]string `json:"tags" binding:"required"`
}

func toSingle(item *devices.DeviceModel) *DeviceViewResponse {
	return &DeviceViewResponse{
		Identifier:      *item.Identifier,
		DeviceType:      *item.DeviceType,
		Language:        *item.Language,
		Timezone:        item.Timezone,
		AppVersion:      *item.AppVersion,
		DeviceVendor:    *item.DeviceVendor,
		DeviceModel:     *item.DeviceModel,
		DeviceOS:        *item.DeviceOS,
		DeviceOSVersion: *item.DeviceOSVersion,
		ADID:            *item.ADID,
		Tags:            item.Tags,
		CreatedAt:       item.CreatedAt,
		LastActiveAt:    item.UpdatedAt,
		ExternalUserID:  *item.ExternalUserID,
		BadgeCount:      *item.BadgeCount,
	}
}

func toList(list []*devices.DeviceModel) []*DeviceViewResponse {
	results := make([]*DeviceViewResponse, 0)
	for _, res := range list {
		results = append(results, toSingle(res))
	}
	return results
}
