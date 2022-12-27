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

type DemoMitrics struct{}

func init() {
	registerCollector("demo_metrics", defaultEnabled, &DemoMitrics{})
}

func (m *DemoMitrics) Update() error {
	demoNum.Set(1)
	return nil
}
