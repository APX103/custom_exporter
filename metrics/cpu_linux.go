package metrics

import (
	"bufio"
	"os/exec"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/shirou/gopsutil/v3/cpu"
	log "github.com/sirupsen/logrus"
)

var (
	cpuUtilization = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_utilization",
		Help: "CPU utilization of this PC",
	})
	cpuCoreNum = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "cpu_core_number",
		Help: "CPU Core Number of this PC",
	})
)

type CPUMetrics struct{}

func init() {
	registerCollector("cpu_metrics", defaultEnabled, &CPUMetrics{})
}

func getCPUUtilizationMetric() (float64, error) {
	percent, err := cpu.Percent(time.Second, false)
	if err != nil {
		return 0, err
	}
	return float64(percent[0]), nil
}

func getCPUCoreNumMetric() (float64, error) {
	cmd := exec.Command("sh", "-c", "cat /proc/cpuinfo|grep processor|wc -l")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Debugf("[Get cpu info error]: %v \n", err)
		return 0, err
	}

	if err := cmd.Start(); err != nil {
		log.Debugf("Error:The command is err: %v \n", err)
		return 0, err
	}

	outputBuf := bufio.NewReader(stdout)
	output, _, err := outputBuf.ReadLine()
	cmd.Wait()
	if err != nil {
		log.Debugf("Read CPU Core failed!")
		return 0, err
	}
	cn, err := strconv.Atoi(string(output))
	if err != nil {
		log.Debugf("Atoi error in getCPUCoreNumMetric")
		return 0, err
	}

	// now cn is number of cpu, int
	return float64(cn), nil
}

func (m *CPUMetrics) Update() error {
	cu, err := getCPUUtilizationMetric()
	cn, err := getCPUCoreNumMetric()
	if err != nil {
		return err
	}
	cpuUtilization.Set(cu)
	cpuCoreNum.Set(cn)
	return nil
}
