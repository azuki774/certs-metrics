package metrics

import (
	"certs-metrics/internal/usecase"
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type MetricsServer struct {
	Logger *zap.Logger
	Us     *usecase.Usecase
}

type metrics struct {
	dir                  string // ca.crt directory, for labeling
	notBefore            prometheus.Gauge
	notAfter             prometheus.Gauge
	validPeriod          prometheus.Gauge
	remainingValidPeriod prometheus.Gauge
}

func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		notBefore: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "cert_not_before",
			Help: "",
		}),
		notAfter: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "cert_not_after",
			Help: "",
		}),
		validPeriod: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "cert_valid_period",
			Help: "",
		}),
		remainingValidPeriod: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "cert_remaining_valid_period",
			Help: "",
		}),
	}
	reg.MustRegister(m.notBefore)
	reg.MustRegister(m.notAfter)
	reg.MustRegister(m.validPeriod)
	reg.MustRegister(m.remainingValidPeriod)
	return m
}

// 60 秒ごとに情報を更新
func refresh(m *metrics) {
	for {
		rand.Seed(time.Now().UnixNano())
		n := -1 + rand.Float64()*2
		m.notBefore.Set(n)
		m.notAfter.Set(n)
		m.validPeriod.Set(n)
		m.remainingValidPeriod.Set(n)
		time.Sleep(5 * time.Second)
	}
}

func (s *MetricsServer) Start(ctx context.Context) error {
	s.Logger.Info("cert-metrics server start")
	// Create a non-global registry.
	reg := prometheus.NewRegistry()

	// Create new metrics and register them using the custom registry.
	m := NewMetrics(reg)
	m.notBefore.Set(123.4)
	go refresh(m)

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	http.ListenAndServe(":8334", nil)

	return nil
}
