package logger

import (
	"log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func InitLogger(pEnv string) {
	var lConfig zap.Config
	if pEnv == "production" {
		lConfig = zap.NewProductionConfig()
		lConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	} else {
		lConfig = zap.NewDevelopmentConfig()
		lConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	lLogger, lErr := lConfig.Build()
	if lErr != nil {
		log.Fatalf("Failed to initialize zap logger: %v", lErr)
	}

	Log = lLogger
	zap.ReplaceGlobals(lLogger)
}
