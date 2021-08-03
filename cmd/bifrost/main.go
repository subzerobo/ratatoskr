package main

import (
	"context"
	"crypto/md5"
	"flag"
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"github.com/subzerobo/ratatoskr/pkg/logger"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	GitCommit     string = "Development"
	BuildTime     string = time.Now().Format(time.RFC1123Z)
	ContainerName string
	StartTime     time.Time = time.Now()
	
	// CommitMetric holds the version information
	commitMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: "rataroskr",
			Subsystem: "bifrost",
			Name:      "version",
			Help:      "version of application",
		},
		[]string{"commit", "build_time"},
	)
)

// @title Ratatoskr(Bifrost)
// @version 1.0
// @description Ratatoskr SDK Public REST APIs(Bifrost)
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email ali.kaviani@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:6060
// @BasePath /
func main() {
	// Set Main Operation
	const op = "ratatoskr.bifrost.api"
	
	// Set Binary Start Time
	StartTime = time.Now()
	
	// Default Config file based on the environment variable
	defaultConfigFile := "configs/bifrost/config-local.yaml"
	if env := os.Getenv("APP_MODE"); env != "" {
		defaultConfigFile = fmt.Sprintf("configs/bifrost/config-%s.yaml", env)
	}
	
	// Load Master Config File
	var configFile string
	flag.StringVar(&configFile, "c", defaultConfigFile, "The environment configuration file of application")
	flag.StringVar(&configFile, "config", defaultConfigFile, "The environment configuration file of application")
	flag.Usage = usage
	flag.Parse()
	
	// Print Start Ascii Art
	printAsciiArt()
	
	// Set Commit Metrics
	gitMetrics()
	
	// Setting up the main context
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	
	// Loading the config file
	cfg, err := LoadConfig(configFile)
	if err != nil {
		logrus.Fatal(errors.Wrapf(err, "failed to load config: %s", op))
	}
	
	// Setup Logger
	logger := logger.CreateLogger(cfg.Logger)
	logger.Info("[OK] Logger Configured")
	
	// Show the loaded config file
	logger.Infof("[OK] Loaded config file: %s", configFile)
	
	// Get OS Container Name
	hostname, err := os.Hostname()
	if err != nil {
		logger.Fatal(errors.WithMessage(err, op))
	}
	ContainerName = hostname
	logger.Infof("[OK] Hostname acquired :%s", hostname)
	
	// Commit, BuildTime
	logger.Infof("[OK] Commit Number:%s, Build Time: %s", GitCommit, BuildTime)
	
	// Connect to NATS Streaming Server
	logger.Infof("[...] Trying to connect to nats urls: %s", cfg.STAN.GenerateURLs())
	opts := []nats.Option{nats.Name("App Mode, Nats Streaming Connection")}
	// Connect to Nats
	nc, err := nats.Connect(cfg.STAN.GenerateURLs(), opts...)
	if err != nil {
		logger.Fatal(errors.Wrap(err, "failed to connect to nats server"))
	}
	// Connect to Stan
	clientId := fmt.Sprintf("ratatoskr_bifrost_%s", fmt.Sprintf("%x", md5.Sum([]byte(hostname))))
	logger.Infof("[OK] Nats Client ID: %s", clientId)
	sc, err := stan.Connect(cfg.STAN.ClusterID, clientId, stan.NatsConn(nc), stan.SetConnectionLostHandler(func(conn stan.Conn, err error) {
		logger.Fatalf("Nats Connection Lost, reason: %v", err)
	}))
	if err != nil {
		logger.Fatal(errors.WithMessage(err, op))
	}
	defer sc.Close()
	
	// Create New Server
	server := NewServer(cfg, sc)
	
	// Initialize the Server Dependencies
	err = server.Initialize(logger)
	if err != nil {
		logger.Fatal(errors.Wrapf(err, "failed to initialize server: %s", op))
	}
	
	done := make(chan bool, 1)
	quiteSignal := make(chan os.Signal, 1)
	signal.Notify(quiteSignal, syscall.SIGINT, syscall.SIGTERM)
	
	// Graceful shutdown goroutine
	go server.GracefulShutdown(quiteSignal, done)
	
	// Start server in blocking mode
	server.Start(ctx)
	
	// Wait for HTTP Server to be killed gracefully !
	<-done
	
	// Killing other background jobs !
	cancel()
	logger.Info("Waiting for background jobs to finish their works...")
	
	// Wait for all other background jobs to finish their works
	server.Wait()
	
	logger.Info("Ratatoskr App Shutdown successfully, see you next time ;-)")
}

func usage() {
	usageStr := `
Usage: ratatoskr [options]
Options:
	-c,  --config   <config file name>   Path of yaml configuration file
`
	fmt.Printf("%s\n", usageStr)
	os.Exit(0)
}

func printAsciiArt() {

	fmt.Println(aurora.Blue(`
╔══════════════════════════════════════════════════════════════════════════════════════════════════════╗
║               |\=.                                                                                   ║
║               /  6',   ██████╗  █████╗ ████████╗ █████╗ ████████╗ ██████╗ ███████╗██╗  ██╗██████╗    ║
║       .--.    \  .-'   ██╔══██╗██╔══██╗╚══██╔══╝██╔══██╗╚══██╔══╝██╔═══██╗██╔════╝██║ ██╔╝██╔══██╗   ║
║      /_   \   /  (_()  ██████╔╝███████║   ██║   ███████║   ██║   ██║   ██║███████╗█████╔╝ ██████╔╝   ║
║        )   | / ';--'   ██╔══██╗██╔══██║   ██║   ██╔══██║   ██║   ██║   ██║╚════██║██╔═██╗ ██╔══██╗   ║
║       /   / /   (      ██║  ██║██║  ██║   ██║   ██║  ██║   ██║   ╚██████╔╝███████║██║  ██╗██║  ██║   ║
║     (    '"    _)_     ╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝   ╚═╝  ╚═╝   ╚═╝    ╚═════╝ ╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝   ║
║      '-==-'""""""                                                                                    ║
║                                      **  Bifrost - SDK's Public API Project **                       ║
╚══════════════════════════════════════════════════════════════════════════════════════════════════════╝
`))
}

func gitMetrics() {
	prometheus.MustRegister(commitMetric)
	commitMetric.WithLabelValues(GitCommit, BuildTime)
}
