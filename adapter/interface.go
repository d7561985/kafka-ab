package adapter

type Producer interface {
	Emit(topic string, key []byte, body []byte) error
}
