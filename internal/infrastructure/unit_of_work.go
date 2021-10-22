package infrastructure

import (
	"github.com/cockroachdb/errors"
	"github.com/hywmongous/example-service/internal/domain/authentication"
	"github.com/hywmongous/example-service/pkg/es"
	"github.com/hywmongous/example-service/pkg/es/kafka"
	"github.com/hywmongous/example-service/pkg/es/mediator"
	"github.com/hywmongous/example-service/pkg/es/mongo"
)

type UnitOfWork struct {
	store    es.EventStore
	stream   es.EventStream
	mediator *mediator.Mediator

	identityRepository authentication.Repository
}

const (
	producer = es.ProducerID("ia")
	topic    = es.Topic("ia")
)

var ErrEmptyCommit = errors.New("attempting to commit an empty stage")

func (uow *UnitOfWork) IdentityRepository() authentication.Repository {
	return uow.identityRepository
}

func MongoStoreFactory() es.EventStore {
	return mongo.CreateMongoEventStore()
}

func KafkaStreamFactory() es.EventStream {
	return kafka.CreateKafkaStream(topic)
}

func UnitOfWorkFactory(
	store es.EventStore,
	stream es.EventStream,
	mediator *mediator.Mediator,
	identityRepository authentication.Repository,
) UnitOfWork {
	uow := UnitOfWork{
		store:              store,
		stream:             stream,
		mediator:           mediator,
		identityRepository: identityRepository,
	}

	mediator.Listen(uow.receiveEvent)

	return uow
}

func (uow *UnitOfWork) receiveEvent(subject es.SubjectID, data es.Data) {
	if err := uow.store.Load(
		producer,
		subject,
		data,
	); err != nil {
		panic(err)
	}
}

func (uow *UnitOfWork) Commit() error {
	events := uow.store.Stage().Events()
	if len(events) == 0 {
		return ErrEmptyCommit
	}

	if err := uow.store.Ship(); err != nil {
		return errors.Wrap(err, "UnitOfWork store failed shipping the events")
	}
	// if err := uow.stream.Publish(events); err != nil {
	// 	return errors.Wrap(err, "UnitOfWork stream failed publishing the events")
	// }
	return nil
}

func (uow *UnitOfWork) Clear() {
	uow.store.Clear()
}

func (uow *UnitOfWork) Mediator() *mediator.Mediator {
	return uow.mediator
}
