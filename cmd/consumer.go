package cmd

import (
	"github.com/urfave/cli/v2"
	conf_kafka "kafka-bench/adapter/conf-kafka"
	"kafka-bench/events"
	super_consumer "kafka-bench/usecase/super-consumer"
)

const (
	fTopic       = "topic"
	fKafkaServer = "kafka-server"
)

type consumerCMD struct{}

func ConsumerCMD() *cli.Command {
	x := consumerCMD{}
	return x.Command()
}

func (c *consumerCMD) Command() *cli.Command {
	return &cli.Command{
		Name:    "consumer",
		Aliases: []string{"c"},
		Action:  c.Action,
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
				Name:    fThreads,
				Value:   10,
				EnvVars: []string{"THREADS"},
			},
		},
	}
}

func (c *consumerCMD) Action(ctx *cli.Context) error {
	threads := ctx.Int(fThreads)
	var list []<-chan events.EventResponse

	for i := 0; i < threads; i++ {
		e := conf_kafka.NewConsumer(conf_kafka.Config{BootStrap: ctx.String(fKafkaServer)}).
			Subscribe(ctx.Context, ctx.String(fTopic))
		list = append(list, e)
	}

	sc := super_consumer.New()
	sc.Register(ctx.Context, list)

	<-ctx.Done()

	return nil
}
