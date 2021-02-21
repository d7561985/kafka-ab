package main

import (
	"kafka-bench/cmd"
	"log"
)

func main() {
	if err := cmd.Start(); err != nil {
		log.Print(err)
	}
}
