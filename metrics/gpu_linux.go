//go:build !nogpu
// +build !nogpu

package metrics

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type GpuInfo struct {
	GPUID             string
	GpuUtilization    float64
	GpuMemUtilization float64
	GpuMemTotal       float64
	GpuMemUsed        float64
	GpuMemFree        float64
	GpuName           string
}

var (
	gpuid             = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "GPUID", Help: "gpu id of this PC"}, []string{"gpu_id", "gpu_name"})
	gpuUtilization    = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "GPUUtilization", Help: "gpu utilization of this PC"}, []string{"gpu_id"})
	gpuMemUtilization = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "GPUMemUtilization", Help: "gpu mem utilization of this PC"}, []string{"gpu_id"})
	gpuMemTotal       = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "GPUMemTotal", Help: "gpu mem of this PC"}, []string{"gpu_id"})
	gpuMemUsed        = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "GPUMemUsed", Help: "gpu mem used of this PC"}, []string{"gpu_id"})
	gpuMemFree        = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "GPUMemFree", Help: "gpu mem free of this PC"}, []string{"gpu_id"})
)

func Str2Num(s string) float64 {
	pattern := regexp.MustCompile(`(\d+)`)
	numberStrings := pattern.FindAllStringSubmatch(strings.Split(s, ", ")[0], -1)
	numbers := make([]float64, len(numberStrings))
	for i, numberString := range numberStrings {
		number, err := strconv.ParseFloat(numberString[1], 64)
		if err != nil {
			panic(err)
		}
		numbers[i] = number
	}
	return numbers[0]
}

func GetGpuInfo() ([]GpuInfo, error) {
	cmd := exec.Command("nvidia-smi", "--query-gpu=memory.total,memory.free,memory.used,utilization.gpu,utilization.memory,name", "--format=csv")

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error:can not obtain stdout pipe for command:%s\nmaybe because there is no nvidia-smi", err)
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Error:The command is err,", err)
		return nil, err
	}

	outputBuf := bufio.NewReader(stdout)
	outputBuf.ReadLine()
	var info []string
	for {
		output, _, err := outputBuf.ReadLine()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, err
		}
		info = append(info, string(output))
	}
	cmd.Wait()
	var gpuInfoList []GpuInfo

	for i, output := range info {
		r := strings.Split(string(output), ", ")
		gpuInfoList = append(gpuInfoList, GpuInfo{
			GPUID:             strconv.Itoa(i),
			GpuUtilization:    Str2Num(r[3]),
			GpuMemUtilization: Str2Num(r[4]),
			GpuMemTotal:       Str2Num(r[0]),
			GpuMemUsed:        Str2Num(r[2]),
			GpuMemFree:        Str2Num(r[1]),
			GpuName:           r[5],
		})
	}
	return gpuInfoList, nil
}

type GPUMitrics struct {
	GPUDevices           []GpuInfo
	GPUDevicesStatsMutex sync.Mutex
}

func (g *GPUMitrics) updateGPUMitrics(gpu_metrics []GpuInfo) {
	g.GPUDevicesStatsMutex.Lock()
	defer g.GPUDevicesStatsMutex.Unlock()

	if len(g.GPUDevices) != len(gpu_metrics) {
		g.GPUDevices = make([]GpuInfo, len(gpu_metrics))

	}

	copy(g.GPUDevices, gpu_metrics)
}

func init() {
	registerCollector("gpu_metrics", defaultEnabled, &GPUMitrics{})
}

func (g *GPUMitrics) Update() error {
	gi, err := GetGpuInfo()
	if err != nil {
		fmt.Println("Get GPU Mitrics Error")
		return err
	}
	g.updateGPUMitrics(gi)
	for _, i := range g.GPUDevices {
		gpuid.WithLabelValues(i.GPUID, i.GpuName).Add(0)
		gpuUtilization.WithLabelValues(i.GPUID).Set(i.GpuUtilization)
		gpuMemUtilization.WithLabelValues(i.GPUID).Set(i.GpuMemUtilization)
		gpuMemTotal.WithLabelValues(i.GPUID).Set(i.GpuMemTotal)
		gpuMemUsed.WithLabelValues(i.GPUID).Set(i.GpuMemUsed)
		gpuMemFree.WithLabelValues(i.GPUID).Set(i.GpuMemFree)
	}
	return nil
}
