package kafka

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/hywmongous/example-service/pkg/es"
	"github.com/segmentio/kafka-go"
)

type KafkaStream struct {
	// We only use a single topic per "bounded context"
	// For this reason it fits the "topic" to be a field
	// within the kafka stream struct instance
	// https://www.confluent.io/blog/put-several-event-types-kafka-topic/
	topic es.Topic
}

const (
	defaultPartition              = 0
	defaultQueueCapacity          = 100
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
	defaultMaxAttempts            = 3

	broker = "ia_kafka:9092"
	group  = "ia"
)

var (
	defaultLogger      kafka.Logger = nil
	defaultErrorLogger kafka.Logger = nil
)

var ErrInvalidKafkaReaderConfig = errors.New("kafka reader config is invalid")

func CreateKafkaStream(topic es.Topic) *KafkaStream {
	return &KafkaStream{
		topic: topic,
	}
}

func (stream *KafkaStream) write(ctx context.Context, config kafka.WriterConfig, events []es.Event) error {
	writer := kafka.NewWriter(config)
	for _, event := range events {
		value, err := event.Marshall()
		if err != nil {
			return err
		}

		message := kafka.Message{
			Key:   []byte(event.Subject),
			Value: value,
		}

		if err := writer.WriteMessages(ctx, message); err != nil {
			return err
		}
	}

	return writer.Close()
}

func (stream *KafkaStream) Publish(events []es.Event) error {
	config := kafka.WriterConfig{
		Brokers: []string{broker},
		Topic:   string(stream.topic),
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	return stream.write(ctx, config, events)
}

func (stream *KafkaStream) read(ctx context.Context, config kafka.ReaderConfig, events chan es.Event, errors chan error) {
	reader := kafka.NewReader(config)
	defer reader.Close() // TODO: What if closing the reader fails?
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

func (stream *KafkaStream) Subscribe(topic es.Topic, ctx context.Context) (chan es.Event, chan error) {
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
