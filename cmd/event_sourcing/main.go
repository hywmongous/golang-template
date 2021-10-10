package main

import (
	"log"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
	"github.com/hywmongous/example-service/pkg/es"
	"github.com/hywmongous/example-service/pkg/es/kafka"
	"github.com/hywmongous/example-service/pkg/es/mediator"
	"github.com/hywmongous/example-service/pkg/es/mongo"
)

// Infrastructure
type UnitOfWork struct {
	store  es.EventStore
	stream es.EventStream
}

// Aggregate
type Identity struct {
	Id    string
	Name  string
	Age   int
	Email string
}

// Snapshot(s)
type IdentitySnapshotV1 struct {
	Id    string
	Name  string
	Age   int
	Email string
}

// Actors
type UnregisteredUser struct{}

// Use cases
type IdentityRegistrationRequest struct {
	Name  string
	Age   int
	Email string
}
type IdentityRegistrationResponse struct {
	Id string
}
type RegisterIdentityUseCase func(request IdentityRegistrationRequest) (IdentityRegistrationResponse, error)

type IdentityChangeNameRequest struct {
	Id   string
	Name string
}
type IdentityChangeNameResponse struct {
	Success bool
}
type ChangeIdentityNameUseCase func(request IdentityChangeNameRequest) (IdentityChangeNameResponse, error)

// Events
type IdentityRegistered struct {
	Id    string
	Name  string
	Age   int
	Email string
}
type IdentityChangedName struct {
	Name string
}
type IdentityChangedAge struct {
	Age int
}
type IdentityChangedEmail struct {
	Email string
}

var (
	Producer = es.ProducerID("EventSourcingExample")
)

func main() {
	// Scenario: A user reigstration is issued and afterwards the events are read to construct the aggregate

	// We sleep here in this example jsut to ensure initialization of external services are done
	time.Sleep(15 * time.Second)

	log.Println("EVENT SOURCING EXAMPLE")

	// Step 1: Create the unit of work
	mongoStore := mongo.CreateMongoEventStore()
	kafkaStram := kafka.CreateKafkaStream(
		es.Topic("ia"),
	)
	uow := UnitOfWork{
		store:  mongoStore,
		stream: kafkaStram,
	}
	mediator.Listen(uow.receiveEvent)

	// Create actor
	unregisterUser := UnregisteredUser{}

	// Do registration
	var registrationUseCase RegisterIdentityUseCase = unregisterUser.Register
	registrationRequest := IdentityRegistrationRequest{
		Name:  "Andreas",
		Age:   22,
		Email: "andreasbrandhoej@hotmail.com",
	}
	registrationResponse, err := registrationUseCase(registrationRequest)
	if err != nil {
		log.Println(err)
		if err = uow.Rollback(); err != nil {
			log.Println("rollback failed")
			log.Println(err)
		}
		log.Fatal("Registration use case failed")
	}
	// log.Println("Identity:", registrationResponse.Id, "was registered")

	// Change name
	var changeNameUseCase ChangeIdentityNameUseCase = unregisterUser.ChangeName
	nameChangeRequest := IdentityChangeNameRequest{
		Id:   registrationResponse.Id,
		Name: "Andreas K. Brandhøj",
	}
	_, err = changeNameUseCase(nameChangeRequest)
	if err != nil {
		log.Println(err)
		if err = uow.Rollback(); err != nil {
			log.Println("rollback failed")
			log.Println(err)
		}
		log.Fatal("Name change of identity failed")
	}
	// log.Println("Identity:", registrationResponse.Id, "changed name to '"+nameChangeRequest.Name+"'")

	// Commit changes made through the use cases
	if err = uow.Commit(); err != nil {
		log.Fatal(err)
	}

	// For this example we wait a bit after the registration
	time.Sleep(2 * time.Second)

	// Query all the events for the aggregate
	events, err := uow.store.Concerning(es.SubjectID(registrationResponse.Id))
	if err != nil {
		log.Fatal(err)
	}

	// Construct the aggregate
	var identity Identity
	identity = VisitEvents(identity, events)

	// Snapshot the aggregate
	snapshotData1 := IdentitySnapshotV1{
		Id:    identity.Id,
		Name:  identity.Name,
		Age:   420,
		Email: identity.Email,
	}

	err = uow.store.Snapshot(Producer, es.SubjectID(identity.Id), snapshotData1)
	if err != nil {
		log.Fatal(err)
	}
	// log.Println("Create snapshot version", snapshot.Version)
	// log.Println("Snapshot", snapshot)

	// Do something
	_, err = changeNameUseCase(IdentityChangeNameRequest{
		Id:   identity.Id,
		Name: "Dat Tommy Than Dieu",
	})
	if err != nil {
		log.Println(err)
		if err = uow.Rollback(); err != nil {
			log.Println("rollback failed")
			log.Println(err)
		}
		log.Fatal("Name change of identity failed")
	}

	// Snapshot the aggregate
	snapshotData2 := IdentitySnapshotV1{
		Id:    identity.Id,
		Name:  identity.Name,
		Age:   69,
		Email: "ddieu19@student.aau.dk",
	}

	// Make another snapshot
	err = uow.store.Snapshot(Producer, es.SubjectID(identity.Id), snapshotData2)
	if err != nil {
		log.Fatal(err)
	}
	// log.Println("Create snapshot version", snapshot.Version)
	// log.Println("Snapshot", snapshot)

	// Do something
	_, err = changeNameUseCase(IdentityChangeNameRequest{
		Id:   identity.Id,
		Name: "Dat Tommy Than Dieu",
	})
	if err != nil {
		log.Println(err)
		if err = uow.Rollback(); err != nil {
			log.Println("rollback failed")
			log.Println(err)
		}
		log.Fatal("Name change of identity failed")
	}
	// Do something
	_, err = changeNameUseCase(IdentityChangeNameRequest{
		Id:   identity.Id,
		Name: "Dat Tommy Than Dieu",
	})
	if err != nil {
		log.Println(err)
		if err = uow.Rollback(); err != nil {
			log.Println("rollback failed")
			log.Println(err)
		}
		log.Fatal("Name change of identity failed")
	}
	// Do something
	_, err = changeNameUseCase(IdentityChangeNameRequest{
		Id:   identity.Id,
		Name: "Dat Tommy Than Dieu",
	})
	if err != nil {
		log.Println(err)
		if err = uow.Rollback(); err != nil {
			log.Println("rollback failed")
			log.Println(err)
		}
		log.Fatal("Name change of identity failed")
	}
	// Do something
	_, err = changeNameUseCase(IdentityChangeNameRequest{
		Id:   identity.Id,
		Name: "Dat Tommy Than Dieu",
	})
	if err != nil {
		log.Println(err)
		if err = uow.Rollback(); err != nil {
			log.Println("rollback failed")
			log.Println(err)
		}
		log.Fatal("Name change of identity failed")
	}

	// Commit with the snapshot
	if err = uow.Commit(); err != nil {
		log.Fatal(err)
	}

	// Create Identity from fetched snapshot
	latestSnapshot, err := uow.store.LatestSnapshot(es.SubjectID(identity.Id))
	if err != nil {
		log.Fatal(err)
	}
	finished_identity := VisitSnapshot(latestSnapshot)

	// And the append proceeding events made after the snapshot
	proceedingEvents, err := uow.store.With(es.SubjectID(identity.Id), latestSnapshot.Version)
	if err != nil {
		log.Fatal(err)
	}
	// log.Println("Found", len(proceedingEvents), "events with snapshot version", latestSnapshot.Version)
	finished_identity = VisitEvents(finished_identity, proceedingEvents)

	// print the final identity
	log.Println("finished_identity", finished_identity.Id, finished_identity.Name, finished_identity.Age, finished_identity.Email)

	// For this example we sleep here such that you have the time to verify the output
	time.Sleep(5 * time.Second)
}

