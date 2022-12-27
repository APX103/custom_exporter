package metrics

// [jetson-exporter] lspci
// 00:02.0 PCI bridge: NVIDIA Corporation Device 0faf (rev a1)
// 01:00.0 Ethernet controller: Realtek Semiconductor Co., Ltd. RTL8111/8168/8411 PCI Express Gigabit Ethernet Controller (rev 15)

import (
	"bufio"
	"fmt"
	"os/exec"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	pciList = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "pci_list_len",
		Help: "Length of PCI List of this PC",
	})
)

func getPCIListMetric() (float64, error) {
	var pci_list []string
	cmd := exec.Command("lspci")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("Error:can not obtain stdout pipe for command:%s\nmaybe because there is no lspci", err)
		return 0, err
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Error:The command is err,", err)
		return 0, err
	}

	outputBuf := bufio.NewReader(stdout)
	for {
		output, _, err := outputBuf.ReadLine()

		if err != nil {
			break
		}
		// fmt.Printf(string(output) + "\n")
		pci_list = append(pci_list, string(output))
	}
	cmd.Wait()
	return float64(len(pci_list)), nil
}

type PCIMitrics struct{}

func init() {
	registerCollector("pci_metrics", defaultEnabled, &PCIMitrics{})
}

func (m *PCIMitrics) Update() error {
	lp, err := getPCIListMetric()
	if err != nil {
		fmt.Println("Get PCI Mitrics Error")
		return err
	}
	pciList.Set(lp)
	return nil
}
