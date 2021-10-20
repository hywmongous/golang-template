package kafka

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/hywmongous/example-service/pkg/es"
	"github.com/segmentio/kafka-go"
)

type Stream struct {
	// We only use a single topic per "bounded context"
	// For this reason it fits the "topic" to be a field
	// within the kafka stream struct instance
	// https://www.confluent.io/blog/put-several-event-types-kafka-topic/
	topic es.Topic
}

const (
	// Read/write shared kafka conf.
	defaultMaxAttempts   = 8
	defaultQueueCapacity = 100

	// Read kafka conf.
	defaultPartition              = 0
	defaultMinBytes               = 1
	defaultMaxBytes               = 1 << 20 // 1MB
	defaultMaxWait                = 10 * time.Second
	defaultReadLagInterval        = 1 * time.Second
	defaultHeartbeatInterval      = 3 * time.Second
	defaultCommitInterval         = 0 // Synchronous commits
	defaultPartitionWatchInterval = 5 * time.Second
	defaultWatchPartitionChanges  = false
	defaultSessionTimeout         = 30 * time.Second
	defaultRebalanceTimeout       = 30 * time.Second
	defaultJoinGroupBackoff       = 30 * time.Second
	defaultRetentionTime          = 30 * time.Hour
	defaultStartOffset            = kafka.FirstOffset
	defaultReadBackoffMin         = 100 * time.Millisecond
	defaultReadBackoffMax         = 1 * time.Second
	defaultIsolationLevel         = kafka.ReadCommitted

	// Write kafka conf.
	defaultBatchSize         = 100
	defaultBatchBytes        = 1 << 20
	defaultBatchTimeout      = 1 * time.Second
	defaultReadTimeout       = 10 * time.Second
	defaultWriteTimeout      = 10 * time.Second
	defaultRebalanceInterval = 10 * time.Second // Deprecated
	defaultIdleConnTimeout   = 10 * time.Second // Deprecated
	defaultRequiredAcks      = -1               // Wait for all replicas
	defaultAsync             = false            // By using false errors are not ignored
	broker                   = "ia_kafka:9092"
	group                    = "ia"

	defaultOffset        = 0
	defaultHighWaterMark = 0
)

var (
	ErrInvalidKafkaReaderConfig = errors.New("kafka reader config is invalid")
	ErrInvalidKafkaWriterConfig = errors.New("kafka writer config is invalid")
	ErrKafkaCouldNotWriteEvent  = errors.New("kafka writer failed writing the event")
	ErrKafkaCouldNotCloseWriter = errors.New("kafka writer failed closing")
	ErrEventMarhsallingFailed   = errors.New("failed marshalling event")
	ErrFailedSteamingEvents     = errors.New("failed publishing the events through the kafka stream")
)

func CreateKafkaStream(topic es.Topic) *Stream {
	return &Stream{
		topic: topic,
	}
}

func (stream *Stream) write(ctx context.Context, config kafka.WriterConfig, events []es.Event) error {
	var defaultHeaders []kafka.Header

	writer := kafka.NewWriter(config)

	for _, event := range events {
		value, err := event.Marshall()
		if err != nil {
			return errors.Wrap(err, ErrEventMarhsallingFailed.Error())
		}

		message := kafka.Message{
			Key:           []byte(event.Subject),
			Value:         value,
			Topic:         string(stream.topic),
			Partition:     defaultPartition,
			Offset:        defaultOffset,
			HighWaterMark: defaultHighWaterMark,
			Headers:       defaultHeaders,
			Time:          time.Now(),
		}

		if err := writer.WriteMessages(ctx, message); err != nil {
			return errors.Wrap(err, ErrKafkaCouldNotWriteEvent.Error())
		}
	}

	return errors.Wrap(
		writer.Close(),
		ErrKafkaCouldNotCloseWriter.Error(),
	)
}