func VisitSnapshot(snapshot es.Snapshot) Identity {
	identity := Identity{}
	log.Println("Visit:", snapshot.Name, snapshot.Version, ",", snapshot.Data)
	switch snapshot.Name {
	case "IdentitySnapshotV1":
		var data IdentitySnapshotV1
		if err := snapshot.Unmarshal(&data); err != nil {
			panic(err)
		}
		identity.Id = data.Id
		identity.Name = data.Name
		identity.Age = data.Age
		identity.Email = data.Email
	}
	return identity
}

func VisitEvents(identity Identity, events []es.Event) Identity {
	for _, event := range events {
		log.Println("Visit:", event.Name, event.Version, ":", event.SnapshotVersion, ",", event.Data)
		switch event.Name {
		case "IdentityRegistered":
			var data IdentityRegistered
			if err := event.Unmarshal(&data); err != nil {
				panic(err)
			}
			identity.Id = data.Id
			identity.Age = data.Age
			identity.Email = data.Email
			identity.Name = data.Name
		case "IdentityChangedName":
			var data IdentityChangedName
			if err := event.Unmarshal(&data); err != nil {
				panic(err)
			}
			identity.Name = data.Name
		case "IdentityChangedAge":
			var data IdentityChangedAge
			if err := event.Unmarshal(&data); err != nil {
				panic(err)
			}
			identity.Age = data.Age
		case "IdentityChangedEmail":
			var data IdentityChangedEmail
			if err := event.Unmarshal(&data); err != nil {
				panic(err)
			}
			identity.Email = data.Email
		}
	}
	return identity
}

func (unregisteredUser UnregisteredUser) Register(request IdentityRegistrationRequest) (IdentityRegistrationResponse, error) {
	response := IdentityRegistrationResponse{
		Id: uuid.New().String(),
	}
	mediator.Publish(
		es.SubjectID(response.Id),
		IdentityRegistered{
			Id:    response.Id,
			Name:  request.Name,
			Age:   request.Age,
			Email: request.Email,
		},
	)
	return response, nil
}

func (unregisteredUser UnregisteredUser) ChangeName(request IdentityChangeNameRequest) (IdentityChangeNameResponse, error) {
	response := IdentityChangeNameResponse{
		Success: true,
	}
	mediator.Publish(
		es.SubjectID(request.Id),
		IdentityChangedName{
			Name: request.Name,
		},
	)
	return response, nil
}

func (uow *UnitOfWork) receiveEvent(subject es.SubjectID, data es.Data) {
	err := uow.store.Load(
		Producer,
		subject,
		data,
	)
	panic(err)
}

func (uow *UnitOfWork) Commit() error {
	events := uow.store.Stage().Events()
	if err := uow.store.Ship(); err != nil {
		log.Panicln("UnitOfWork store failed shipping/storing the events")
		return errors.Wrap(err, "UnitOfWork storing failed")
	}
	if err := uow.stream.Publish(events); err != nil {
		log.Panicln("UnitOfWork store failed shipping/publishing the events")
		return errors.Wrap(err, "UnitOfWork publishing failed")
	}
	return nil
}

func (uow *UnitOfWork) Rollback() error {
	uow.store.Clear()
	return nil
}
