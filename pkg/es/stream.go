package es

import (
	"context"
	"log"
)

type EventStream interface {
	Publish(events []Event) error
	Subscribe(topic Topic, ctx context.Context) (chan Event, chan error)
}

type (
	EventReceived func(event Event)
	Topic         string
)

func CreateEventStream(events []Event) chan Event {
	channel := make(chan Event, len(events))
	for _, event := range events {
		// We dont check for "ok" because we wont close it
		channel <- event
	}
	return channel
}

func CreateErrorPrinter() chan error {
	errors := make(chan error)

	printer := func() {
		for {
			err, ok := <-errors
			if !ok {
				break
			}

			log.Println(err)
		}
	}
	go printer()
	return errors
}
