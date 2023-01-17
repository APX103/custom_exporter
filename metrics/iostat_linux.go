//go:build !noiostat
// +build !noiostat

package metrics

import (
	"bufio"
	"fmt"
	"os/exec"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type deviceStats struct {
	device             string
	tps                float64
	kB_read_per_second float64
	kB_wrtn_per_second float64
	kB_read            float64
	kB_wrtn            float64
}

var (
	device = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "iostat_device",
		Help: "iostat_device",
	}, []string{"device"})
	tps = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "iostat_tps",
		Help: "iostat_tps",
	}, []string{"device"})
	kB_read_per_second = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "iostat_kB_read_per_second",
		Help: "iostat_kB_read_per_second",
	}, []string{"device"})
	kB_wrtn_per_second = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "iostat_kB_wrtn_per_second",
		Help: "iostat_kB_wrtn_per_second",
	}, []string{"device"})
	kB_read = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "iostat_kB_read",
		Help: "iostat_kB_read",
	}, []string{"device"})
	kB_wrtn = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "iostat_kB_wrtn",
		Help: "iostat_kB_wrtn",
	}, []string{"device"})
	iostatFormat = "%s %f %f %f %f %f"
)

func (d *IODeviceMetrics) updateIODeviceMetric(newIOStats []deviceStats) {
	d.IODeviceStatsMutex.Lock()
	defer d.IODeviceStatsMutex.Unlock()

	// Reset the cache if the list of CPUs has changed.
	if len(d.IODeviceStats) != len(newIOStats) {
		d.IODeviceStats = make([]deviceStats, len(newIOStats))

	}

	for i, stats := range newIOStats {
		d.IODeviceStats[i] = stats
	}
}

func getIOStats() ([]deviceStats, error) {
	cmd := exec.Command("sh", "-c", "iostat -d")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("[Get iostat info error]: %v \n", err)
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error:The command is err: %v \n", err)
		return nil, err
	}

	outputBuf := bufio.NewReader(stdout)
	outputBuf.ReadLine()
	outputBuf.ReadLine()
	outputBuf.ReadLine()

	var deviceStatslist []deviceStats
	scanner := bufio.NewScanner(outputBuf)
	for scanner.Scan() {
		d := &deviceStats{}
		fmt.Sscanf(scanner.Text(), iostatFormat, &d.device, &d.tps, &d.kB_read_per_second, &d.kB_wrtn_per_second, &d.kB_read, &d.kB_wrtn)
		deviceStatslist = append(deviceStatslist, *d)
	}
	cmd.Wait()
	return deviceStatslist, nil
}

type IODeviceMetrics struct {
	IODeviceStats      []deviceStats
	IODeviceStatsMutex sync.Mutex
}

func init() {
	registerCollector("io_stats_metrics", defaultEnabled, &IODeviceMetrics{})
}

func (iods *IODeviceMetrics) Update() error {
	iostats, err := getIOStats()
	if err != nil {
		fmt.Println("Get err in getIOStats")
		return err
	}
	iods.updateIODeviceMetric(iostats)
	for _, s := range iods.IODeviceStats {
		device.WithLabelValues(s.device).Add(0)
		tps.WithLabelValues(s.device).Set(s.tps)
		kB_read_per_second.WithLabelValues(s.device).Set(s.kB_read_per_second)
		kB_wrtn_per_second.WithLabelValues(s.device).Set(s.kB_wrtn_per_second)
		kB_read.WithLabelValues(s.device).Set(s.kB_read)
		kB_wrtn.WithLabelValues(s.device).Set(s.kB_wrtn)
	}
	return nil
}