func (stream *Stream) Publish(events []es.Event) error {
	var defaultLogger kafka.Logger

	var defaultErrorLogger kafka.Logger

	var defaultCompressionCodec kafka.CompressionCodec

	defaultBalancer := &kafka.RoundRobin{}
	defaultDialer := kafka.DefaultDialer

	config := kafka.WriterConfig{
		Brokers:           []string{broker},
		Topic:             string(stream.topic),
		Dialer:            defaultDialer,
		Balancer:          defaultBalancer,
		MaxAttempts:       defaultMaxAttempts,
		QueueCapacity:     defaultQueueCapacity,
		BatchSize:         defaultBatchSize,
		BatchBytes:        defaultBatchBytes,
		BatchTimeout:      defaultBatchTimeout,
		ReadTimeout:       defaultReadTimeout,
		WriteTimeout:      defaultWriteTimeout,
		RebalanceInterval: defaultRebalanceInterval,
		IdleConnTimeout:   defaultIdleConnTimeout,
		RequiredAcks:      defaultRequiredAcks,
		Async:             defaultAsync,
		CompressionCodec:  defaultCompressionCodec,
		Logger:            defaultLogger,
		ErrorLogger:       defaultErrorLogger,
	}

	if err := config.Validate(); err != nil {
		panic(errors.Wrap(err, ErrInvalidKafkaWriterConfig.Error()))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	return errors.Wrap(
		stream.write(ctx, config, events),
		ErrFailedSteamingEvents.Error(),
	)
}

func (stream *Stream) read(
	ctx context.Context,
	config kafka.ReaderConfig,
	events chan es.Event,
	errors chan error,
) {
	reader := kafka.NewReader(config)
	defer reader.Close() // What if closing the reader fails?

	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			errors <- err
		}

		event, err := es.UnmarshalEvent(msg.Value)
		if err != nil {
			errors <- err
		}
		events <- event
	}
}

func (stream *Stream) Subscribe(ctx context.Context, topic es.Topic) (chan es.Event, chan error) {
	var defaultLogger kafka.Logger

	var defaultErrorLogger kafka.Logger

	config := kafka.ReaderConfig{
		Brokers:         []string{broker},
		GroupID:         group,
		GroupTopics:     nil, // Defined through "Topic"
		Topic:           string(topic),
		Partition:       defaultPartition, // the same as undefined
		Dialer:          nil,
		QueueCapacity:   defaultQueueCapacity,
		MinBytes:        defaultMinBytes,
		MaxBytes:        defaultMaxBytes, // 1 MB,
		MaxWait:         defaultMaxWait,
		ReadLagInterval: defaultReadLagInterval,
		GroupBalancers: []kafka.GroupBalancer{
			kafka.RangeGroupBalancer{},
			kafka.RoundRobinGroupBalancer{},
		},
		HeartbeatInterval:      defaultHeartbeatInterval,
		CommitInterval:         defaultCommitInterval,
		PartitionWatchInterval: defaultPartitionWatchInterval,
		WatchPartitionChanges:  defaultWatchPartitionChanges,
		SessionTimeout:         defaultSessionTimeout,
		RebalanceTimeout:       defaultRebalanceTimeout,
		JoinGroupBackoff:       defaultJoinGroupBackoff,
		RetentionTime:          defaultRetentionTime,
		StartOffset:            defaultStartOffset,
		ReadBackoffMin:         defaultReadBackoffMin,
		ReadBackoffMax:         defaultReadBackoffMax,
		Logger:                 defaultLogger,
		ErrorLogger:            defaultErrorLogger,
		IsolationLevel:         defaultIsolationLevel,
		MaxAttempts:            defaultMaxAttempts,
	}

	if err := config.Validate(); err != nil {
		panic(errors.Wrap(err, ErrInvalidKafkaReaderConfig.Error()))
	}

	events := make(chan es.Event)
	errors := make(chan error)

	go stream.read(ctx, config, events, errors)

	return events, errors
}
