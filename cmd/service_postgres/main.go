package main

import (
	"log"
	"math/rand"

	"github.com/google/uuid"
	"github.com/hywmongous/example-service/pkg/es"
	"github.com/hywmongous/example-service/pkg/es/mongo"
)

// "github.com/hywmongous/example-service/internal/presentation/connectors/gin/bootstrap"
/*"github.com/hywmongous/example-service/pkg/es"
"github.com/hywmongous/example-service/pkg/es/mongo"*/

type Transaction struct {
	AggregateId string
	Amount      int
}

func test(asd struct{}) {

}
func main() {
	// fx.New(bootstrap.Module).Run()

	aggregateId := uuid.NewString()

	producer := es.ProducerID("Andreas")
	subject := es.SubjectID("Insert Aggregate ID")

	var eventStore es.EventStore
	eventStore, err := mongo.CreateMongoEventStore()
	if err != nil {
		log.Fatal(err)
	}

	data := []es.EventData{
		Transaction{
			AggregateId: aggregateId,
			Amount:      rand.Intn(100),
		},
		Transaction{
			AggregateId: aggregateId,
			Amount:      rand.Intn(100),
		},
		Transaction{
			AggregateId: aggregateId,
			Amount:      rand.Intn(100),
		},
	}

	if err := eventStore.Stock(producer, subject, data); err != nil {
		log.Fatal(err)
	}

	callback := func(event es.Event) error {
		var transaction Transaction
		event.Unmarshal(&transaction)

		println("Transaction:")
		println(transaction.AggregateId)
		println(transaction.Amount)

		return nil
	}
	eventStore.Retrieve(subject, callback)

	/*
		latestEvent, err := eventStore.LatestEvent(subject, reflect.TypeOf((*Transaction)(nil)))
		if err != nil {
			log.Fatal(err)
		}

		var transaction Transaction
		latestEvent.Unmarshal(&transaction)
		println("Transaction:")
		println(transaction.AggregateId)
		println(transaction.Amount)
	*/

	/*
		uri := options.Client().ApplyURI("mongodb://root:root@ia_mongo:27017")
		client, err := mongo.NewClient(uri)
		if err != nil {
			log.Fatal(err)
		}
		ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
		err = client.Connect(ctx)
		if err != nil {
			log.Fatal(err)
		}
		defer client.Disconnect(ctx)

		database := client.Database("eventstore")
		if database == nil {
			log.Fatal("Database is nil")
		}
		eventCollection := database.Collection("events")
		if eventCollection == nil {
			log.Fatal("Collection is nil")
		}
		_, err = eventCollection.InsertOne(ctx, bson.D{
			{Key: "Key1", Value: "Values1"},
		})
		if err != nil {
			log.Fatal(err)
		}
	*/

	/*
		producer := es.ProducerID("Andreas")
		println("A")
		registry, err := es.CreateEventRegistry()
		if err != nil {
			println(err)
			return
		}
		registry.RegisterEvent(Transaction{})
		println("B")
		eventstore, err := mongo.CreateEventStore(registry)
		if err != nil {
			println(err)
			return
		}
		println("C")
		err = eventstore.Send(producer, Transaction{amount: 10})
		if err != nil {
			println(err)
			return
		}
		println("D")
		events, err := eventstore.Load(producer)
		if err != nil {
			println(err)
			return
		}
		println("E")
		for _, event := range events {
			println(event.GetData())
		}*/
}
