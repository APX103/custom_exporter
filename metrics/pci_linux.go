//go:build !nopci
// +build !nopci

package metrics

// [jetson-exporter] lspci
// 00:02.0 PCI bridge: NVIDIA Corporation Device 0faf (rev a1)
// 01:00.0 Ethernet controller: Realtek Semiconductor Co., Ltd. RTL8111/8168/8411 PCI Express Gigabit Ethernet Controller (rev 15)

import (
	"bufio"
	"os/exec"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
)

var (
	lenPCIList = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "pci_list_len",
		Help: "Length of PCI List of this PC",
	})
	pci_list = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "pci_list_item",
		Help: "List of the pci device",
	}, []string{"device"})
)

func getPCIListMetric() ([]string, error) {
	var pci_list []string
	cmd := exec.Command("lspci")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Debugf("Error:can not obtain stdout pipe for command:%s\nmaybe because there is no lspci", err)
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		log.Debugln("Error:The command is err,", err)
		return nil, err
	}

	outputBuf := bufio.NewReader(stdout)
	for {
		output, _, err := outputBuf.ReadLine()

		if err != nil {
			break
		}
		pci_list = append(pci_list, string(output))
	}
	cmd.Wait()

	return pci_list, nil
}

type PCIMitrics struct {
	PCIDevicelist       []string
	PCIDeviceStatsMutex sync.Mutex
}

func (p *PCIMitrics) updatePCIMitrics(pci_list []string) {
	p.PCIDeviceStatsMutex.Lock()
	defer p.PCIDeviceStatsMutex.Unlock()

	if len(p.PCIDevicelist) != len(pci_list) {
		p.PCIDevicelist = make([]string, len(pci_list))

	}

	for i, device := range pci_list {
		p.PCIDevicelist[i] = device
	}
}

func init() {
	registerCollector("pci_metrics", defaultEnabled, &PCIMitrics{})
}

func (p *PCIMitrics) Update() error {
	lp, err := getPCIListMetric()
	if err != nil {
		log.Debugf("Get PCI Mitrics Error")
		return err
	}
	p.updatePCIMitrics(lp)
	lenPCIList.Set(float64(len(lp)))
	for _, i := range p.PCIDevicelist {
		pci_list.WithLabelValues(i).Add(0)
	}
	return nil
}
