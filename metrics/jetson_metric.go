//go:build !nojetson
// +build !nojetson

package metrics

// RAM 2015/3964MB (lfb 98x4MB) SWAP 29/1982MB (cached 3MB) CPU [6%@102,5%@102,4%@102,3%@102] EMC_FREQ 0% GR3D_FREQ 0% PLL@28C CPU@31.5C PMIC@100C GPU@30.5C AO@36C thermal@31.25C POM_5V_IN 1388/1388 POM_5V_GPU 122/122 POM_5V_CPU 163/163

// import (
// 	"github.com/prometheus/client_golang/prometheus"
// 	"github.com/prometheus/client_golang/prometheus/promauto"
// 	kingpin "gopkg.in/alecthomas/kingpin.v2"
// )

// type jetsonStat struct {
// 	SWAP    string
// 	IRAM    string
// 	RAM     string
// 	MTS     string
// 	VALS    string
// 	VAL_FRE string
// 	CPU     string
// 	VOLT    string
// 	TEMP    string
// }

// var (
// 	logFilePath    = kingpin.Flag("jetson.logfilepath", "logfilepath.").Default("/tmp/tegrastats.log").String()
// 	tegrastatsPath = kingpin.Flag("jetson.tegrastatspath", "tegrastatspath.").Default("/usr/bin/tegrastats").String()

// 	SWAP    = promauto.NewGauge(prometheus.GaugeOpts{Name: "jetson_swap_info", Help: "swap info of this jetson nano"})
// 	IRAM    = promauto.NewGauge(prometheus.GaugeOpts{Name: "jetson_iram_info", Help: "iram info of this jetson nano"})
// 	RAM     = promauto.NewGauge(prometheus.GaugeOpts{Name: "jetson_ram_info", Help: "ram info of this jetson nano"})
// 	MTS     = promauto.NewGauge(prometheus.GaugeOpts{Name: "jetson_mts_info", Help: "mts info of this jetson nano"})
// 	VALS    = promauto.NewGauge(prometheus.GaugeOpts{Name: "jetson_vals_info", Help: "vals info of this jetson nano"})
// 	VAL_FRE = promauto.NewGauge(prometheus.GaugeOpts{Name: "jetson_val_fre_info", Help: "val_fre info of this jetson nano"})
// 	CPU     = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "jetson_cpu_info", Help: "cpu info of this jetson nano"}, []string{"cpu_id"})
// 	VOLT    = promauto.NewGauge(prometheus.GaugeOpts{Name: "jetson_volt_info", Help: "volt info of this jetson nano"})
// 	TEMP    = promauto.NewGauge(prometheus.GaugeOpts{Name: "jetson_temp_info", Help: "temp info of this jetson nano"})
// )

// func getJetsonMetric() {
// 	// SWAPMatcher, err := regexp.Compile(`SWAP (\d+)\/(\d+)(\w)B( ?)\(cached (\d+)(\w)B\)`)
// 	// IRAMMatcher, err := regexp.Compile(`IRAM (\d+)\/(\d+)(\w)B( ?)\(lfb (\d+)(\w)B\)`)
// 	// RAMMatcher, err := regexp.Compile(`RAM (\d+)\/(\d+)(\w)B( ?)\(lfb (\d+)x(\d+)(\w)B\)`)
// 	// MTSMatcher, err := regexp.Compile(`MTS fg (\d+)% bg (\d+)%`)
// 	// VALSMatcher, err := regexp.Compile(`\b([A-Z0-9_]+) ([0-9%@]+)(?=[^/])\b`)
// 	// VAL_FREMatcher, err := regexp.Compile(`\b(\d+)%@(\d+)`)
// 	// CPUMatcher, err := regexp.Compile(`CPU \[(.*?)\]`)
// 	// VOLTMatcher, err := regexp.Compile(`\b(\w+) ([0-9.]+)\/([0-9.]+)\b`)
// 	// TEMPMatcher, err := regexp.Compile(`\b(\w+)@(-?[0-9.]+)C\b`)
// 	return
// }

// type JetsonMitrics struct{}

// func init() {
// 	registerCollector("jetson_metrics", defaultEnabled, &JetsonMitrics{})
// }

// func jetsonInit() {

// }

// func (j *JetsonMitrics) Update() error {
// 	return nil
// }
