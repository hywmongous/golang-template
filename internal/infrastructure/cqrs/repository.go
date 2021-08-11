package cqrs

import (
	"errors"

	"github.com/hywmongous/example-service/internal/domain/identity"
	"github.com/hywmongous/example-service/internal/infrastructure/cqrs/visitors"
	merr "github.com/hywmongous/example-service/pkg/errors"
	"github.com/hywmongous/example-service/pkg/es"
)

type IdentityRepository struct {
	eventStore es.EventStore
	visitor    visitors.AggregateVisitor
}

var (
	ErrConsumerIsNil = errors.New("consumer cannot be nil")
)

func CreateIdentityRepository(eventStore es.EventStore) (IdentityRepository, error) {
	if eventStore == nil {
		return IdentityRepository{}, merr.CreateInvalidInputError("CreateIdentityRepository", "consumer", ErrConsumerIsNil)
	}

	return IdentityRepository{
		eventStore: eventStore,
	}, nil
}

func (repository IdentityRepository) ReadAllIdentities() ([]identity.Identity, error) {
	return nil, nil
}

func (repository IdentityRepository) ReadIdentityById(id identity.IdentityID) (identity.Identity, error) {
	producerID := es.ProducerID(id)
	events, err := repository.eventStore.Load(producerID)
	if err != nil {
		return identity.Identity{}, err
	}

	for _, event := range events {
		print(event)
	}
	return identity.Identity{}, nil
}

func (repository IdentityRepository) UpdateIdentity(identity identity.Identity) error {
	return nil
}

func (repository IdentityRepository) DeleteIdentityById(id identity.IdentityID) error {
	return nil
}
