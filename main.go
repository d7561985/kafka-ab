package main

import (
	"log"

	"github.com/d7561985/kafka-ab/cmd"
)

func main() {
	if err := cmd.Start(); err != nil {
		log.Print(err)
	}
}
