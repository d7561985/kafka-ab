package cmd

import (
	"github.com/urfave/cli/v2"
	"os"
)

func Start() error {
	app := &cli.App{
		Name:  "Kafka bench",
		Usage: "Produce huge data and try it consume",
		Commands: []*cli.Command{
			ConsumerCMD(),
			ProducerCMD(),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		return err
	}

	return nil
}
