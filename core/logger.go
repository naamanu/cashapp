package core

import (
	"log"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func InitLogger(env Environment) {
	var config zap.Config

	if env == Development {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		config = zap.NewProductionConfig()
		config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	}

	var err error
	Log, err = config.Build()
	if err != nil {
		log.Printf("Failed to initialize logger: %v", err)
		os.Exit(1)
	}
	defer Log.Sync()
}
