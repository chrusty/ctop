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
		ReadCount         int64
		ReadCountTS       int64
		ReadLatency       float64
		ReadRate          float64
		WriteCount        int64
		WriteCountTS      int64
		WriteLatency      float64
		WriteRate         float64
		LiveDiskSpaceUsed int64
		MeanRowSize       int64
		MaxRowSize        int64
	}

	CFMetric struct {
		KeySpace         string
		ColumnFamily     string
		MetricName       string
		MetricIntValue   int64
		MetricFloatValue float64
		MetricTimeStamp  int64
	}
)
