package main

import (
	"log"
	"time"

	"github.com/hywmongous/example-service/pkg/es"
	"github.com/hywmongous/example-service/pkg/es/mediator"
)

type Event1 struct {
	time es.Timestamp
}

type Event2 struct {
	time es.Timestamp
}

const sleepTime = 100 * time.Millisecond

func main() {
	mediator := mediator.Create()
	topic1 := es.Topic("topic1")
	topic2 := es.Topic("topic2")

	// Listeners
	mediator.ListenTo(topic1, listener1)
	mediator.Listen(universalListener1)

	// Connectors
	connector1 := mediator.ChannelTo(topic2)
	connector1Func := func() {
		data := <-connector1
		log.Println("connector1:", data)
	}

	go connector1Func()

	// Create events
	event1 := Event1{
		time: es.Timestamp(time.Now().Unix()),
	}
	event2 := Event2{
		time: es.Timestamp(time.Now().Unix()),
	}

	// Publish
	mediator.Publish(es.SubjectID("me"), event1)
	mediator.Publish(es.SubjectID("me"), event2)

	// We sleep in this example to ensure the channel has received the event data
	time.Sleep(sleepTime)
}

func listener1(subject es.SubjectID, data es.Data) {
	log.Println("listener1:", data)
}

func universalListener1(subject es.SubjectID, data es.Data) {
	log.Println("universalListener1:", data)
}
