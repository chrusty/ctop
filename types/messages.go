package types

// import "time"

type (
	LogMessage struct {
		Severity string
		Message  string
	}

	CFStats struct {
		KeySpace          string
		ColumnFamily      string
		ReadLatency       float64
		ReadRate          float64
		WriteLatency      float64
		WriteRate         float64
		LiveDiskSpaceUsed float64
		MeanRowSize       int64
		MaxRowSize        int64
	}

	CFMetric struct {
		KeySpace         string
		ColumnFamily     string
		MetricName       string
		MetricFloatValue float64
		MetricTimeStamp  int64
	}
)
