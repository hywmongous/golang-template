package cqrs

import (
	"github.com/hywmongous/example-service/internal/domain/identity"
	"github.com/hywmongous/example-service/pkg/es"
)

type IdentityRepository struct {
	store es.EventStore
}

func IdentityRepositoryFactory(
	store es.EventStore,
) identity.Repository {
	return IdentityRepository{
		store: store,
	}
}

func (repository IdentityRepository) FindIdentityByEmail(email string) (identity.Identity, error) {
	events, err := repository.store.Concerning(es.SubjectID(email))
	if err != nil {
		return identity.Identity{}, err
	}

	model := identityModel{}
	if err = visitEvents(events, &model); err != nil {
		return identity.Identity{}, err
	}

	return identity.RecreateIdentity(
		model.id,
		model.email,
		model.password,
		model.sessions,
	), nil
}

func visitEvents(events []es.Event, model readModel) error {
	for _, event := range events {
		switch event.Name {
		case "IdentityRegistered":
			var data identity.IdentityRegistered
			if err := event.Unmarshal(&data); err != nil {
				return err
			}
			model = model.ApplyIdentityRegistered(&data)
		case "IdentityLoggedIn":
			var data identity.IdentityLoggedIn
			if err := event.Unmarshal(&data); err != nil {
				return err
			}
			model = model.ApplyIdentityLoggedIn(&data)
		case "IdentityLoggedOut":
			var data identity.IdentityLoggedOut
			if err := event.Unmarshal(&data); err != nil {
				return err
			}
			model = model.ApplyIdentityLoggedOut(&data)
		}
	}
	return nil
}
