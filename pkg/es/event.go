package es

import (
	"time"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

type Event struct {
	// UUID for the event
	ID Ident

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
	ErrEventDataIsNil       = errors.New("data cannot be nil")
	ErrNoEventData          = errors.New("event data array is length 0")
	ErrFindingLatestVersion = errors.New("failure finding the latest version")
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
		return Event{}, errors.Wrap(err, "could not get event next event version")
	}

	snapshotVersion, err := currentSnapshotVersion(subject, store)
	if err != nil {
		return Event{}, errors.Wrap(err, "could not get current snapshot version")
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
		ID:              Ident(uuid.New().String()),
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
		ID:              id,
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

func EmptyEvent() Event {
	return RecreateEvent(
		Ident(""),
		ProducerID(""),
		SubjectID(""),
		Version(0),
		Version(0),
		Version(0),
		Title(""),
		Timestamp(0),
		nil,
	)
}

func nextEventVersion(subject SubjectID, store EventStore) (Version, error) {
	latestEvent, err := store.LatestEvent(subject)
	if errors.Is(err, ErrNoEvents) {
		return InitialEventVersion, nil
	} else if err != nil {
		return InitialEventVersion, errors.Wrap(err, ErrFindingLatestVersion.Error())
	}

	return latestEvent.Version + 1, nil
}

func currentSnapshotVersion(subject SubjectID, store EventStore) (Version, error) {
	latestSnapshot, err := store.LatestSnapshot(subject)
	if errors.Is(err, ErrNoSnapshots) {
		return InitialSnapshotVersion, nil
	} else if err != nil {
		return InitialSnapshotVersion, errors.Wrap(err, ErrFindingLatestVersion.Error())
	}

	return latestSnapshot.Version, nil
}
