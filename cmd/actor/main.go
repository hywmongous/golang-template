package main

import (
	"log"

	"github.com/hywmongous/example-service/internal/application/actors"
	"github.com/hywmongous/example-service/internal/application/usecases"
	"github.com/hywmongous/example-service/internal/domain/identity"
)

func main() {
	scenario_registraion()
}

func scenario_registraion() {
	actor := actors.CreateUnregisteredIdentity()

	email, _ := identity.CreateEmail("andreasbrandhoej@hotmail.com")
	password, _ := identity.CreatePassword("password")
	registrationRequest := usecases.RegisterIdentityRequest{
		Email:    email,
		Password: password,
	}

	registrationResponse, err := actor.Register(registrationRequest)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Identity '", registrationResponse.IdentityID, "' was registered")
}
