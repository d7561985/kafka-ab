package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/urfave/cli/v2"
)

func Start() error {
	app := &cli.App{
		Name:  "Kafka-ab",
		Usage: "Helps to check slo for kafka as consumer or producer with different requirement to topic kind and producer options.",
		Commands: []*cli.Command{
			ConsumerCMD(),
			ProducerCMD(),
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	go signalListener(cancel)

	err := app.RunContext(ctx, os.Args)
	if err != nil {
		return err
	}

	return nil
}

func signalListener(cancel func()) {
	defer cancel()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)

	<-ch
}
