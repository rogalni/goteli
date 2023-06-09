package goteli

import (
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func setupLogger(logLevel string, isJsonLogging bool, serviceName string) func() error {
	lvl, err := zap.ParseAtomicLevel(logLevel)
	if err != nil {
		lvl = zap.NewAtomicLevelAt(zap.InfoLevel)
	}
	var conf zap.Config
	if isJsonLogging {
		conf = zap.NewProductionConfig()
	} else {
		conf = zap.NewDevelopmentConfig()
	}

	conf.Level = lvl
	conf.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	conf.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder

	zl, _ := conf.Build(
		zap.Fields(zap.String("service", serviceName)),
	)

	l := otelzap.New(zl,
		otelzap.WithMinLevel(zap.DebugLevel),
		otelzap.WithTraceIDField(true))

	otelzap.ReplaceGlobals(l)

	return l.Sync
}
