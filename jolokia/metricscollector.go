package jolokia

import (
	"github.com/hailocab/ctop/types"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	// "strconv"
	// "strings"
	"time"
)

var (
	// Define a list of metrics to collect (sticking to 5 important ones to save on HTTP calls):
	cfMetrics = []string{"ReadTotalLatency/Count", "WriteTotalLatency/Count", "LiveDiskSpaceUsed/Count", "ReadLatency/Mean", "WriteLatency/Mean"}
	metricsChannel chan types.CFMetric
	messageChannel chan types.LogMessage
	// "ReadTotalLatency"  - The number of reads to a CF (Count)
	// "WriteTotalLatency" - The number of writes to a CF (Count)
	// "LiveDiskSpaceUsed" - Disk space used Count)
	// "MeanRowSize"       - Mean row-size
	// "MaxRowSize"        - Max row-size
	// "ReadLatency"       - Read latency (Mean)
	// "WriteLatency"      - Write latency (Mean)
)

// Checks the connection to the metrics provider:
func CheckConnection(metricsURL string) error {
	// Request the root URL:
	URL := fmt.Sprintf(metricsURL)

	_, err := http.Get(URL)
	return err
}


// Retreive metrics from MX4J:
func getCFMetrics(metricsURL string) (error) {

	logToChannel("debug", fmt.Sprintf("Getting metrics from (%s)", metricsURL))

	// Get the CFMetrics:
	for _, cfMetric := range cfMetrics {

		// Create a new JolokiaResponse{} to unmarshal the JSON into:
		jolokiaResponse := types.JolokiaResponse{}

		logToChannel("info", fmt.Sprintf("Getting %s metrics ...", cfMetric))

		// Build the reqest URL:
		URL := fmt.Sprintf("%s/read/org.apache.cassandra.metrics:type=ColumnFamily,keyspace=*,scope=*,name=%s", metricsURL, cfMetric)

		// Request the data from MX4J:
		httpResponse, err := http.Get(URL)
		if err != nil {
			logToChannel("error", fmt.Sprintf("Trouble talking to Jolokia (%s)\n%s", URL, err))
			continue
		} else {
			logToChannel("debug", fmt.Sprintf("Got HTTP response code (%d)", httpResponse.StatusCode))
		}

		// Read the response:
		jsonResponse, err := ioutil.ReadAll(httpResponse.Body)
		if err != nil {
			logToChannel("error", fmt.Sprintf("Couldn't get response body!\n%s", err))
			continue
		}

		// UnMarshal the JSON response:
		err = json.Unmarshal([]byte(jsonResponse), &jolokiaResponse)
		if err != nil {
			logToChannel("error", fmt.Sprintf("Couldn't unmarshal the response!\n%s", err))
			continue
		} else {
			logToChannel("debug", fmt.Sprintf("Got a metric - GREAT SUCCESS!"))

			// Process all of the returned values:
			for jolokiaResponseKey, jolokiaResponseValue := range jolokiaResponse.Value {

				// // Split up the comma-delimited metadata string:
				// columnFamilyMetaData := strings.Split(columnFamilyList.CFList[i].ColmnFamily, ",")

				// // Now split these values up by "=" to get the metadata we're after:
				// keySpaceName := strings.Split(columnFamilyMetaData[1], "=")

				// // Create a new KeySpace{}:
				// cluster.KeySpaces[keySpaceName[1]] = types.KeySpace{
				// 	ColumnFamilies: make(map[string]types.ColumnFamily),
				// }

				// // Make a new Metric struct:
				// cfMetric := types.CFMetric{
				// 	KeySpace:         name,
				// 	ColumnFamily:     columnFamily,
				// 	MetricName:       metric.CFLongData.Name,
				// 	MetricIntValue:   metricIntValue,
				// 	MetricFloatValue: metricFloatValue,
				// 	MetricTimeStamp:  time.Now().Unix(),
				// }

				// // Put it in the metrics channel:
				// select {
				// case metricsChannel <- cfMetric:
				// 	logToChannel("debug", fmt.Sprintf("Sent a metric"))
				// default:
				// 	logToChannel("info", fmt.Sprintf("Couldn't send metric!"))
				// }

				logToChannel("debug", fmt.Sprintf("%s => %s", jolokiaResponseKey, jolokiaResponseValue))
			}
		}
	}

	return nil
}

// Collects actual metrics
func MetricsCollector(metricsChan chan types.CFMetric, messageChan chan types.LogMessage, metricsURL string) {

	messageChannel = messageChan
	metricsChannel = metricsChan

	for {
		// Get metrics for each ColumnFamily from MX4J:
		err := getCFMetrics(metricsURL)
		if err != nil {
			logToChannel("error", fmt.Sprintf("Couldn't get metrics!\n%s", err))
		}
		time.Sleep(5 * time.Second)
	}

}

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
		fmt.Printf("[%s] %s\n", severity, message)

	default:

	}
}