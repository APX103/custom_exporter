package metrics

import (
	"fmt"
	"sync"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const (
	defaultEnabled  = true
	defaultDisabled = false
)

var (
	initiatedCollectorsMtx = sync.Mutex{}
	InitiatedMetrics       = make(map[string]Metrics)
	MetricsState           = make(map[string]*bool)
	forcedMetrics          = map[string]bool{}
)

func registerCollector(collector string, isDefaultEnabled bool, metric Metrics) {
	var helpDefaultState string
	if isDefaultEnabled {
		helpDefaultState = "enabled"
	} else {
		helpDefaultState = "disabled"
	}

	flagName := fmt.Sprintf("collector.%s", collector)
	flagHelp := fmt.Sprintf("Enable the %s collector (default: %s).", collector, helpDefaultState)
	defaultValue := fmt.Sprintf("%v", true)

	flag := kingpin.Flag(flagName, flagHelp).Default(defaultValue).Action(collectorFlagAction(collector)).Bool()
	MetricsState[collector] = flag

	InitiatedMetrics[collector] = metric
}

func collectorFlagAction(collector string) func(ctx *kingpin.ParseContext) error {
	return func(ctx *kingpin.ParseContext) error {
		forcedMetrics[collector] = true
		return nil
	}
}

type Metrics interface {
	Update() error
}
