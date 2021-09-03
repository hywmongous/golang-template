package es

import (
	"context"
	"log"
)

type EventStream interface {
	// TODO: It does not make sense to channel the errors
	// There should be underlying rules which governs how
	// we handle the errors. Possible options:
	//   1. Retry, a specific amount
	//   2. Stop, we completely ignore the rest of the events
	//   3. Ignore, we just continue with the next one and
	//        accept the probability og loss - It would be insane to expect another response immediately
	// I 4. Continue, with the next
	//   5. Store, emit them at the next publish of other events
	// Idea: We have to know which publication (publish of event) failed
	//   when the error occures. However, we have to take into account
	//   the channel of errors is also used to propagate errorsregarding
	//   the client termination of socket listening. I see three solutions
	//     1. add an addition channel for errors and chane the current
	//          one to return the failed publication
	//     2. Merge the above mentioned idea into a channel returning a struct
	//     3. Ignore the client close. What can we do about it? I think: Nothing.
	// This discussion is closed, on the realisation that we should stop streaming
	//   when an error occurs because we are required to publish them in a sequential order
	//   this simplifies the whole process of commiting the changes to topics.
	Publish(events []Event) error
	Subscribe(topic Topic, ctx context.Context) (chan Event, chan error)
	// Publish(events chan Event, errors chan error) context.CancelFunc
	// Subscribe(topic Topic, errors chan error) (chan Event, context.CancelFunc)
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
