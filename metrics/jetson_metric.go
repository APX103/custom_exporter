//go:build !nojetson
// +build !nojetson

package metrics

// RAM 2015/3964MB (lfb 98x4MB) SWAP 29/1982MB (cached 3MB) CPU [6%@102,5%@102,4%@102,3%@102] EMC_FREQ 0% GR3D_FREQ 0% PLL@28C CPU@31.5C PMIC@100C GPU@30.5C AO@36C thermal@31.25C POM_5V_IN 1388/1388 POM_5V_GPU 122/122 POM_5V_CPU 163/163

import (
	"bufio"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/dlclark/regexp2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	log "github.com/sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

type jetsonStat struct {
	SWAP map[string]float64
	IRAM map[string]float64
	RAM  map[string]float64
	MTS  map[string]float64
	VALS map[string]float64
	CPU  map[string]float64
	VOLT map[string]float64
	TEMP map[string]float64
}

var (
	logFilePath    = kingpin.Flag("jetson.logfilepath", "logfilepath.").Default("/tmp/tegrastats.log").String()
	tegrastatsPath = kingpin.Flag("jetson.tegrastatspath", "tegrastatspath.").Default("/usr/bin/tegrastats").String()

	SWAP = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "jetson_swap_info", Help: "swap info of this jetson nano"}, []string{"val"})
	IRAM = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "jetson_iram_info", Help: "iram info of this jetson nano"}, []string{"val"})
	RAM  = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "jetson_ram_info", Help: "ram info of this jetson nano"}, []string{"val"})
	MTS  = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "jetson_mts_info", Help: "mts info of this jetson nano"}, []string{"val"})
	VALS = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "jetson_vals_info", Help: "vals info of this jetson nano"}, []string{"val"})
	CPU  = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "jetson_cpu_info", Help: "cpu info of this jetson nano"}, []string{"cpu_id"})
	VOLT = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "jetson_volt_info", Help: "volt info of this jetson nano"}, []string{"name"})
	TEMP = promauto.NewGaugeVec(prometheus.GaugeOpts{Name: "jetson_temp_info", Help: "temp info of this jetson nano"}, []string{"name"})
)

func regexp2FindAllString(re *regexp2.Regexp, s string) map[string]string {
	matches := make(map[string]string)
	m, _ := re.FindStringMatch(s)

	for m != nil {
		p := m.Groups()
		matches[p[1].Captures[0].String()] = p[2].Captures[0].String()
		m, _ = re.FindNextMatch(m)
	}
	return matches
}

func getJetsonMetric(str string) jetsonStat {
	SWAPMatcher, _ := regexp.Compile(`SWAP (?P<swap_use>\d+)\/(?P<swap_tot>\d+)(?P<swap_unit>\w)B( ?)\(cached (?P<swap_cached_size>\d+)(?P<swap_cached_unit>\w)B\)`)
	IRAMMatcher, _ := regexp.Compile(`IRAM (?P<iram_use>\d+)\/(?P<iram_tot>\d+)(?P<iram_unit>\w)B( ?)\(lfb (?P<iram_lfb_size>\d+)(?P<iram_lfb_unit>\w)B\)`)
	RAMMatcher, _ := regexp.Compile(`RAM (?P<ram_use>\d+)\/(?P<ram_tot>\d+)(?P<ram_unit>\w)B( ?)\(lfb (?P<ram_lfb_nblock>\d+)x(?P<ram_lfb_size>\d+)(?P<ram_lfb_unit>\w)B\)`)
	MTSMatcher, _ := regexp.Compile(`MTS fg (?P<mts_fg>\d+)% bg (?P<mts_bg>\d+)%`)
	VALSMatcher, _ := regexp2.Compile(`\b([A-Z0-9_]+) ([0-9%@]+)(?=[^/])\b`, 0)
	VAL_FREMatcher, _ := regexp.Compile(`\b(?P<val>\d+)%@(?P<frq>\d+)`)
	CPUMatcher, _ := regexp.Compile(`CPU \[(.*?)\]`)
	VOLTMatcher, _ := regexp.Compile(`\b(?P<name>\w+) (?P<cur>[0-9.]+)\/(?P<avg>[0-9.]+)\b`)
	TEMPMatcher, _ := regexp.Compile(`\b(?P<name>\w+)@(?P<val>-?[0-9.]+)C\b`)

	return jetsonStat{
		SWAP: getSWAP(SWAPMatcher, str),
		IRAM: getIRAM(IRAMMatcher, str),
		RAM:  getRAM(RAMMatcher, str),
		MTS:  getMTS(MTSMatcher, str),
		VALS: getVALS(VALSMatcher, VAL_FREMatcher, str),
		CPU:  getCPU(CPUMatcher, VAL_FREMatcher, str),
		VOLT: getVOLT(VOLTMatcher, str),
		TEMP: getTEMP(TEMPMatcher, str),
	}
}

