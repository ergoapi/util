package exsink

type SinkFactory struct {
}

type EventSink interface {
	SendEvent(any) error
}

func NewSinkFactory() *SinkFactory {
	return &SinkFactory{}
}
