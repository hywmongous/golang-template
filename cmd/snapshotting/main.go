package main

import (
	"log"
	"time"

	"github.com/hywmongous/example-service/pkg/es"
	"github.com/hywmongous/example-service/pkg/es/mongo"
)

func main() {
	// The goal of this example is to show how snapshotting
	// can eliminate the need of the previous events

	time.Sleep(10 * time.Second)

	// producer := es.ProducerID("Producer")
	subject := es.SubjectID("Subject")

	store := mongo.CreateMongoEventStore()

	snapshot, err := store.LatestSnapshot(subject)
	if err != nil {
		log.Fatal("LatestSnapshot:", err)
	}

	events, err := store.With(subject, snapshot.Version)
	if err != nil {
		log.Fatal("With:", err)
	}
	for _, event := range events {
		log.Println(event)
	}

	time.Sleep(10 * time.Second)
}
