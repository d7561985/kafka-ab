package cmd

import (
	"github.com/urfave/cli/v2"
	conf_kafka "kafka-bench/adapter/conf-kafka"
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
		Name:   "consumer",
		Aliases: []string{"c"},
		Action: c.Action,
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

		},
	}
}

func (c *consumerCMD) Action(ctx *cli.Context) error {
	cons := conf_kafka.NewConsumer(conf_kafka.Config{BootStrap: ctx.String(fKafkaServer)})
	e := cons.Subscribe(ctx.Context, ctx.String(fTopic))

	sc := super_consumer.New()
	go sc.Register(ctx.Context, e)

	<-ctx.Done()

	return nil
}
