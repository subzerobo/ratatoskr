package handlers

import (
	"context"
	"fmt"
	sigar "github.com/cloudfoundry/gosigar"
	"github.com/gin-gonic/gin"
	"github.com/subzerobo/ratatoskr/internal/services/applications"
	"github.com/subzerobo/ratatoskr/internal/services/authentication"
	"github.com/subzerobo/ratatoskr/internal/services/devices"
	"github.com/subzerobo/ratatoskr/pkg/errors"
	"github.com/subzerobo/ratatoskr/pkg/logger"
	"github.com/subzerobo/ratatoskr/pkg/rest"
	"net/http"
	"runtime"
	"time"
)

// YggdrasilHandler is rest handler for Yggdrasil rest handler
type YggdrasilHandler struct {
	HealthCheckInfo struct {
		GitCommit     string
		BuildTime     string
		ContainerName string
		StartTime     time.Time
	}
	Logger         *logger.StandardLogger
	HTTPServer     *http.Server
	accountSvc     authentication.Service
	applicationSvc applications.Service
	deviceSvc      devices.Service
}

func CreateYggdrasilHandler(
	accountSvc authentication.Service,
	applicationSvc applications.Service,
	deviceSvd devices.Service,
	logger *logger.StandardLogger,
) *YggdrasilHandler {
	return &YggdrasilHandler{
		Logger:         logger,
		accountSvc:     accountSvc,
		applicationSvc: applicationSvc,
		deviceSvc:      deviceSvd,
	}
}

// Start starts the http server
func (h *YggdrasilHandler) Start(ctx context.Context, r *gin.Engine, defaultPort int) {
	const op = "http.rest.start"
	
	addr := fmt.Sprintf(":%v", defaultPort)
	
	h.HTTPServer = &http.Server{
		Addr:    addr,
		Handler: r,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
	
	h.Logger.Infof("[OK] Starting HTTP REST Server on %s ", addr)
	err := h.HTTPServer.ListenAndServe()
	if err != http.ErrServerClosed {
		h.Logger.Fatal(errors.WithMessage(err, op))
	}
	// Code Reach Here after HTTP Server Shutdown!
	h.Logger.Info("[OK] HTTP REST Server is shutting down!")
}

// Stop handles the http server in graceful shutdown
func (h *YggdrasilHandler) Stop() {
	const op = "http.rest.stop"
	
	// Create an 5s timeout context or waiting for app to shutdown after 5 seconds
	ctxTimeout, cancelTimeout := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelTimeout()
	
	h.HTTPServer.SetKeepAlivesEnabled(false)
	if err := h.HTTPServer.Shutdown(ctxTimeout); err != nil {
		h.Logger.Error(errors.WithMessage(err, op))
	}
	h.Logger.Info("HTTP REST Server graceful shutdown completed")
	
}

// HealthCheck godoc
// @Summary Pod/Container health check
// @Description Pod/Container health check
// @ID handle_health_check
// @Tags General
// @Accept	json
// @Produce	json
// @Success 200 {object} rest.StandardResponse{data=HealthCheckResponse} "Success Result"
// @Failure 400 {object} rest.StandardResponse
// @Failure 500 {object} rest.StandardResponse
// @Router / [get]
func (h *YggdrasilHandler) HealthCheck(c *gin.Context) {
	uptime := sigar.Uptime{}
	uptime.Get()
	avg := sigar.LoadAverage{}
	avg.Get()
	hcDTO := HealthCheckResponse{
		Status:       "ok",
		GitCommit:    h.HealthCheckInfo.GitCommit,
		BuildTime:    h.HealthCheckInfo.BuildTime,
		Container:    h.HealthCheckInfo.ContainerName,
		Version:      runtime.Version(),
		Uptime:       uptime.Format(),
		BinaryUptime: time.Since(h.HealthCheckInfo.StartTime).String(),
		LAOne:        avg.One,
		LAFive:       avg.Five,
		LAFifteen:    avg.Fifteen,
	}
	c.JSON(http.StatusOK, rest.GetSuccessResponse(hcDTO))
}
