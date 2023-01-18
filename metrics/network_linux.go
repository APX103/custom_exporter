package metrics

import (
	"bufio"
	"fmt"
	"os/exec"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type network struct {
	device    string
	ipv4_addr string
	receive   string
	transmit  string
}

var (
	net_dev = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "net_dev",
		Help: "net_dev, and it's ipv4 addr",
	}, []string{"net_dev", "ipv4_addr"})
	receive = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "receive",
		Help: "receive",
	}, []string{"net_dev"})
	transmit = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "transmit",
		Help: "transmit",
	}, []string{"net_dev"})
)

type NetWorkMetrics struct {
	networks []network
}

func init() {
	registerCollector("NetWork_metrics", defaultEnabled, &NetWorkMetrics{})
}

func getIPv4Addr(device string) (string, error) {
	_cmd := fmt.Sprintf(`ip -4 addr | grep inet | grep %v | grep -Eo "[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+" | head -n 1`, device)
	cmd := exec.Command("sh", "-c", _cmd)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("[Get ipv4 addr error]: %v \n", err)
		return "", err
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error:The command is err: %v \n", err)
		return "", err
	}

	outputBuf := bufio.NewReader(stdout)
	output, _, err := outputBuf.ReadLine()
	cmd.Wait()
	if err != nil {
		// fmt.Println("Get ipv4 addr failed!")
		return "", err
	}
	return string(output), nil
}

func getReceive(device string) (string, error) {
	_cmd := fmt.Sprintf(`cat /proc/net/dev | grep %v | awk '{print $2}' | grep -Eo "[0-9]+"`, device)
	cmd := exec.Command("sh", "-c", _cmd)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("[Get receive error]: %v \n", err)
		return "", err
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error:The command is err: %v \n", err)
		return "", err
	}

	outputBuf := bufio.NewReader(stdout)
	output, _, err := outputBuf.ReadLine()
	cmd.Wait()
	if err != nil {
		fmt.Println("Get receive failed!")
		return "", err
	}
	return string(output), nil
}

func getTransmit(device string) (string, error) {
	_cmd := fmt.Sprintf(`cat /proc/net/dev | grep %v | awk '{print $10}' | grep -Eo "[0-9]+"`, device)
	cmd := exec.Command("sh", "-c", _cmd)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("[Get transmit error]: %v \n", err)
		return "", err
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error:The command is err: %v \n", err)
		return "", err
	}

	outputBuf := bufio.NewReader(stdout)
	output, _, err := outputBuf.ReadLine()
	cmd.Wait()
	if err != nil {
		fmt.Println("Get transmit failed!")
		return "", err
	}
	return string(output), nil
}

func getNetWork() ([]network, error) {
	cmd := exec.Command("sh", "-c", "for i in $(cat /proc/net/dev | grep 0 | grep -v lo | awk '{print $1}'); do echo ${i%?};done")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Printf("[Get network error]: %v \n", err)
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("Error:The command is err: %v \n", err)
		return nil, err
	}

	outputBuf := bufio.NewReader(stdout)

	var network_list []network

	for {
		rl, _, err := outputBuf.ReadLine()
		if err != nil {
			break
		}
		de := string(rl)
		ia, _ := getIPv4Addr(de)
		rx, _ := getReceive(de)
		tx, _ := getTransmit(de)

		n := &network{
			device:    de,
			ipv4_addr: ia,
			receive:   rx,
			transmit:  tx,
		}
		network_list = append(network_list, *n)
	}
	return network_list, nil
}

func (n *NetWorkMetrics) updateNetWorkMetrics(new []network) {
	if len(n.networks) != len(new) {
		n.networks = make([]network, len(new))
	}

	copy(n.networks, new)
}

func (n *NetWorkMetrics) Update() error {
	networks, err := getNetWork()
	if err != nil {
		fmt.Println("Get err in getNetWork")
		return err
	}
	n.updateNetWorkMetrics(networks)
	for _, s := range n.networks {
		_receive, err := strconv.ParseFloat(s.receive, 64)
		if err != nil {
			fmt.Println("parse receive error")
			return err
		}
		_transmit, err := strconv.ParseFloat(s.transmit, 64)
		if err != nil {
			fmt.Println("parse transmit error")
			return err
		}
		net_dev.WithLabelValues(s.device, s.ipv4_addr).Add(0)
		receive.WithLabelValues(s.device).Set(_receive)
		transmit.WithLabelValues(s.device).Set(_transmit)
	}
	return nil

}
