package main

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/subzerobo/ratatoskr/cmd/bifrost/handlers"
	"github.com/subzerobo/ratatoskr/internal/services/applications"
	"github.com/subzerobo/ratatoskr/internal/services/devices"
	"github.com/subzerobo/ratatoskr/internal/storage/postgres"
	rs "github.com/subzerobo/ratatoskr/internal/storage/redis"
	"github.com/subzerobo/ratatoskr/pkg/logger"
	pg "github.com/subzerobo/ratatoskr/platform/postgres"
	"github.com/subzerobo/ratatoskr/platform/redis"
	"os"
	"sync"
)

type server struct {
	sync.WaitGroup
	Stan        stan.Conn
	Config      *Config
	Logger      *logger.StandardLogger
	RESTHandler *handlers.BifrostHandler
	NatsHandler *nats.Handler
}

// NewServer Create a new instance of server application
func NewServer(cfg *Config, stan stan.Conn) *server {
	return &server{
		Stan:   stan,
		Config: cfg,
	}
}

// Initialize is responsible for app initialization and wrapping required dependencies
func (s *server) Initialize(logger *logger.StandardLogger) error {
	// Initialize Database
	connection := pg.CreateConnection(s.Config.Database, "ratatoskr.io")
	
	// Initialize GORM
	gorm, err := connection.OpenGORM()
	if err != nil {
		return err
	}
	
	// Initialize Postgres Backed Repository
	repository, err := postgres.CreateRepository(gorm)
	if err != nil {
		return err
	}
	
	// Initialize Redis
	redisConnection := redis.Initialize(s.Config.Redis, fmt.Sprintf("Ratatoskr:%s", s.Config.Database.HOST))
	redisClient, err := redisConnection.Open()
	if err != nil {
		return err
	}
	rs.CreateRedisStore(redisClient)
	
	// Create Services
	applicationService := applications.CreateService(repository)
	deviceService := devices.CreateService(repository)
	
	// REST Handler
	restHandler := handlers.CreateBifrostHandler(applicationService, deviceService, logger)
	
	// Update GitCommit and BuildTime in handler
	restHandler.HealthCheckInfo.GitCommit = GitCommit
	restHandler.HealthCheckInfo.BuildTime = BuildTime
	restHandler.HealthCheckInfo.ContainerName = ContainerName
	restHandler.HealthCheckInfo.StartTime = StartTime
	
	s.RESTHandler = restHandler
	s.Logger = logger
	return nil
}

// Start starts the application in blocking mode
func (s *server) Start(ctx context.Context) {
	// Create Router for HTTP Server
	router := SetupRouter(s.RESTHandler, s.Config.Prometheus, s.Config.Authentication.JWT)
	
	// // Start Nats Worker
	// err := s.NatsHandler.Start(ctx)
	// if err != nil {
	// 	x, ok := err.(WithKindContextError)
	// 	if ok {
	// 		fields := logrus.Fields{}
	// 		for k, v := range x.Context() {
	// 			fields[k] = fmt.Sprintf("%v", v)
	// 		}
	// 		s.Logger.WithFields(fields).Error(err.Error())
	// 	} else {
	// 		s.Logger.Error(err.Error())
	// 	}
	//
	// }
	
	// Start REST Server in Blocking mode
	s.RESTHandler.Start(ctx, router, s.Config.Port)
}

// GracefulShutdown listen over the quitSignal to graceful shutdown the app
func (s *server) GracefulShutdown(quitSignal <-chan os.Signal, done chan<- bool) {
	// Wait for OS signals
	<-quitSignal
	
	// Kill the API Endpoints first
	s.RESTHandler.Stop()
	
	close(done)
}
