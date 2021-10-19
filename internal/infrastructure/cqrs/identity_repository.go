package cqrs

import (
	"github.com/cockroachdb/errors"
	"github.com/hywmongous/example-service/internal/domain/authentication"
	"github.com/hywmongous/example-service/pkg/es"
)

type IdentityRepository struct {
	store es.EventStore
}

var (
	ErrVisitForEventFailed       = errors.New("visiting event failed")
	ErrCouldNotFindEntity        = errors.New("could not find entity in event store")
	ErrCouldNotReconstructEntity = errors.New("could not construct entity")
)

func IdentityRepositoryFactory(
	store es.EventStore,
) authentication.Repository {
	return IdentityRepository{
		store: store,
	}
}

func (repository IdentityRepository) FindIdentityByEmail(email string) (authentication.Identity, error) {
	events, err := repository.store.Concerning(es.SubjectID(email))
	if err != nil {
		return authentication.Identity{}, errors.Wrap(err, ErrCouldNotFindEntity.Error())
	}

	model := defaultIdentityModel()
	if err = visitEvents(events, &model); err != nil {
		return authentication.Identity{}, errors.Wrap(err, ErrCouldNotReconstructEntity.Error())
	}

	return authentication.RecreateIdentity(
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
			var data authentication.IdentityRegistered
			if err := event.Unmarshal(&data); err != nil {
				return errors.Wrap(err, ErrVisitForEventFailed.Error())
			}
			model = model.ApplyIdentityRegistered(&data)
		case "IdentityLoggedIn":
			var data authentication.IdentityLoggedIn
			if err := event.Unmarshal(&data); err != nil {
				return errors.Wrap(err, ErrVisitForEventFailed.Error())
			}
			model = model.ApplyIdentityLoggedIn(&data)
		case "IdentityLoggedOut":
			var data authentication.IdentityLoggedOut
			if err := event.Unmarshal(&data); err != nil {
				return errors.Wrap(err, ErrVisitForEventFailed.Error())
			}
			model = model.ApplyIdentityLoggedOut(&data)
		}
	}
	return nil
}
