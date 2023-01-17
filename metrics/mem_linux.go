package metrics

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/shirou/gopsutil/v3/mem"
)

var (
	memUtilization = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "mem_utilization",
		Help: "Mem utilization of this PC",
	})
	memTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "mem_total",
		Help: "Mem Total of this PC",
	})
)

// type CPUMitrics struct{}

type MemMitrics struct{}

func init() {
	registerCollector("mem_metrics", defaultEnabled, &MemMitrics{})
}

func getMemUtilizationMetric() (float64, error) {
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}
	return memInfo.UsedPercent, err
}

func getMemTotalMetric() (float64, error) {
	cmd := exec.Command("sh", "-c", "cat /proc/meminfo | grep MemTotal | grep -Eo '[0-9]+'")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("[Get mem total info error]: %v \n", err)
		return 0, err
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error:The command is err: %v \n", err)
		return 0, err
	}

	outputBuf := bufio.NewReader(stdout)
	output, _, err := outputBuf.ReadLine()
	cmd.Wait()
	if err != nil {
		fmt.Println("Read Mem Total failed!")
		return 0, err
	}
	cn, err := strconv.Atoi(string(output))
	if err != nil {
		fmt.Println("Atoi error in getMemTotalMetric")
		return 0, err
	}
	return float64(cn), nil
}

func (m *MemMitrics) Update() error {
	mu, err := getMemUtilizationMetric()
	if err != nil {
		return err
	}
	memUtilization.Set(mu)
	mt, err := getMemTotalMetric()
	if err != nil {
		return err
	}
	memTotal.Set(mt)

	return nil
}
