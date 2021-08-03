package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/subzerobo/ratatoskr/internal/services/devices"
	"github.com/subzerobo/ratatoskr/pkg/errors"
	"github.com/subzerobo/ratatoskr/pkg/rest"
	"net/http"
)

// HandleAddDevice godoc
// @Summary Register a new device to one of your Ratatosk apps
// @Description Register a new device to one of your Ratatosk apps
// @ID handle_add_device
// @Tags Devices
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
		DeviceType:         req.DeviceType,
		Identifier:         req.Identifier,
		Language:           req.Language,
		Timezone:           req.Timezone,
		AppVersion:         req.AppVersion,
		DeviceVendor:       req.DeviceVendor,
		DeviceModel:        req.DeviceModel,
		DeviceOS:           req.DeviceOS,
		DeviceOSVersion:    req.DeviceOSVersion,
		ADID:               req.ADID,
		SDK:                req.SDK,
		SessionCount:       0,
		NotificationTypes:  req.NotificationTypes,
		Long:               req.Long,
		Lat:                req.Lat,
		Country:            req.Country,
		ExternalUserID:     req.ExternalUserID,
		ExternalUserIDHash: req.ExternalUserIdHash,
		Tags:               req.Tags,
	}, req.AppId)
	if err != nil {
		kind, _ := errors.AsKindContext(err)
		c.JSON(kind.GetHttpStatus(), rest.GetFailMessageResponse(err.Error()))
		return
	}
	
	c.JSON(http.StatusOK, rest.GetSuccessResponse(DeviceSuccessResponse{res.UUID}))
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

type DeviceSuccessResponse struct {
	UUID string `json:"uuid"`
}
