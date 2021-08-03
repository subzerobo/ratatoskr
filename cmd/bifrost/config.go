package main

import (
	"github.com/kelseyhightower/envconfig"
	authentication2 "github.com/subzerobo/ratatoskr/internal/services/authentication"
	"github.com/subzerobo/ratatoskr/pkg/logger"
	"github.com/subzerobo/ratatoskr/platform/postgres"
	"github.com/subzerobo/ratatoskr/platform/redis"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

type Config struct {
	BasePath       string                 `yaml:"BASE_PATH" envconfig:"BASE_PATH"`
	Port           int                    `yaml:"PORT" envconfig:"PORT"`
	Logger         logger.Config          `yaml:"LOGGER"`
	Prometheus     PrometheusConfig       `yaml:"PROMETHEUS"`
	STAN           STAN                   `yaml:"STAN"`
	Database       postgres.Config        `yaml:"DB"`
	Redis          redis.Config           `yaml:"REDIS"`
	Mailer         MailerConfig           `yaml:"MAILER"`
	Authentication authentication2.Config `yaml:"AUTHENTICATION"`
}

type PrometheusConfig struct {
	Path     string `yaml:"PATH" envconfig:"PROMETHEUS_PATH"`
	UseAuth  bool   `yaml:"USE_AUTH" envconfig:"PROMETHEUS_USE_AUTH"`
	UserName string `yaml:"USERNAME" envconfig:"PROMETHEUS_USERNAME"`
	Password string `yaml:"PASSWORD" envconfig:"PROMETHEUS_PASSWORD"`
}

type STAN struct {
	ClusterID string   `yaml:"CLUSTER_ID" envconfig:"STAN_CLUSTER_ID"`
	NatsURLs  []string `yaml:"NATS_URLS" envconfig:"STAN_NATS_URLS"`
}

type MailerConfig struct {
	Token string `yaml:"SG_TOKEN" envconfig:"SG_MAILER_TOKEN"`
	Email string `yaml:"SG_EMAIL" envconfig:"SG_MAILER_EMAIL"`
	Name  string `yaml:"SG_NAME" envconfig:"SG_MAILER_NAME"`
}

func (s *STAN) GenerateURLs() string {
	return strings.Join(s.NatsURLs, ",")
}

// LoadConfig loads configs form provided yaml file or overrides it with env variables
func LoadConfig(filePath string) (*Config, error) {
	cfg := Config{}
	if filePath != "" {
		err := readFile(&cfg, filePath)
		if err != nil {
			return nil, err
		}
	}
	err := readEnv(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

func readFile(cfg *Config, filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		return err
	}
	return nil
}

func readEnv(cfg *Config) error {
	return envconfig.Process("", cfg)
}
