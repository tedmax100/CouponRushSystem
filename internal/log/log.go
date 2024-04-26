package log

import (
	"context"
	"errors"
	"fmt"
	l "log"
	"syscall"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	SEVERITY_INFRA_DEPENDENCY = 21
)

var (
	log *otelzap.Logger
)

func Initital(gitCommit, gitBranch, goVersion, buildDate string) {
	config := zap.NewProductionConfig()
	config.Level.SetLevel(zap.DebugLevel)
	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = EncodeTime
	logger, err := config.Build(zap.AddCaller(), zap.AddCallerSkip(1),
		zap.Fields(zap.String("service", "coupon_rush_server")),
		zap.Fields(zap.String("commit", gitCommit)),
		zap.Fields(zap.String("branch", gitBranch)),
		zap.Fields(zap.String("go", goVersion)),
		zap.Fields(zap.String("build", buildDate)))

	if err != nil {
		l.Println(err)
	}
	err = logger.Sync()
	if err != nil && !errors.Is(err, syscall.EINVAL) {
		l.Println(err)
	}
	log = otelzap.New(logger, otelzap.WithTraceIDField(true))
}

func EncodeTime(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	encoder.AppendInt64(t.UnixMilli())
}

func Debug(ctx context.Context, msg string, args ...zap.Field) {
	log.DebugContext(ctx, msg, args...)
}

func Info(ctx context.Context, msg string, args ...zap.Field) {
	log.InfoContext(ctx, msg, args...)
}

func Warn(ctx context.Context, msg string, args ...zap.Field) {
	log.WarnContext(ctx, msg, args...)
}

func Error(ctx context.Context, err error, args ...zap.Field) {
	log.ErrorContext(ctx, err.Error(), args...)
}

func Fatal(ctx context.Context, err error, args ...zap.Field) {
	log.FatalContext(ctx, err.Error(), append(args, zap.Int("severity", SEVERITY_INFRA_DEPENDENCY))...)
}

func Sync() {
	err := log.Sync()
	if err != nil && !errors.Is(err, syscall.EINVAL) {
		log.Error(err.Error())
	}
}

type FxLogger struct{}

func (*FxLogger) Printf(msg string, args ...interface{}) {
	log.Debug(fmt.Sprintf(msg, args...), zap.String("module", "fx"))
}
