package es

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type Event struct {
	// UUID for the event
	Id Ident

	// The producer of the event
	// This id be the service ID or name
	Producer ProducerID

	// Who is this event regarding?
	// In case of DDD it could be the aggregate root ID
	Subject SubjectID

	// The version of the ID, which is used to
	// sort the events in the created order
	Version Version

	// The version of the event data
	SchemaVersion Version

	// The snapshot version which this event is under
	SnapshotVersion Version

	// The name of the Event.
	// For instance: "IdentityRegistered"
	// The name can be generated with "CreateEventName"
	Name Title

	// The time of which the event was created
	Timestamp Timestamp

	// The data regarding the event
	// For isntance, if the event is "IdentityRegistered"
	// then the data could be the time of registration
	// and the ID of the registrated identity.
	Data Data
}

var (
	ErrEventDataIsNil = errors.New("data cannot be nil")
	ErrNoEventData    = errors.New("event data array is length 0")
)

func CreateEvent(
	producer ProducerID,
	subject SubjectID,
	schemaVersion Version,
	data Data,
	store EventStore,
) (Event, error) {
	if data == nil {
		return Event{}, ErrEventDataIsNil
	}

	nextVersion, err := nextEventVersion(subject, store)
	if err != nil {
		return Event{}, err
	}

	snapshotVersion, err := currentSnapshotVersion(subject, store)
	if err != nil {
		return Event{}, err
	}

	return createEvent(
		producer,
		subject,
		nextVersion,
		schemaVersion,
		snapshotVersion,
		data,
	), nil
}

func CreateEventBatch(
	producer ProducerID,
	subject SubjectID,
	schemaVersion Version,
	data []Data,
	store EventStore,
) ([]Event, error) {
	if len(data) == 0 {
		return []Event{}, ErrNoEventData
	}

	nextEventVersion, err := nextEventVersion(subject, store)
	if err != nil {
		return nil, err
	}

	snapshotVersion, err := currentSnapshotVersion(subject, store)
	if err != nil {
		return nil, err
	}

	events := make([]Event, len(data))
	for idx, elem := range data {
		events[idx] = createEvent(
			producer,
			subject,
			nextEventVersion,
			schemaVersion,
			snapshotVersion,
			elem,
		)
	}

	return events, nil
}

func createEvent(
	producer ProducerID,
	subject SubjectID,
	version Version,
	schemaVersion Version,
	snapshotVersion Version,
	data Data,
) Event {
	return Event{
		Id:              Ident(uuid.New().String()),
		Producer:        producer,
		Subject:         subject,
		Version:         version,
		SchemaVersion:   schemaVersion,
		SnapshotVersion: snapshotVersion,
		Name:            CreateTitleForData(data),
		Timestamp:       Timestamp(time.Now().Unix()),
		Data:            data,
	}
}

func RecreateEvent(
	id Ident,
	producer ProducerID,
	subject SubjectID,
	version Version,
	schemaVersion Version,
	snapshotVersion Version,
	name Title,
	timestamp Timestamp,
	data Data,
) Event {
	return Event{
		Id:              id,
		Producer:        producer,
		Subject:         subject,
		Version:         version,
		SchemaVersion:   schemaVersion,
		SnapshotVersion: snapshotVersion,
		Name:            name,
		Timestamp:       timestamp,
		Data:            data,
	}
}

func nextEventVersion(subject SubjectID, store EventStore) (Version, error) {
	latestEvent, err := store.LatestEvent(subject)
	if err == mongo.ErrNoDocuments {
		return InitialEventVersion, nil
	} else if err != nil {
		return InitialEventVersion, err
	}
	return latestEvent.Version + 1, nil
}

func currentSnapshotVersion(subject SubjectID, store EventStore) (Version, error) {
	latestSnapshot, err := store.LatestSnapshot(subject)
	if err == mongo.ErrNoDocuments {
		return InitialSnapshotVersion, nil
	} else if err != nil {
		return InitialSnapshotVersion, err
	}
	return latestSnapshot.Version, nil
}
