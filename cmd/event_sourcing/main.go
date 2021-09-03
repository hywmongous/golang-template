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

// Actors
type UnregisteredUser struct{}

// Use cases
type IdentityRegistrationRequest struct {
	Name  string
	Age   int
	Email string
}
type IdentityRegistrationResponse struct {
	id string
}
type RegisterIdentityUseCase func(request IdentityRegistrationRequest) (IdentityRegistrationResponse, error)

// Events
type IdentityRegistered struct {
	Id    string
	Name  string
	Age   int
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
	request := IdentityRegistrationRequest{
		Name:  "Andreas",
		Age:   22,
		Email: "andreasbrandhoej@hotmail.com",
	}
	response, err := registrationUseCase(request)
	if err != nil {
		uow.Rollback()
		log.Panicln("Registration use case failed")
	}
	log.Println("Identity:", response.id, "was registered")

	// Commit changes
	if err = uow.Commit(); err != nil {
		uow.Rollback()
		log.Panicln(err)
	}

	// For this example we wait a bit after the registration
	time.Sleep(2 * time.Second)

	// Query all the events for the aggregate
	events, err := uow.store.Concerning(es.SubjectID(response.id))
	if err != nil {
		log.Panic(err)
	}

	// Construct the aggregate
	identity := Visit(events)

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
			VisitIdentityRegistered(&identity, data)
		}
	}
	return identity
}

func VisitIdentityRegistered(identity *Identity, event IdentityRegistered) {
	identity.Id = event.Id
	identity.Age = event.Age
	identity.Email = event.Email
	identity.Name = event.Name
}

func (unregisteredUser UnregisteredUser) Register(request IdentityRegistrationRequest) (IdentityRegistrationResponse, error) {
	response := IdentityRegistrationResponse{
		id: uuid.New().String(),
	}
	mediator.Publish(
		es.SubjectID(response.id),
		IdentityRegistered{
			Id:    response.id,
			Name:  request.Name,
			Age:   request.Age,
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
