package types

type (
	JolokiaResponse struct {
		Request   map[string]string
		Timestamp int64
		Status    int
		Value     map[string]JolokiaReponseValue
	}

	JolokiaReponseValue map[string]float64
)
