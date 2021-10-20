package main

import (
	"log"

	"github.com/hywmongous/example-service/pkg/es"
)

const subject = es.SubjectID("Andreas")

func createEvent(version es.Version) es.Event {
	return es.Event{
		Subject: subject,
		Version: version,
	}
}

func createSnapshot(version es.Version) es.Snapshot {
	return es.Snapshot{
		Subject: subject,
		Version: version,
	}
}

func main() {
	log.Println("Staging example")

	stage := createStage()

	log.Println("Is empty (after adding)", stage.IsEmpty(subject))

	log.Println("---Stages---")

	stageEvents(stage)

	log.Println("---Subjects---")

	stageSubjects(stage)

	log.Println("Is empty (after clearing)", stage.IsEmpty(subject))
}

func createStage() es.Stage {
	stage := es.CreateStage()
	currEventVersion := 1

	event1 := createEvent(es.Version(currEventVersion))
	currEventVersion++

	event2 := createEvent(es.Version(currEventVersion))

	currEventVersion++

	event3 := createEvent(es.Version(currEventVersion))

	currEventVersion++

	event4 := createEvent(es.Version(currEventVersion))

	currEventVersion++

	event5 := createEvent(es.Version(currEventVersion))

	currEventVersion++

	event6 := createEvent(es.Version(currEventVersion))

	currSnapshotVersion := 1

	snapshot1 := createSnapshot(es.Version(currSnapshotVersion))

	currSnapshotVersion++

	snapshot2 := createSnapshot(es.Version(currSnapshotVersion))

	currSnapshotVersion++

	snapshot3 := createSnapshot(es.Version(currSnapshotVersion))

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

	return stage
}

func stageEvents(stage es.Stage) {
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
}

func stageSubjects(stage es.Stage) {
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
}
