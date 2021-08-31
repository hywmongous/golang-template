package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hywmongous/example-service/pkg/es"
	"github.com/hywmongous/example-service/pkg/es/kafka"
	"github.com/hywmongous/example-service/pkg/es/mongo"
)

type Transaction struct {
	AggregateId string
	Amount      int
}

func main() {
	// We do this because we have to wait for kafka initialization
	time.Sleep(30 * time.Second)

	log.Print("Event creation")
	producer := es.ProducerID("Producer")
	subject := es.SubjectID("Subject")
	eventData := Transaction{
		AggregateId: string(subject),
		Amount:      rand.Intn(100),
	}

	log.Print("Event store")
	store := mongo.CreateMongoEventStore()

	log.Print("Commit event data")
	ware, err := store.Load(producer, subject, eventData)
	log.Print("Loaded: ", ware.Id)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Event shipping")
	events, err := store.Ship()
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Event querying")
	stock, err := store.After(subject, es.Timestamp(time.Now().AddDate(0, 0, -2).Unix()))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Stock")
	for _, event := range stock {
		log.Println(event.Id)
	}

	log.Print("Event streaming")
	stream := kafka.CreateKafkaStream()
	errors := stream.CreateErrorPrinter()
	topic := es.Topic("ia.identity")

	log.Print("Event publications")
	publications := es.CreateEventStream(events)
	stream.Publish(topic, publications, errors)

	log.Print("Event subscription")
	subscriptions, cancel := stream.Subscribe(topic, errors)
	go func() {
		log.Println("Subscription")
		for {
			event, ok := <-subscriptions
			if !ok {
				break
			}
			log.Println(event.Id)
		}
		cancel()
	}()

	// We do this because we have to ensure we have received the
	// events and the goroutines have finished their word
	time.Sleep(30 * time.Second)
}
