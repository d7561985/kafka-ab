package conf_kafka

import (
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"kafka-bench/adapter"
	"log"
)

type producer struct {
	srv *kafka.Producer
}

func NewProducer(c Config) adapter.Producer {
	conf := kafka.ConfigMap{"bootstrap.servers": c.BootStrap}

	p, err := kafka.NewProducer(&conf)
	if err != nil {
		log.Fatalf("kafka: create producer error: %s", err)
	}

	return &producer{srv: p}
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
	}

	//log.Printf("send topoc: %s", topic)
	return nil
}
