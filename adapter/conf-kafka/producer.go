package conf_kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"kafka-bench/adapter"
	"log"
)

type producer struct {
	Config
	srv *kafka.Producer
}

func NewProducer(c Config) adapter.Producer {
	conf := kafka.ConfigMap{
		"bootstrap.servers":   c.BootStrap,
		"delivery.timeout.ms": int(c.TimeOut.Milliseconds()),
	}

	p, err := kafka.NewProducer(&conf)
	if err != nil {
		log.Fatalf("kafka: create producer error: %s", err)
	}

	return &producer{srv: p, Config: c}
}

func (p *producer) Emit(topic string, key []byte, body []byte) error {
	c := make(chan kafka.Event)
	defer close(c)

	err := p.srv.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
		Value:          body,
	}, c)

	if err != nil {
		return fmt.Errorf("producer: %s key %s  erorr: %w", topic, string(key), err)
	}

	v := <-c

	switch ev := v.(type) {
	case *kafka.Message:
		if ev.TopicPartition.Error != nil {
			return fmt.Errorf("producer: %s user %s failed to deliver: %v", topic, string(key), ev.TopicPartition)
		}
	default:
		log.Printf("kafka producer unhandled response: %T => %[1]v", v)
	}

	if p.Verbosity >= 2 {
		log.Printf("send topoc: %s", topic)
	}

	return nil
}
