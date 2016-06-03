package main

import (
	"fmt"
	"github.com/hailocab/ctop/types"
	"github.com/nsf/termbox-go"
)

// Logging to a channel (from anywhere):
func logToChannel(severity string, message string) {
	// Make a new LogMessage struct:
	logMessage := types.LogMessage{
		Severity: severity,
		Message:  message,
	}

	// Put it in the messages channel:
	select {
	case messageChannel <- logMessage:

	default:

	}
}

// Takes metrics off the channel and adds them up:
func handleMetrics() {

	var cfStats types.CFStats

	for {
		// Get a metric from the channel:
		cfMetric := <-metricsChannel
		logToChannel("debug", fmt.Sprintf("Received a metric! %s", cfMetric.MetricName))

		// Build the key:
		statName := cfMetric.KeySpace + ":" + cfMetric.ColumnFamily

		statsMutex.Lock()
		defer statsMutex.Unlock()

		// See if we already have a stats-entry:
		if _, ok := stats[statName]; ok {
			// Use the existing stats-entry:
			logToChannel("debug", fmt.Sprintf("Updating existing stat (%s)", statName))
			cfStats = stats[statName]
		} else {
			// Add a new entry to the map:
			logToChannel("debug", fmt.Sprintf("Adding new stat (%s)", statName))
			cfStats = types.CFStats{
				ReadLatency:  0.0,
				ReadRate:     0.0,
				WriteLatency: 0.0,
				WriteRate:    0.0,
				KeySpace:     cfMetric.KeySpace,
				ColumnFamily: cfMetric.ColumnFamily,
			}
		}

		// Figure out which metric we need to update:
		if cfMetric.MetricName == "ReadLatency/OneMinuteRate" {
			// Read rate(s):
			cfStats.ReadRate = cfMetric.MetricFloatValue
			stats[statName] = cfStats

		} else if cfMetric.MetricName == "WriteLatency/OneMinuteRate" {
			// Write rate(s):
			cfStats.WriteRate = cfMetric.MetricFloatValue
			stats[statName] = cfStats

		} else if cfMetric.MetricName == "LiveDiskSpaceUsed/Count" {
			// Total disk space used(k):
			cfStats.LiveDiskSpaceUsed = cfMetric.MetricFloatValue
			stats[statName] = cfStats

		} else if cfMetric.MetricName == "ReadLatency/Mean" {
			// ReadLatency (MicroSeconds):
			if cfMetric.MetricFloatValue > 0 {
				cfStats.ReadLatency = cfMetric.MetricFloatValue / 1000
				stats[statName] = cfStats
			}

		} else if cfMetric.MetricName == "WriteLatency/Mean" {
			// WriteLatency (MicroSeconds):
			if cfMetric.MetricFloatValue > 0 {
				cfStats.WriteLatency = cfMetric.MetricFloatValue / 1000
				stats[statName] = cfStats
			}
		}

		statsMutex.Unlock()

	}

}

// Returns the key-code:
func handleKeypress(ev *termbox.Event) {
	logToChannel("debug", fmt.Sprintf("Key pressed: %s", ev.Ch))
}
