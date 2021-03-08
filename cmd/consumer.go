package cmd

import (
	conf_kafka "kafka-bench/adapter/conf-kafka"
	"kafka-bench/events"
	super_consumer "kafka-bench/usecase/super-consumer"

	"github.com/urfave/cli/v2"
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
				EnvVars: []string{Topic},
			},
			&cli.StringFlag{
				Name:    fKafkaServer,
				Value:   "PLAINTEXT://127.0.0.1:9094",
				EnvVars: []string{KafkaBootstrap},
			},
			&cli.IntFlag{
				Name:    fConcurrency,
				Value:   10,
				EnvVars: []string{Concurrency},
			},
		},
	}
}

func (c *consumerCMD) Action(ctx *cli.Context) error {
	threads := ctx.Int(fConcurrency)
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
