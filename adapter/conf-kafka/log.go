package conf_kafka

import (
	"log"
	"log/syslog"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

// KafkaLog read log if blocking operation
func KafkaLog(k <-chan kafka.LogEvent) {
	for event := range k {

		switch syslog.Priority(event.Level) {
		case syslog.LOG_EMERG, syslog.LOG_ALERT, syslog.LOG_CRIT, syslog.LOG_ERR, syslog.LOG_WARNING:
			log.Fatalf("%s: %s", syslogToString(event.Level), event.Message)
		default:
			log.Printf("%s: %s", syslogToString(event.Level), event.Message)
		}
	}
}

func syslogToString(l int) string {
	switch syslog.Priority(l) {
	case syslog.LOG_EMERG:
		return "LOG_EMERG"
	case syslog.LOG_ALERT:
		return "LOG_ALERT"
	case syslog.LOG_CRIT:
		return "LOG_CRIT"
	case syslog.LOG_ERR:
		return "LOG_ERR"
	case syslog.LOG_WARNING:
		return "LOG_WARNING"
	case syslog.LOG_NOTICE:
		return "LOG_NOTICE"
	case syslog.LOG_INFO:
		return "LOG_INFO"
	case syslog.LOG_DEBUG:
		return "LOG_DEBUG"
	default:
		return "LOG_UNKNOWN"
	}
}
