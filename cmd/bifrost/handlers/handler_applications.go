package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/subzerobo/ratatoskr/pkg/errors"
	"github.com/subzerobo/ratatoskr/pkg/rest"
	"net/http"
)

// HandleAndroidParams godoc
// @Summary Apps Android Params
// @Description Gets the android params for one of your Ratatoskr apps
// @ID handle_android_params_device
// @Tags App,SDK
// @Produce	json
// @Param app_uuid path string true "UUID of user-owned application"
// @Success 200 {object} rest.StandardResponse "Success Result"
// @Failure 400 {object} rest.StandardResponse "Validation error"
// @Failure 500 {object} rest.StandardResponse
// @Router /apps/{app_uuid}/android_params [get]
func (h *BifrostHandler) HandleAndroidParams(c *gin.Context) {
	appUUID := c.Param("app_uuid")
	res, err := h.applicationSvc.GetAndroidParams(appUUID)
	if err != nil {
		kind, _ := errors.AsKindContext(err)
		c.JSON(kind.GetHttpStatus(), rest.GetFailMessageResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, rest.GetSuccessResponse(res))
}
