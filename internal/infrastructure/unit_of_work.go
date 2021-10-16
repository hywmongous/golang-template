package infrastructure

import (
	"github.com/cockroachdb/errors"
	"github.com/hywmongous/example-service/internal/domain/identity"
	"github.com/hywmongous/example-service/pkg/es"
	"github.com/hywmongous/example-service/pkg/es/kafka"
	"github.com/hywmongous/example-service/pkg/es/mediator"
	"github.com/hywmongous/example-service/pkg/es/mongo"
)

type UnitOfWork struct {
	store  es.EventStore
	stream es.EventStream

	identityRepository identity.Repository
}

var (
	Producer = es.ProducerID("ia")
	Topic    = es.Topic("ia")
)

func (uow *UnitOfWork) IdentityRepository() identity.Repository {
	return uow.identityRepository
}

func MongoStoreFactory() es.EventStore {
	return mongo.CreateMongoEventStore()
}

func KafkaStreamFactory() es.EventStream {
	return kafka.CreateKafkaStream(Topic)
}

func UnitOfWorkFactory(
	store es.EventStore,
	stream es.EventStream,
	identityRepository identity.Repository,
) UnitOfWork {
	uow := UnitOfWork{
		store:              store,
		stream:             stream,
		identityRepository: identityRepository,
	}
	mediator.Listen(uow.receiveEvent)
	return uow
}

func (uow *UnitOfWork) receiveEvent(subject es.SubjectID, data es.Data) {
	if err := uow.store.Load(
		Producer,
		subject,
		data,
	); err != nil {
		panic(err)
	}
}

func (uow *UnitOfWork) Commit() error {
	// events := uow.store.Stage().Events()
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
