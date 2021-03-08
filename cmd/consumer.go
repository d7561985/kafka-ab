package cmd

import (
	"context"
	conf_kafka "kafka-bench/adapter/conf-kafka"
	"kafka-bench/events"
	super_consumer "kafka-bench/usecase/super-consumer"
	"log"
	"time"

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
			&cli.DurationFlag{
				Name:        fTimeLimit,
				EnvVars:     []string{TimeLimit},
				Aliases:     []string{"t"},
				Value:       time.Minute * 10,
				DefaultText: `Seconds to max. to spend on benchmarking.`,
				Usage:       "Stop all tasks instantly. In case of desired request not reach will exist with status 1 ",
			},
			&cli.IntFlag{
				Name:        fRequests,
				Aliases:     []string{"n"},
				DefaultText: "Number of requests to perform/consume",
				Usage:       "require duration postfix: s - seconds, h - hours and etc",
				EnvVars:     []string{Requests},
				Value:       0,
			},
			&cli.IntFlag{
				Name:        fVerbosity,
				EnvVars:     []string{Verbosity},
				Aliases:     []string{"v"},
				Value:       0,
				DefaultText: "How much troubleshooting info to print",
			},
			&cli.StringFlag{
				Name:        fForceName,
				EnvVars:     []string{ForceName},
				Aliases:     []string{"f"},
				Value:       "",
				DefaultText: "all consumer's groups have the same group name which pass throughout",
			},
			&cli.BoolFlag{
				Name:        fStaticGroupName,
				EnvVars:     []string{Static},
				Aliases:     []string{"s"},
				Value:       false,
				DefaultText: "all consumer's groups hase unique but static name (group-1, group-2 and group-...",
			},
			&cli.BoolFlag{
				Name:        fAutoCommit,
				EnvVars:     []string{AutoCommit},
				Aliases:     []string{"a"},
				Value:       true,
				DefaultText: "auto-commit option",
			},
			&cli.BoolFlag{
				Name:        fEarliest,
				EnvVars:     []string{EARLIEST},
				Aliases:     []string{"e"},
				Value:       true,
				DefaultText: "read from the very beginning of log",
			},
		},
	}
}

func (c *consumerCMD) Action(root *cli.Context) error {
	threads := root.Int(fConcurrency)
	var list []<-chan events.EventResponse

	ctx, cancel := context.WithTimeout(root.Context, root.Duration(fTimeLimit))
	defer cancel()

	for i := 0; i < threads; i++ {
		e := conf_kafka.NewConsumer(i, conf_kafka.Config{
			BootStrap:       root.String(fKafkaServer),
			Verbosity:       root.Int(fVerbosity),
			ForceName:       root.String(fForceName),
			StaticGroupName: root.Bool(fStaticGroupName),
			AutoCommit:      root.Bool(fAutoCommit),
			Earliest:        root.Bool(fEarliest),
		}).Subscribe(ctx, root.String(fTopic))
		list = append(list, e)
	}

	requests := root.Uint64(fRequests)
	sc := super_consumer.New(super_consumer.Config{
		Verbosity: root.Int(fVerbosity),
		Requests:  requests,
	})

	done := sc.Start(ctx, list)

	select {
	case <-ctx.Done(): // timeout
		if requests > 0 {
			// not meat target
			log.Fatalf("timeout reached")
		} else {
			log.Println("operation completed without issues")
		}
	case <-root.Done(): // signal ok
		log.Println("termination request")
	case <-done:
		log.Println("operation completed without issues")
	}

	return nil
}
