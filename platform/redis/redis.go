package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

type Config struct {
	HOST     string `yaml:"HOST" envconfig:"REDIS_HOST"`
	PORT     int    `yaml:"PORT" envconfig:"REDIS_PORT"`
	PASSWORD string `yaml:"PASSWORD" envconfig:"REDIS_PASSWORD"`
}

type connectionString struct {
	connection string
	password   string
	domain     string
}

// Connections contains the functions to handle the redis platform
type Connections interface {
	Open() (*redis.Client, error)
}

// Initialize to init the redis platform with connection string and password
func Initialize(cfg Config, domain string) Connections {
	return &connectionString{
		connection: fmt.Sprintf("%s:%d",cfg.HOST, cfg.PORT),
		password:   cfg.PASSWORD,
		domain:     domain,
	}
}

// Open is to open a connection to redis server
func (cs *connectionString) Open() (*redis.Client, error) {
	logrus.WithFields(logrus.Fields{
		"platform": "redis",
		"domain":   cs.domain,
	}).Info("Connecting to Redis Server")
	client := redis.NewClient(&redis.Options{
		Addr:     cs.connection,
		Password: cs.password,
		DB:       0,
	})
	err := client.Ping(context.Background()).Err()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"connection": cs.connection,
			"password":   cs.password,
		}).Fatal(err)
		return nil, err
	}
	return client, nil
}