// TODO Add unit to vals
func getSWAP(m *regexp.Regexp, txt string) map[string]float64 {
	res := make(map[string]float64)
	sub_list := m.SubexpNames()
	vals := m.FindStringSubmatch(txt)
	if len(vals) == 0 {
		return res
	}
	for i, vl := range sub_list {
		if vl == "swap_use" || vl == "swap_tot" || vl == "swap_cached_size" {
			val, err := strconv.ParseFloat(vals[i], 64)
			if err != nil {
				panic("parse swap value err, can not parse to float")
			}
			res[vl] = val
		}
	}
	return res
}

func getIRAM(m *regexp.Regexp, txt string) map[string]float64 {
	res := make(map[string]float64)
	sub_list := m.SubexpNames()
	vals := m.FindStringSubmatch(txt)
	if len(vals) == 0 {
		return res
	}
	for i, vl := range sub_list {
		if vl == "iram_use" || vl == "iram_tot" || vl == "iram_lfb_size" {
			val, err := strconv.ParseFloat(vals[i], 64)
			if err != nil {
				panic("parse iram value err, can not parse to float")
			}
			res[vl] = val
		}
	}
	return res
}

func getRAM(m *regexp.Regexp, txt string) map[string]float64 {
	res := make(map[string]float64)
	sub_list := m.SubexpNames()
	vals := m.FindStringSubmatch(txt)
	if len(vals) == 0 {
		return res
	}
	for i, vl := range sub_list {
		if vl == "ram_use" || vl == "ram_tot" || vl == "ram_lfb_nblock" || vl == "ram_lfb_size" {
			val, err := strconv.ParseFloat(vals[i], 64)
			if err != nil {
				panic("parse ram value err, can not parse to float")
			}
			res[vl] = val
		}
	}
	return res
}

func getMTS(m *regexp.Regexp, txt string) map[string]float64 {
	res := make(map[string]float64)
	sub_list := m.SubexpNames()
	vals := m.FindStringSubmatch(txt)
	if len(vals) == 0 {
		return res
	}
	for i, vl := range sub_list {
		if vl == "mts_fg" || vl == "mts_bg" {
			val, err := strconv.ParseFloat(vals[i], 64)
			if err != nil {
				panic("parse mts value err, can not parse to float")
			}
			res[vl] = val
		}
	}
	return res
}

func getVALS(m *regexp2.Regexp, vf *regexp.Regexp, txt string) map[string]float64 {
	res := make(map[string]float64)
	matches := regexp2FindAllString(m, txt)
	for k, v := range matches {
		vmap := getVAL_FRE(vf, v)
		val, ok := vmap["val"]
		if ok {
			res[k+"_val"] = val
			res[k+"_frq"] = vmap["frq"]
		} else {
			res[k], _ = strconv.ParseFloat(v, 64)
		}

	}
	return res
}

func getVAL_FRE(m *regexp.Regexp, txt string) map[string]float64 {
	res := make(map[string]float64)
	sub_list := m.SubexpNames()
	vals := m.FindStringSubmatch(txt)
	if len(vals) == 0 {
		return res
	}
	if strings.Contains(txt, "@") {
		for i, vl := range sub_list {
			if vl == "val" || vl == "frq" {
				val, err := strconv.ParseFloat(vals[i], 64)
				if err != nil {
					panic("parse val_frq value err, can not parse to float")
				}
				res[vl] = val
			}
		}
	} else {
		res["frq"], _ = strconv.ParseFloat(vals[0], 64)
	}
	return res
}

