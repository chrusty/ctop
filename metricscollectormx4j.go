package main

import (
	types "./types"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	MX4JPort = "8081"
)

var (
	// Define a list of metrics to collect (sticking to 5 important ones to save on HTTP calls):
	cfMetrics = []string{"ReadCount", "WriteCount", "LiveDiskSpaceUsed", "RecentReadLatencyMicros", "RecentWriteLatencyMicros"}

	// "ReadCount"                - The number of reads to a CF
	// "WriteCount"               - The number of writes to a CF
	// "LiveDiskSpaceUsed"        - Disk space used
	// "MeanRowSize"              - Mean row-size
	// "MaxRowSize"               - Max row-size
	// "RecentReadLatencyMicros"  - Read latency
	// "RecentWriteLatencyMicros" - Write latency
)

// Checks the connection to MX4J:
func checkConnection(cassandraIP string) error {
	// Request the root URL:
	URL := fmt.Sprintf("http://%s:%s/", cassandraHost, MX4JPort)

	_, err := http.Get(URL)
	return err
}

// Return a list of keySpaces and columnFamilies from MX4J:
func getCluster(cassandraIP string) (types.Cluster, error) {

	logToChannel("info", fmt.Sprintf("Getting list of KeySpaces and ColumnFamilies from (%s:%s)", cassandraIP, MX4JPort))

	// Create a new MX4JCFList{} to unmarshal the XML into:
	columnFamilyList := types.MX4JCFList{}

	// Build the reqest URL:
	URL := fmt.Sprintf("http://%s:%s/server?instanceof=org.apache.cassandra.db.ColumnFamilyStore&template=identity", cassandraIP, MX4JPort)

	// Request the data from MX4J:
	httpResponse, err := http.Get(URL)
	if err != nil {
		logToChannel("error", fmt.Sprintf("Trouble talking to MX4J (%s)\n%s", URL, err))
	} else {
		logToChannel("debug", fmt.Sprintf("Got HTTP response code (%d)", httpResponse.StatusCode))
	}

	// Read the response:
	xmlResponse, err := ioutil.ReadAll(httpResponse.Body)
	if err != nil {
		logToChannel("error", fmt.Sprintf("Couldn't get response body!\n%s", err))
	}

	// UnMarshal the XML response:
	err = xml.Unmarshal([]byte(xmlResponse), &columnFamilyList)
	if err != nil {
		logToChannel("error", fmt.Sprintf("Couldn't unmarshal the response!\n%s", err))
	} else {
		logToChannel("debug", fmt.Sprintf("Got a ColumnFamily list - great success!"))
		//log.Debugf("- %s", columnFamilyList)
	}

	// Create a new types.Cluster{}:
	cluster := types.Cluster{
		Name:      "cruft",
		KeySpaces: make(map[string]types.KeySpace),
	}

	// Populate the Cluster{} with the results returned from MX4J:
	for i := range columnFamilyList.CFList {
		// Split up the comma-delimited metadata string:
		columnFamilyMetaData := strings.Split(columnFamilyList.CFList[i].ColmnFamily, ",")

		// Now split these values up by "=" to get the metadata we're after:
		keySpaceName := strings.Split(columnFamilyMetaData[1], "=")

		// Create a new KeySpace{}:
		cluster.KeySpaces[keySpaceName[1]] = types.KeySpace{
			ColumnFamilies: make(map[string]types.ColumnFamily),
		}
	}

	for i := range columnFamilyList.CFList {
		// Split up the comma-delimited metadata string:
		columnFamilyMetaData := strings.Split(columnFamilyList.CFList[i].ColmnFamily, ",")

		// Now split these values up by "=" to get the metadata we're after:
		columnFamilyType := strings.Split(columnFamilyMetaData[0], "=")
		keySpaceName := strings.Split(columnFamilyMetaData[1], "=")
		columnFamilyName := strings.Split(columnFamilyMetaData[2], "=")

		// Create a new ColumnFamily{}:
		if columnFamilyType[1] == "ColumnFamilies" {
			logToChannel("debug", fmt.Sprintf("Found KS:CF - %s:%s (%s)", keySpaceName[1], columnFamilyName[1], columnFamilyType[1]))
			cluster.KeySpaces[keySpaceName[1]].ColumnFamilies[columnFamilyName[1]] = types.ColumnFamily{}
		}
	}

	return cluster, nil
}

