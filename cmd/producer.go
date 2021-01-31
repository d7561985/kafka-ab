package cmd

import (
	"github.com/urfave/cli/v2"
	conf_kafka "kafka-bench/adapter/conf-kafka"
	super_producer "kafka-bench/usecase/super-producer"
)

const(
	fThreads = "threads"
	fMsgSize = "msgSize"
)
type producerCMD struct{}

func ProducerCMD() *cli.Command {
	x := producerCMD{}
	return x.Command()
}

func (p *producerCMD)Command() *cli.Command  {
	return &cli.Command{
		Name: "producer",
		Aliases: []string{"p"},
		Action: p.Action,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    fTopic,
				Value:   "my-topic",
				EnvVars: []string{"TOPIC"},
			},
			&cli.StringFlag{
				Name:    fKafkaServer,
				Value:   "127.0.0.1:9094",
				EnvVars: []string{"KAFKA_SERVER"},
			},
			&cli.IntFlag{
				Name:    fMsgSize,
				Usage:   "size of msg in bytes",
				Value:   1 << 10,
				EnvVars: []string{"MSG_SIZE_IN_BYTES"},
			},
			&cli.IntFlag{
				Name:    fThreads,
				Value:   100,
				EnvVars: []string{"THREADS"},
			},
		},
	}
}

func (p *producerCMD)Action(ctx *cli.Context) error  {
	e := conf_kafka.NewProducer(conf_kafka.Config{
		BootStrap: ctx.String(fKafkaServer),
	})

	sp := super_producer.New(e, ctx.String(fTopic), ctx.Int(fThreads), ctx.Int(fMsgSize))
	sp.Run(ctx.Context)

	<-ctx.Done()

	return nil
}
