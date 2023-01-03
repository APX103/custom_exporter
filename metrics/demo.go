package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	demoNum = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "demo_number",
		Help: "DEMO NUmber",
	})
)

type DemoMetrics struct{}

func init() {
	registerCollector("demo_metrics", defaultEnabled, &DemoMetrics{})
}

func (m *DemoMetrics) Update() error {
	demoNum.Set(1)
	return nil
}
