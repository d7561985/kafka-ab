package conf_kafka

import (
	"context"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/icrowley/fake"
	"kafka-bench/events"
	"log"
	"os"
)

type consumer struct {
	srv *kafka.Consumer

	group string
}

func NewConsumer(c Config) *consumer {
	host, err := os.Hostname()
	if err != nil {
		host = "common"
	}

	group := fake.FullName()

	conf := kafka.ConfigMap{
		"bootstrap.servers":               c.BootStrap,
		"session.timeout.ms":              6000,
		"group.id":                        group,
		"group.instance.id":               host,
		"auto.offset.reset":               "earliest", //earliest/latest
		"enable.auto.commit":              true,
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true, // app rebalance
		"enable.partition.eof":            true, // Enable generation of PartitionEOF when the end of a partition is reached.
	}

	p, err := kafka.NewConsumer(&conf)
	if err != nil {
		log.Fatalf("kafka:consumer start error: %s", err)
	}

	return &consumer{srv: p, group: group}
}

func (s *consumer) Subscribe(ctx context.Context, topic string) <-chan events.EventResponse {
	err := s.srv.Subscribe(topic, nil)
	if err != nil {
		log.Fatalf("kafka:consumer subscribe error: %s", err)
	}

	log.Printf("consumer[group: %s] listen: %v",
		s.group, topic)

	// first part of getting all partition offsets
	//metadata, err := s.p.GetMetadata(&topic, false, 10000)
	//if err == nil {
	//    for _, topicMetadata := range metadata.Topics {
	//        for _, partition := range topicMetadata.Partitions {
	//            //fmt.Println("A",i, partition)
	//
	//            low, high, err := s.p.QueryWatermarkOffsets(topic, partition.ID, 10000)
	//            if err == nil {
	//                fmt.Printf("%q Partition[%d][%d,%d]\n", topic, partition.ID, low, high)
	//            }
	//        }
	//    }
	//} else {
	//    panic(fmt.Errorf("metadata error: %w, %[1]T", err))
	//    return nil
	//}

	res := make(chan events.EventResponse)

	go func() {
		for {
			select {
			case <-ctx.Done():
				_ = s.srv.Close()
				return
			case m := <-s.srv.Events():
				switch e := m.(type) {
				case kafka.AssignedPartitions:
					//log.Printf("AssignedPartitions: %s", m.String())

					if err := s.srv.Assign(e.Partitions); err != nil {
						log.Printf("a err: %s", err)
					}
				case kafka.RevokedPartitions:
					//log.Printf("RevokedPartitions")

					if err := s.srv.Unassign(); err != nil {
						log.Printf("unsign error:%s", err)
					}
				case kafka.Error:
					// Errors should generally be considered as informational, the client will try to automatically recover
					log.Printf("KafkaConsumer: Error: %v: %v\n", e.Code(), e)
				case *kafka.Message:
					res <- events.EventResponse{
						Key:   string(e.Key),
						Value: e.Value,
					}
				default:
					//fmt.Println("unhandled event", m)
				}
			default:
				//fmt.Printf("%s\n", m.TopicPartition.String())
			}
		}
	}()

	return res
}