func getCPU(m *regexp.Regexp, vf *regexp.Regexp, txt string) map[string]float64 {
	k := "cpu_"
	res := make(map[string]float64)
	_vals := m.FindStringSubmatch(txt)
	if len(_vals) == 0 {
		return res
	}
	vals := strings.Split(_vals[1], ",")
	for i, val := range vals {
		vmap := getVAL_FRE(vf, val)
		v, ok := vmap["val"]
		if ok {
			res[k+strconv.Itoa(i)+"_val"] = v
			res[k+strconv.Itoa(i)+"_frq"] = vmap["frq"]
		} else {
			res[k+strconv.Itoa(i)], _ = vmap["frq"]
		}
	}
	return res
}

func getVOLT(m *regexp.Regexp, txt string) map[string]float64 {
	res := make(map[string]float64)
	sub_list := m.SubexpNames()
	items := m.FindAllString(txt, 10)
	if len(items) == 0 {
		return res
	}
	for _, i := range items {
		vals := m.FindStringSubmatch(i)
		res[vals[1]+"_"+sub_list[2]], _ = strconv.ParseFloat(vals[2], 64)
		res[vals[1]+"_"+sub_list[3]], _ = strconv.ParseFloat(vals[3], 64)
	}
	return res
}

func getTEMP(m *regexp.Regexp, txt string) map[string]float64 {
	res := make(map[string]float64)
	sub_list := m.SubexpNames()
	items := m.FindAllString(txt, 10)
	if len(items) == 0 {
		return res
	}
	for _, i := range items {
		vals := m.FindStringSubmatch(i)
		res[vals[1]+"_"+sub_list[2]], _ = strconv.ParseFloat(vals[2], 64)
	}
	return res
}

func updateJetsonMetric(s jetsonStat) {
	for k, v := range s.SWAP {
		SWAP.WithLabelValues(k).Set(v)
	}
	for k, v := range s.IRAM {
		IRAM.WithLabelValues(k).Set(v)
	}
	for k, v := range s.RAM {
		RAM.WithLabelValues(k).Set(v)
	}
	for k, v := range s.MTS {
		MTS.WithLabelValues(k).Set(v)
	}
	for k, v := range s.VALS {
		VALS.WithLabelValues(k).Set(v)
	}
	for k, v := range s.CPU {
		CPU.WithLabelValues(k).Set(v)
	}
	for k, v := range s.VOLT {
		VOLT.WithLabelValues(k).Set(v)
	}
	for k, v := range s.TEMP {
		TEMP.WithLabelValues(k).Set(v)
	}
}

type JetsonMitrics struct {
	jetsonCmdFlag bool
}

func init() {
	registerCollector("jetson_metrics", defaultEnabled, &JetsonMitrics{jetsonCmdFlag: false})
}

func jetsonInit() error {
	cmd := exec.Command("tegrastats", "--stop")
	_, err := cmd.StdoutPipe()
	if err != nil {
		log.Errorf("[jetsonInit stop error]: %v \n", err)
		return err
	}

	if err := cmd.Start(); err != nil {
		log.Errorf("Error:The jetsonInit stop command is err: %v \n", err)
		return err
	}
	cmd.Wait()
	cmd = exec.Command(*tegrastatsPath, "--logfile", *logFilePath, "--start")
	_, err = cmd.StdoutPipe()
	if err != nil {
		log.Errorf("[jetsonInit start error]: %v \n", err)
		return err
	}

	if err := cmd.Start(); err != nil {
		log.Errorf("Error:The jetsonInit start command is err: %v \n", err)
		return err
	}
	cmd.Wait()
	return nil
}

func getTegrastatStr(log_path string) (string, error) {
	cmd := exec.Command("tail", "-n", "1", *logFilePath)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Errorf("[jetsonInit stop error]: %v \n", err)
		return "", err
	}

	if err := cmd.Start(); err != nil {
		log.Errorf("Error:The jetsonInit stop command is err: %v \n", err)
		return "", err
	}
	outputBuf := bufio.NewReader(stdout)
	output, _, err := outputBuf.ReadLine()
	cmd.Wait()
	return string(output), nil
}

func (j *JetsonMitrics) Update() error {
	if !j.jetsonCmdFlag {
		err := jetsonInit()
		if err != nil {
			return nil
		}
		j.jetsonCmdFlag = true
	}
	s, err := getTegrastatStr(*logFilePath)
	if err != nil {
		return nil
	}
	i := getJetsonMetric(s)
	updateJetsonMetric(i)
	return nil
}
