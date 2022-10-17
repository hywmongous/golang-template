package main

/*
import (
	"context"
	"crypto/rand"
	"log"
	"math/big"
	"time"

	"github.com/hywmongous/example-service/pkg/es"
	"github.com/hywmongous/example-service/pkg/es/kafka"
	"github.com/hywmongous/example-service/pkg/es/mongo"
)



type Transaction struct {
	AggregateID string
	Amount      int
}

const (
	sleepTime = 30 * time.Second
	producer  = es.ProducerID("Producer")
	subject   = es.SubjectID("Subject")
)

func main() {
	// We do this because we have to wait for kafka initialization
	time.Sleep(sleepTime)

	log.Print("Event creation")

	var maxAmount int64 = 100

	amount, err := rand.Int(rand.Reader, big.NewInt(maxAmount))
	if err != nil {
		log.Fatal(err)
	}

	eventData := Transaction{
		AggregateID: string(subject),
		Amount:      int(amount.Int64()),
	}

	log.Print("Event store")

	store := mongo.CreateMongoEventStore()

	log.Print("Commit event data")

	if err := store.Load(producer, subject, eventData); err != nil {
		log.Fatal(err)
	}

	log.Print("Event shipping")

	if err := store.Ship(); err != nil {
		log.Fatal(err)
	}

	log.Print("Event querying")

	stock, err := store.After(subject, es.Timestamp(time.Now().AddDate(0, 0, -2).Unix()))
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Stock")

	for _, event := range stock {
		log.Println(event.ID)
	}

	log.Print("Event streaming")

	topic := es.Topic("ia")
	stream := kafka.CreateKafkaStream(topic)

	log.Print("Event publications")

	events := store.Stage().Events()
	if err = stream.Publish(events); err != nil {
		log.Fatal(err)
	}

	log.Print("Event subscription")

	ctx := context.Background()
	subscriptions, errs := stream.Subscribe(ctx, topic)

	go func() {
		log.Println("Subscription")

		for {
			select {
			case event, ok := <-subscriptions:
				if !ok {
					break
				}

				log.Println(event.ID)
			case err, ok := <-errs:
				if !ok {
					break
				}

				log.Fatal(err)
			case <-ctx.Done():
				if err := ctx.Err(); err != nil {
					log.Fatal(err)
				}
				// this if/else displays the different reasons for Done
				// if errors.Is(err, context.Canceled) {
				// 	log.Fatal(err)
				// } else if errors.Is(err, context.DeadlineExceeded) {
				// 	log.Fatal(err)
				// } else {
				// 	log.Fatal(err)
				// }
				return
			}
		}
	}()

	time.Sleep(sleepTime)
}
*/
