package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/hywmongous/example-service/pkg/es"
	"github.com/hywmongous/example-service/pkg/es/kafka"
	"github.com/hywmongous/example-service/pkg/es/mongo"
)

// "github.com/hywmongous/example-service/internal/presentation/connectors/gin/bootstrap"
/*"github.com/hywmongous/example-service/pkg/es"
"github.com/hywmongous/example-service/pkg/es/mongo"*/

type Transaction struct {
	AggregateId string
	Amount      int
}

func main() {
	// fx.New(bootstrap.Module).Run()

	// We do this because we have to wait for kafka initialization
	time.Sleep(20 * time.Second)

	// Event
	producer := es.ProducerID("Producer")
	subject := es.SubjectID("Subject")
	eventData := []es.EventData{
		Transaction{
			AggregateId: string(subject),
			Amount:      rand.Intn(100),
		},
	}

	// Event Store
	store := mongo.CreateMongoEventStore()

	// Event Stocking
	events, err := store.Send(producer, subject, eventData)
	if err != nil {
		log.Fatal(err)
	}

	// Event Retriving
	stock, err := store.After(subject, es.EventTimestamp(time.Now().AddDate(0, 0, -2).Unix()))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Stock")
	for _, event := range stock {
		log.Println(time.Unix(int64(event.Timestamp), 0).String())
	}

	// Event Streaming
	stream := kafka.CreateKafkaStream()
	errors := stream.PrintErrors()
	topic := es.Topic("ia.identity")

	// Event publications
	publications := es.CreateEventStream(events)
	stream.Publish(topic, publications, errors)

	// Event subscription
	subscriptions := make(chan es.Event)
	stream.Subscribe(topic, subscriptions, errors)
	go func() {
		log.Println("Subscription")
		for {
			event, ok := <-subscriptions
			if !ok {
				break
			}
			log.Println(event.Id)
		}
	}()

	// We do this because we have to ensure we have received the
	// events and the goroutines have finished their word
	time.Sleep(20 * time.Second)
}
