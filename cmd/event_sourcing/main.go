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

type IdentityChangeAgeRequest struct {
	Id  string
	Age int
}
type IdentityChangeAgeResponse struct {
	Success bool
}
type ChangeIdentityAgeUseCase func(request IdentityChangeNameRequest) (IdentityChangeNameResponse, error)

type IdentityChangeEmailRequest struct {
	Id    string
	Email string
}
type IdentityChangeEmailResponse struct {
	Success bool
}
type ChangeIdentityEmailUseCase func(request IdentityChangeNameRequest) (IdentityChangeNameResponse, error)

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
		store:  &mongoStore,
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
		uow.Rollback()
		log.Fatal("Registration use case failed")
	}
	log.Println("Identity:", registrationResponse.Id, "was registered")

	// Change name
	var changeNameUseCase ChangeIdentityNameUseCase = unregisterUser.ChangeName
	nameChangeRequest := IdentityChangeNameRequest{
		Id:   registrationResponse.Id,
		Name: "Andreas K. Brandhøj",
	}
	_, err = changeNameUseCase(nameChangeRequest)
	if err != nil {
		uow.Rollback()
		log.Fatal("Name change of identity failed")
	}
	log.Println("Identity:", registrationResponse.Id, "changed name to '"+nameChangeRequest.Name+"'")

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
	identity := Visit(events)

	// Snapshot the aggregate
	snapshotData := IdentitySnapshotV1{
		Id:    identity.Id,
		Name:  identity.Name,
		Age:   identity.Age,
		Email: identity.Email,
	}
	snapshot, err := uow.store.Snapshot(Producer, es.SubjectID(identity.Id), snapshotData)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Snapshot", snapshot)

	// Commit with the snapshot
	if err = uow.Commit(); err != nil {
		log.Fatal(err)
	}

	// print the final identity
	log.Println(identity.Id, identity.Name, identity.Age, identity.Email)

	// For this example we sleep here such that you have the time to verify the output
	time.Sleep(5 * time.Second)
}

func Visit(events []es.Event) Identity {
	identity := Identity{}
	for _, event := range events {
		log.Println("Visit:", event)
		switch event.Name {
		case "IdentityRegistered":
			var data IdentityRegistered
			event.Unmarshal(&data)
			identity.Id = data.Id
			identity.Age = data.Age
			identity.Email = data.Email
			identity.Name = data.Name
		case "IdentityChangedName":
			var data IdentityChangedName
			event.Unmarshal(&data)
			identity.Name = data.Name
		case "IdentityChangedAge":
			var data IdentityChangedAge
			event.Unmarshal(&data)
			identity.Age = data.Age
		case "IdentityChangedEmail":
			var data IdentityChangedEmail
			event.Unmarshal(&data)
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

func (unregisteredUser UnregisteredUser) ChangeAge(request IdentityChangeAgeRequest) (IdentityChangeAgeResponse, error) {
	response := IdentityChangeAgeResponse{
		Success: true,
	}
	mediator.Publish(
		es.SubjectID(request.Id),
		IdentityChangedAge{
			Age: request.Age,
		},
	)
	return response, nil
}

func (unregisteredUser UnregisteredUser) ChangeEmail(request IdentityChangeEmailRequest) (IdentityChangeEmailResponse, error) {
	response := IdentityChangeEmailResponse{
		Success: true,
	}
	mediator.Publish(
		es.SubjectID(request.Id),
		IdentityChangedEmail{
			Email: request.Email,
		},
	)
	return response, nil
}

func (uow *UnitOfWork) receiveEvent(subject es.SubjectID, data es.Data) {
	uow.store.Load(
		Producer,
		subject,
		data,
	)
}

func (uow *UnitOfWork) Commit() error {
	events, err := uow.store.Ship()
	if err != nil {
		log.Panicln("UnitOfWork store failed shipping the events")
		return errors.Wrap(err, "UnitOfWork Commiting failed")
	}
	uow.stream.Publish(events)
	return nil
}

func (uow *UnitOfWork) Rollback() error {
	uow.store.Clear()
	return nil
}
