package kafka

import "time"

const (
	minBytes               = 10e3
	maxBytes               = 10e6
	queueCapacity          = 100
	heartbeatInterval      = 1 * time.Second
	commitInterval         = 0
	partitionWatchInterval = 1 * time.Second
	maxAttempts            = 10
	dialTimeout            = 3 * time.Minute
	maxWait                = 1 * time.Second

	writerReadTimeout  = 1 * time.Second
	writerWriteTimeout = 1 * time.Second

	deadLetterQueueTopic = "dead-letter-queue"
)