// Retreive metrics from MX4J:
func getCFMetrics(cluster types.Cluster, cassandraIP string) (types.Cluster, error) {

	logToChannel("debug", fmt.Sprintf("Getting metrics from (%s:%s)", cassandraIP, MX4JPort))

	// Iterate through our Cluster{}:
	for name, keySpace := range cluster.KeySpaces {
		for columnFamily := range keySpace.ColumnFamilies {

			// Get the CFMetrics:
			for i := range cfMetrics {

				// Create a new MX4JCassandraCFLongData{} to unmarshal the XML into:
				metric := types.MX4JCassandraCFLongData{}

				logToChannel("info", fmt.Sprintf("Getting %s:%s:%s", name, columnFamily, cfMetrics[i]))

				// Build the reqest URL:
				URL := fmt.Sprintf("http://%s:%s/getattribute?objectname=org.apache.cassandra.db:type=ColumnFamilies,keyspace=%s,columnfamily=%s&attribute=%s&format=long&template=identity", cassandraIP, MX4JPort, name, columnFamily, cfMetrics[i])

				// Request the data from MX4J:
				httpResponse, err := http.Get(URL)
				if err != nil {
					logToChannel("error", fmt.Sprintf("Trouble talking to MX4J (%s)\n%s", URL, err))
				} else {
					logToChannel("debug", fmt.Sprintf("Got HTTP response code (%d)", httpResponse.StatusCode))
				}

				// Read the response:
				xmlResponse, err := ioutil.ReadAll(httpResponse.Body)
				if err != nil {
					logToChannel("error", fmt.Sprintf("Couldn't get response body!\n%s", err))
				}

				// UnMarshal the XML response:
				err = xml.Unmarshal([]byte(xmlResponse), &metric)
				if err != nil {
					logToChannel("error", fmt.Sprintf("Couldn't unmarshal the response!\n%s", err))
				} else {
					logToChannel("debug", fmt.Sprintf("Got a metric - GREAT SUCCESS!"))

					// Make an int64:
					metricIntValue, _ := strconv.ParseInt(metric.CFLongData.Value, 0, 64)
					metricFloatValue, _ := strconv.ParseFloat(metric.CFLongData.Value, 64)

					// Make a new Metric struct:
					cfMetric := types.CFMetric{
						KeySpace:         name,
						ColumnFamily:     columnFamily,
						MetricName:       metric.CFLongData.Name,
						MetricIntValue:   metricIntValue,
						MetricFloatValue: metricFloatValue,
						MetricTimeStamp:  time.Now().Unix(),
					}

					// Put it in the metrics channel:
					select {
					case metricsChannel <- cfMetric:
						logToChannel("debug", fmt.Sprintf("Sent a metric."))
					default:
						logToChannel("info", fmt.Sprintf("Couldn't send metric!"))
					}
				}
			}
		}
	}

	return cluster, nil
}

func MetricsCollector(cassandraHost string) {

	// Get a list of cluster KeySpaces and ColumnFamilies from MX4J:
	cluster, err := getCluster(cassandraHost)
	if err != nil {
		logToChannel("error", fmt.Sprintf("Couldn't get cluster schema!\n%s", err))
	}

	for {
		// Get metrics for each ColumnFamily from MX4J:
		cluster, err = getCFMetrics(cluster, cassandraHost)
		if err != nil {
			logToChannel("error", fmt.Sprintf("Couldn't get metrics!\n%s", err))
		}
		time.Sleep(5 * time.Second)
	}

}
