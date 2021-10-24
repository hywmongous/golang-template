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
	mediator *mediator.Mediator
}

func main() {
	// create UOW
	uow := UnitOfWork{
		events: make([]es.Data, 0),
	}

	// Create actor
	actor := &UnregisteredUser{
		mediator: mediator.Create(),
	}

	// Do "Register Identity" usecase by the actor
	age := 22
	request := RegisterIdentityRequest{
		name:  "Andreas",
		age:   age,
		email: "andreasbrandhoej@hotmail.com",
	}

	var registrationUseCase RegisterIdentityUseCase = actor.RegisterIdentity

	response, err := registrationUseCase(request)
	if err != nil {
		log.Println(err)
		uow.Rollback()
	}

	if response.success {
		log.Println("Registration was successful")
	} else {
		log.Println(err)
		uow.Rollback()
	}

	// Commit changes
	uow.Commit()
}

func (user *UnregisteredUser) RegisterIdentity(request RegisterIdentityRequest) (RegisterIdentityResponse, error) {
	response := RegisterIdentityResponse{
		success: true,
	}

	user.mediator.Publish(
		es.SubjectID("Me"),
		IdentityRegistered{
			time: es.Timestamp(time.Now().Unix()),
			name: request.name,
		},
	)

	return response, nil
}

func (uow *UnitOfWork) Commit() {
	for _, event := range uow.events {
		log.Println("Committing:", event)
	}

	uow.events = make([]es.Data, 0)
}

func (uow *UnitOfWork) Rollback() {
	uow.events = make([]es.Data, 0)
}
