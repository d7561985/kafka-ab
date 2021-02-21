package super_producer

import "time"

type Config struct {
	Topic       string
	Concurrency int
	WindowSize  int
	Requests    int
	TimeOut     time.Duration

	Verbosity int
}
