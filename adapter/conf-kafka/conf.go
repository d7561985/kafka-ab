package conf_kafka

import "time"

type Config struct {
	BootStrap string
	TimeOut   time.Duration
	Verbosity int
}
