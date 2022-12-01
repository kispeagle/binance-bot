package logger

import (
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	godotenv.Load(".env")

}

var (
	logger *zap.SugaredLogger
	dblog  *mongo.Collection
)

func GetLogger() *zap.SugaredLogger {
	if logger == nil {
		conf := zap.NewProductionConfig()

		if os.Getenv("level_logger") == "debug" {
			conf.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		}

		conf.EncoderConfig.TimeKey = "timestamp"
		conf.EncoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
		conf.OutputPaths = append(conf.OutputPaths, "logger/logfiles/log")

		l, _ := conf.Build()
		logger = l.Sugar()
	}
	return logger
}
