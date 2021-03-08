package conf_kafka

import (
	"context"
	"fmt"
	"kafka-bench/events"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/icrowley/fake"
)

type consumer struct {
	srv *kafka.Consumer

	group string
}

func NewConsumer(i int, c Config) *consumer {
	var instanceID string

	host, err := os.Hostname()
	if err != nil {
		host = "common"
	}

	group := fake.FullName()
	if len(c.ForceName) > 0 {
		group = c.ForceName
		instanceID = fake.FullName()
	} else {
		instanceID = host
	}

	if c.StaticGroupName {
		group = fmt.Sprintf("group#%d", i)
	}

	var offset = "earliest"
	if !c.Earliest {
		offset = "latest"
	}

	conf := kafka.ConfigMap{
		"bootstrap.servers":               c.BootStrap,
		"session.timeout.ms":              6000,
		"group.id":                        group,
		"group.instance.id":               instanceID,
		"auto.offset.reset":               offset,
		"enable.auto.commit":              c.AutoCommit,
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true, // app rebalance
		"enable.partition.eof":            true, // Enable generation of PartitionEOF when the end of a partition is reached.
		"go.logs.channel.enable":          true, // handle kafka logs
	}

	p, err := kafka.NewConsumer(&conf)
	if err != nil {
		log.Fatalf("kafka:consumer start error: %s", err)
	}

	go KafkaLog(p.Logs())

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
					log.Printf("kafka.consumer AssignedPartitions: %s", m.String())

					if err := s.srv.Assign(e.Partitions); err != nil {
						log.Printf("kafka.consumer assign partition err: %s", err)
					}
				case kafka.RevokedPartitions:
					log.Printf("kafka.consumer RevokedPartitions: %s", m.String())

					if err := s.srv.Unassign(); err != nil {
						log.Printf("kafka.consumer unsign partition error:%s", err)
					}
				case kafka.Error:
					if checkFatalKafka(e.Code()) {
						log.Fatalf("kafka.consumer fatal error: %s", e.String())
					} else {
						// Errors should generally be considered as informational, the client will try to automatically recover
						log.Printf("kafka.consumer error: %v: %s", e.Code(), e.String())
					}

				case *kafka.Message:
					res <- events.EventResponse{
						Key:   string(e.Key),
						Value: e.Value,
					}
				default:
					//fmt.Println("unhandled event", m)
				}
			}
		}
	}()

	return res
}

func checkFatalKafka(c kafka.ErrorCode) bool {
	return c == kafka.ErrDestroy || c == kafka.ErrFail || c == kafka.ErrTransport || c == kafka.ErrCritSysResource ||
		c == kafka.ErrResolve || c == kafka.ErrAllBrokersDown || c == kafka.ErrRetry || c == kafka.ErrFatal ||
		c == kafka.ErrOffsetOutOfRange || c == kafka.ErrBrokerNotAvailable || c == kafka.ErrNetworkException ||
		c == kafka.ErrOffsetNotAvailable || c == kafka.ErrPreferredLeaderNotAvailable || c == kafka.ErrGroupMaxSizeReached
}
