package main

import (
	"log"

	"github.com/hywmongous/example-service/pkg/es"
)

var subject = es.SubjectID("Andreas")

func main() {
	log.Println("Staging example")

	stage := es.CreateStage()

	event1 := es.Event{
		Version: 1,
		Subject: subject,
	}
	event2 := es.Event{
		Version: 2,
		Subject: subject,
	}
	event3 := es.Event{
		Version: 3,
		Subject: subject,
	}
	event4 := es.Event{
		Version: 4,
		Subject: subject,
	}
	event5 := es.Event{
		Version: 5,
		Subject: subject,
	}
	event6 := es.Event{
		Version: 6,
		Subject: subject,
	}
	snapshot1 := es.Snapshot{
		Subject: subject,
		Version: 1,
	}
	snapshot2 := es.Snapshot{
		Subject: subject,
		Version: 2,
	}
	snapshot3 := es.Snapshot{
		Subject: subject,
		Version: 3,
	}

	log.Println("Is empty (before adding)", stage.IsEmpty(subject))

	stage.AddEvent(event1)
	stage.AddEvent(event2)
	stage.AddSnapshot(snapshot1)
	stage.AddEvent(event3)

	stage.Clear(subject)

	stage.AddEvent(event4)
	stage.AddEvent(event5)
	stage.AddSnapshot(snapshot2)
	stage.AddEvent(event6)
	stage.AddSnapshot(snapshot3)

	log.Println("Is empty (after adding)", stage.IsEmpty(subject))

	log.Println("---Stages---")
	eventStages := stage.EventStages(subject)
	for idx, eventStage := range eventStages {
		log.Println("Stage", idx)
		for _, event := range eventStage.Events() {
			log.Println(event)
		}

		if eventStage.Snapshot() != nil {
			log.Println("Snapshot version", eventStage.Snapshot().Version)
		} else {
			log.Println("No snapshot")
		}
	}

	if firstEvent, found := stage.FirstEvent(subject); found {
		log.Println("First event version", firstEvent.Version)
	} else {
		log.Println("No first event found")
	}

	if latestEvent, found := stage.LatestEvent(subject); found {
		log.Println("Latest event version", latestEvent.Version)
	} else {
		log.Println("No latest event found")
	}

	if latestSnapshot, found := stage.LatestSnapshot(subject); found {
		log.Println("Latest snapshot version", latestSnapshot.Version)
	} else {
		log.Println("No latest snapshot found")
	}

	log.Println("---Subjects---")
	subjects := stage.Subjects()
	for _, subject := range subjects {
		log.Println(subject)
	}

	stage.Clear(subject)

	log.Println("---Stages (After clearing)---")
	eventStagesAfterClearing := stage.EventStages(subject)
	for idx, eventStage := range eventStagesAfterClearing {
		log.Println("Stage", idx)
		for _, event := range eventStage.Events() {
			log.Println(event)
		}

		if eventStage.Snapshot() != nil {
			log.Println("Snapshot version", eventStage.Snapshot().Version)
		} else {
			log.Println("No snapshot")
		}
	}

	log.Println("Is empty (after clearing)", stage.IsEmpty(subject))
}
