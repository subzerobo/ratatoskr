package logger

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func TestCreateLogger(t *testing.T) {
	cfg := Config{
		LogLevel:      "ERROR",
		GrayLogActive: true,
		GrayLogServer: "1.1.1.1",
		GrayLogStream: "my_stream",
	}
	logger := CreateLogger(cfg)
	if logger.GetLevel() != logrus.ErrorLevel {
		t.Fatalf("we got %v as logger level but expected %v", logger.GetLevel(), logrus.ErrorLevel)
	}
}
