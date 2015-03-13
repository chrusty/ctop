package types

type (
	ColumnFamily struct {
		ReadCount                         int64
		WriteCount                        int64
		LiveDiskSpaceUsed                 int64
		MeanRowSize                       int64
		MaxRowSize                        int64
		RecentSSTablesPerReadHistogram    map[int]int
		RecentReadLatencyHistogramMicros  map[int]int
		RecentWriteLatencyHistogramMicros map[int]int
		EstimatedColumnCountHistogram     map[int]int
		EstimatedRowSizeHistogram         map[int]int
	}

	KeySpace struct {
		ColumnFamilies map[string]ColumnFamily
	}

	Cluster struct {
		Name      string
		KeySpaces map[string]KeySpace
	}
)
