package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/prometheus/exporter-toolkit/web/kingpinflag"
	log "github.com/sirupsen/logrus"
	kingpin "gopkg.in/alecthomas/kingpin.v2"

	"pjlab.org.cn/qa/qa_exporter/metrics"
)

var (
	metricsPath = kingpin.Flag(
		"web.telemetry-path",
		"Path under which to expose metrics.",
	).Default("/metrics").String()
	maxProcs = kingpin.Flag(
		"runtime.gomaxprocs",
		"The target number of CPUs Go will run on (GOMAXPROCS)",
	).Envar("GOMAXPROCS").Default("1").Int()
	logLevel = kingpin.Flag(
		"log.loglevel",
		"loglevel of this app",
	).Default("info").String()
	toolkitFlags = kingpinflag.AddFlags(kingpin.CommandLine, ":2112")
)

func UpdateMetrics() {
	for key, m := range metrics.MetricsState {
		if *m {
			err := metrics.InitiatedMetrics[key].Update()
			if err != nil {
				log.Debug(fmt.Sprintf("Get Param Err: %v \n", err))
			}
		}

	}
}

func recordMetrics() {
	go func() {
		for {
			UpdateMetrics()
			time.Sleep(2 * time.Second)
		}
	}()
}

func parseLogLevel(l string) log.Level {
	switch l {
	case "info":
		return log.InfoLevel
	case "debug":
		return log.DebugLevel
	case "error":
		return log.ErrorLevel
	default:
		return log.FatalLevel
	}
}

func main() {
	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Version(version.Print("deploy_platform_exporter"))
	kingpin.CommandLine.UsageWriter(os.Stdout)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(parseLogLevel(*logLevel))

	log.Infof("server start info")
	log.Errorf("server start error")
	log.Debugf("server start debug")

	recordMetrics()

	logger := promlog.New(promlogConfig)
	runtime.GOMAXPROCS(*maxProcs)
	level.Debug(logger).Log("msg", "Go MAXPROCS", "procs", runtime.GOMAXPROCS(0))
	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>Node Exporter</title></head>
			<body>
			<h1>Node Exporter</h1>
			<p><a href="` + *metricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})
	server := &http.Server{}
	if err := web.ListenAndServe(server, toolkitFlags, logger); err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}
}
