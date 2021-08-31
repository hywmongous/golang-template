package main

import (
	"log"
	"time"

	"github.com/hywmongous/example-service/pkg/es"
	"github.com/hywmongous/example-service/pkg/es/mediator"
)

type IdentityRegistered struct {
	time es.Timestamp
	name string
}

type UnitOfWork struct {
	events []es.Data
}

type RegisterIdentityRequest struct {
	name  string
	age   int
	email string
}

type RegisterIdentityResponse struct {
	success bool
}

type RegisterIdentityUseCase func(request RegisterIdentityRequest) (RegisterIdentityResponse, error)

type UnregisteredUser struct {
	Register RegisterIdentityUseCase
}

func main() {
	// create UOW
	uow := UnitOfWork{
		events: make([]es.Data, 0),
	}
	mediator.Listen(uow.receiver)

	// Create actor
	actor := UnregisteredUser{
		Register: RegisterIdentity,
	}

	// Do "Register Identity" usecase by the actor
	request := RegisterIdentityRequest{
		name:  "Andreas",
		age:   22,
		email: "andreasbrandhoej@hotmail.com",
	}
	response, err := actor.Register(request)
	if err != nil {
		log.Panic(err)
		uow.Rollback()
	}

	if response.success {
		log.Println("Registration was successfull")
	} else {
		log.Panic("Registration failed")
		uow.Rollback()
	}

	// Commit changes
	if err = uow.Commit(); err != nil {
		log.Panic(err)
		uow.Rollback()
	}
}

func RegisterIdentity(request RegisterIdentityRequest) (RegisterIdentityResponse, error) {
	response := RegisterIdentityResponse{
		success: true,
	}
	mediator.Publish(
		es.SubjectID("Me"),
		IdentityRegistered{
			time: es.Timestamp(time.Now().Unix()),
			name: request.name,
		},
	)
	return response, nil
}

func (uow *UnitOfWork) Commit() error {
	for _, event := range uow.events {
		log.Println("Committing:", event)
	}
	uow.events = make([]es.Data, 0)
	return nil
}

func (uow *UnitOfWork) Rollback() error {
	uow.events = make([]es.Data, 0)
	return nil
}

func (uow *UnitOfWork) receiver(subject es.SubjectID, data es.Data) {
	uow.events = append(uow.events, data)
}
