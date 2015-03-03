package types

type (
	// This is a response from a query for a CFHistogram array:
	// ("http://%s:8081/getattribute?objectname=org.apache.cassandra.db:type=ColumnFamilies,keyspace=%s,columnfamily=%s&attribute=%s&format=array&template=viewarray&template=identity", cassandraIP, name, columnFamily, cfHistograms[i])
	MX4JCassandraCFHistogram struct {
		//XMLName xml.Name `xml:"MBean"`
		CFHistogram []MX4JCassandraCFHistogramElement `xml:"Attribute>Array>Element"`
	}

	// This is one of the array elements:
	MX4JCassandraCFHistogramElement struct {
		Index string `xml:"index,attr"`
		Value string `xml:"element,attr"`
	}

	// This is a response from a query for an individual bit of data:
	// ("http://%s:8081/getattribute?objectname=org.apache.cassandra.db:type=ColumnFamilies,keyspace=%s,columnfamily=%s&attribute=%s&format=long&template=identity", cassandraIP, name, columnFamily, cfMetrics[i])
	MX4JCassandraCFLongData struct {
		//XMLName xml.Name `xml:"MBean"`
		CFLongData MX4JCassandraCFLongDataAttribute `xml:"Attribute"`
	}

	// This is the bit of data itself:
	MX4JCassandraCFLongDataAttribute struct {
		Name  string `xml:"name,attr"`
		Value string `xml:"value,attr"`
	}

	// This is the response from a query for the list of ColumnFamilies:
	// ("http://%s:8081/server?instanceof=org.apache.cassandra.db.ColumnFamilyStore&template=identity", cassandraIP)
	MX4JCFList struct {
		CFList []MX4JCFListColumnFamily `xml:"MBean"`
	}

	MX4JCFListColumnFamily struct {
		ColmnFamily string `xml:"objectname,attr"`
	}
)
