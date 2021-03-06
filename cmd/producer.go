package cmd

import (
	"context"
	"log"
	"time"

	conf_kafka "github.com/d7561985/kafka-ab/adapter/conf-kafka"
	super_producer "github.com/d7561985/kafka-ab/usecase/super-producer"

	"github.com/urfave/cli/v2"
)

type producerCMD struct{}

func ProducerCMD() *cli.Command {
	x := producerCMD{}
	return x.Command()
}

func (p *producerCMD) Command() *cli.Command {
	return &cli.Command{
		Name:    "producer",
		Aliases: []string{"p"},
		Action:  p.Action,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    fTopic,
				Value:   "my-topic",
				EnvVars: []string{Topic},
			},
			&cli.StringFlag{
				Name:    fKafkaServer,
				Value:   "PLAINTEXT://127.0.0.1:9094",
				Aliases: []string{"srv"},
				EnvVars: []string{KafkaBootstrap},
			},
			&cli.IntFlag{
				Name:    fWindowSize,
				Aliases: []string{"b"},
				Usage:   "Size of message send/receive buffer, in bytes",
				Value:   1 << 10,
				EnvVars: []string{WindowSize},
			},
			&cli.IntFlag{
				Name:        fConcurrency,
				Aliases:     []string{"c"},
				DefaultText: "Number of multiple requests to make/read at a time",
				Value:       100,
				EnvVars:     []string{Concurrency},
			},
			&cli.IntFlag{
				Name:        fRequests,
				Aliases:     []string{"n"},
				DefaultText: "Number of requests to perform/consume",
				EnvVars:     []string{Requests},
				Value:       0,
			},
			&cli.DurationFlag{
				Name:        fTimeLimit,
				EnvVars:     []string{TimeLimit},
				Aliases:     []string{"t"},
				Value:       time.Minute * 10,
				DefaultText: `Seconds to max. to spend on benchmarking.`,
				Usage:       "Stop all tasks instantly. In case of desired request not reach will exist with status 1 ",
			},
			&cli.DurationFlag{
				Name:        fTimeOut,
				EnvVars:     []string{TimeOut},
				Aliases:     []string{"s"},
				Value:       time.Second * 30,
				Usage:       "require duration postfix: s - seconds, h - hours and etc",
				DefaultText: "Seconds to max. wait for each response",
			},
			&cli.IntFlag{
				Name:        fVerbosity,
				EnvVars:     []string{Verbosity},
				Aliases:     []string{"v"},
				Value:       0,
				DefaultText: "How much troubleshooting info to print",
			},
		},
	}
}

func (p *producerCMD) Action(c *cli.Context) error {
	e := conf_kafka.NewProducer(conf_kafka.Config{
		BootStrap: c.String(fKafkaServer),
		TimeOut:   c.Duration(fTimeOut),
		Verbosity: c.Int(fVerbosity),
	})

	requests := c.Int(fRequests)

	sp := super_producer.New(e, super_producer.Config{
		Topic:       c.String(fTopic),
		Concurrency: c.Int(fConcurrency),
		Requests:    requests,
		TimeOut:     c.Duration(fTimeOut),
		Verbosity:   c.Int(fVerbosity),
		WindowSize:  c.Int(fWindowSize),
	})

	ctx, cancel := context.WithTimeout(c.Context, c.Duration(fTimeLimit))
	defer cancel()

	defer func(start time.Time) {
		log.Printf("duration: %s", time.Since(start))
	}(time.Now())

	go sp.Run(ctx)

	select {
	case <-ctx.Done(): // timeout
		if requests > 0 {
			// not meat target
			log.Fatalf("timeout reached")
		} else {
			log.Println("operation completed without issues")
		}
	case <-c.Done(): // signal ok
		log.Println("termination request")
	case <-sp.Done(): // app finished
		log.Println("operation completed without issues")
	}

	return nil
}
