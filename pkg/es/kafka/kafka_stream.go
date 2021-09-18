package kafka

import (
	"context"

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
	broker = "ia_kafka:9092"
	group  = "ia"
)

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
		event, err := es.Unmarshal(msg.Value)
		if err != nil {
			errors <- err
		}
		events <- event
	}
}

func (stream *KafkaStream) Subscribe(topic es.Topic, ctx context.Context) (chan es.Event, chan error) {
	config := kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   string(topic),
		GroupID: group,
	}

	events := make(chan es.Event)
	errors := make(chan error)

	go stream.read(ctx, config, events, errors)
	return events, errors
}
