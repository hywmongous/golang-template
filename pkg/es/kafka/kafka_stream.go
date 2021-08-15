package kafka

import (
	"context"
	"log"

	"github.com/hywmongous/example-service/pkg/es"
	"github.com/segmentio/kafka-go"
)

type KafkaStream struct{}

const (
	broker = "ia_kafka:9092"
)

func CreateKafkaStream() KafkaStream {
	return KafkaStream{}
}

func (stream KafkaStream) write(ctx context.Context, config kafka.WriterConfig, events chan es.Event, errors chan error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	writer := kafka.NewWriter(config)
	for {
		event, ok := <-events
		if !ok {
			break
		}

		value, err := event.Marshall()
		if err != nil {
			errors <- err
			break
		}

		message := kafka.Message{
			Key:   []byte(event.Subject),
			Value: value,
		}

		if err := writer.WriteMessages(ctx, message); err != nil {
			errors <- err
			break
		}
	}

	if err := writer.Close(); err != nil {
		errors <- err
	}
}

func (stream KafkaStream) Publish(topic es.Topic, events chan es.Event, errors chan error) context.CancelFunc {
	config := kafka.WriterConfig{
		Brokers: []string{broker},
		Topic:   string(topic),
	}

	ctx, cancel := context.WithCancel(context.Background())
	go stream.write(ctx, config, events, errors)
	return cancel
}

func (stream KafkaStream) read(ctx context.Context, config kafka.ReaderConfig, events chan es.Event, errors chan error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	reader := kafka.NewReader(config)
	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			errors <- err
			break
		}

		event, err := es.Unmarshal(msg.Value)
		if err != nil {
			errors <- err
			break
		}
		events <- event
	}

	if err := reader.Close(); err != nil {
		errors <- err
	}
}

func (stream KafkaStream) Subscribe(topic es.Topic, events chan es.Event, errors chan error) context.CancelFunc {
	config := kafka.ReaderConfig{
		Brokers: []string{broker},
		Topic:   string(topic),
	}

	ctx, cancel := context.WithCancel(context.Background())
	go stream.read(ctx, config, events, errors)
	return cancel
}

func (stream KafkaStream) PrintErrors() chan error {
	errors := make(chan error)
	go func() {
		for {
			err, ok := <-errors
			if !ok {
				break
			}
			log.Println(err)
		}
	}()
	return errors
}
