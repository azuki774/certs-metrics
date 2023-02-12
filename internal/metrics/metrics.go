package metrics

import (
	"certs-metrics/internal/usecase"
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type MetricsServer struct {
	Logger *zap.Logger
	Us     *usecase.Usecase
	Port   string
	Dirs   []string // ca.crt directory, for labeling
}

type metrics struct {
	notBefore            prometheus.GaugeVec
	notAfter             prometheus.GaugeVec
	validPeriod          prometheus.GaugeVec
	remainingValidPeriod prometheus.GaugeVec
}

func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		notBefore: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "cert_not_before",
			Help: "certification notBefore value (unixtime)",
		},
			[]string{"cert_name", "cert_full_path"},
		),
		notAfter: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "cert_not_after",
			Help: "certification notAfter value (unixtime)",
		},
			[]string{"cert_name", "cert_full_path"},
		),
		validPeriod: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "cert_valid_period",
			Help: "certification valid period (min)",
		},
			[]string{"cert_name", "cert_full_path"},
		),
		remainingValidPeriod: *prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "cert_remaining_valid_period",
			Help: "certification remaining valid period (min)",
		},
			[]string{"cert_name", "cert_full_path"},
		),
	}
	reg.MustRegister(m.notBefore)
	reg.MustRegister(m.notAfter)
	reg.MustRegister(m.validPeriod)
	reg.MustRegister(m.remainingValidPeriod)
	return m
}

// refresh: 60 秒ごとに情報を更新
func (s *MetricsServer) refresh(ctx context.Context, m *metrics) {
	go func() {
		for {
			cis, err := s.Us.LoadCertsInfo(s.Dirs)
			if err != nil {
				s.Logger.Error("fetch error", zap.Error(err))
			}

			for _, ci := range cis {
				pl := prometheus.Labels{"cert_name": ci.FileName, "cert_full_path": ci.FullPath}

				m.notBefore.With(pl).Set(float64(ci.NotBefore.Unix()))
				m.notAfter.With(pl).Set(float64(ci.NotAfter.Unix()))
				m.validPeriod.With(pl).Set(ci.ValidPeriod.Minutes())
				m.remainingValidPeriod.With(pl).Set(ci.RemainingValidPeriod.Minutes())
			}
			time.Sleep(60 * time.Second)
		}
	}()
	<-ctx.Done()
	s.Logger.Info("refresh routine close")
}

func (s *MetricsServer) Start(ctx context.Context) error {
	s.Logger.Info("cert-metrics server start", zap.String("port", s.Port))
	// Create a non-global registry.
	reg := prometheus.NewRegistry()

	// Create new metrics and register them using the custom registry.
	m := NewMetrics(reg)

	go s.refresh(ctx, m)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", s.Port),
		Handler: nil,
	}
	go func() {
		<-ctx.Done()
		s.Logger.Info("shutdown signal catch")
		ctx2, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		nerr := server.Shutdown(ctx2)
		if nerr != nil {
			s.Logger.Error("gracefully shutdown error", zap.Error(nerr))
		}
	}()

	s.Logger.Info("start listening", zap.String("port", s.Port))
	err := server.ListenAndServe()

	if errors.Is(err, http.ErrServerClosed) {
		// expected error
		err = nil
	} else {
		s.Logger.Error("metrics server close error", zap.Error(err))
		return err
	}

	s.Logger.Info("metrics server shutdown")
	return nil
}
