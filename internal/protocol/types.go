package protocol

type APIKey int16

const (
	Produce APIKey = 0
	Fetch   APIKey = 1
	// joingroup, metadata, etc
)

type RequestHeader struct {
	APIKey        APIKey
	APIVersion    int16
	CorrelationID int32
	ClientID      string
}

type Message struct {
	Offset int64 //asign by broker
	Key    []byte
	Value  []byte
	// Timestamp
}

type ProduceRequest struct {
	Topic     string
	Partition int32
	Messages  []Message
}

type FetchRequest struct {
	Topic     string
	Partition int32
	Offset    int64
	MaxBytes  int32
}

type ProduceResponse struct {
	Topic     string
	Partition int32
	Offset    int64
	ErrorCode int16
}

type FetchResponse struct {
	Topic     string
	Partition int32
	Messages  []Message
	ErrorCode int16
}
