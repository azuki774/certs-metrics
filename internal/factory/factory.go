package factory

import (
	"certs-metrics/internal/metrics"
	"certs-metrics/internal/repository"
	"certs-metrics/internal/usecase"
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger() (*zap.Logger, error) {
	config := zap.NewProductionConfig()
	// config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	config.EncoderConfig.EncodeTime = JSTTimeEncoder
	l, err := config.Build()

	l.WithOptions(zap.AddStacktrace(zap.ErrorLevel))
	if err != nil {
		fmt.Printf("failed to create logger: %v\n", err)
	}
	return l, err
}

func JSTTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	const layout = "2006-01-02T15:04:05+09:00"
	jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	enc.AppendString(t.In(jst).Format(layout))
}

func NewUsecase(l *zap.Logger) (us *usecase.Usecase) {
	return &usecase.Usecase{L: l, C: &repository.CertsLoader{}}
}

func NewMetricsServer(l *zap.Logger, us *usecase.Usecase, dirs []string) (ms *metrics.MetricsServer) {
	return &metrics.MetricsServer{Logger: l, Us: us, Dirs: dirs}
}
