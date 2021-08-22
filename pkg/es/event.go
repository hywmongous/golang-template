package es

import (
	"encoding/json"
	"errors"
	"math"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	merr "github.com/hywmongous/example-service/pkg/errors"
)

type (
	// UUID for the event
	EventId            string
	ProducerID         string
	SubjectID          string
	EventVersion       uint
	EventSchemaVersion uint
	SnapshotVersion    uint
	EventName          string
	EventTimestamp     int64
	EventData          interface{}
)

type Event struct {
	// UUID for the event
	Id EventId

	// The producer of the event
	// This id be the service ID or name
	Producer ProducerID

	// Who is this event regarding?
	// In case of DDD it could be the aggregate root ID
	Subject SubjectID

	// The version of the ID, which is used to
	// sort the events in the created order
	Version EventVersion

	// The version of the event data
	SchemaVersion EventSchemaVersion

	// The snapshot version which this event is under
	SnapshotVersion SnapshotVersion

	// The name of the Event.
	// For instance: "IdentityRegistered"
	// The name can be generated with "CreateEventName"
	Name EventName

	// The time of which the event was created
	Timestamp EventTimestamp

	// The data regarding the event
	// For isntance, if the event is "IdentityRegistered"
	// then the data could be the time of registration
	// and the ID of the registrated identity.
	Data EventData
}

var (
	ErrEventDataIsNil = errors.New("data cannot be nil")
	ErrNoEventData    = errors.New("event data array is length 0")
)

const (
	InitialEventVersion       = EventVersion(0)
	InitialEventSchemaVersion = EventSchemaVersion(0)
	InitialSnapshotVersion    = SnapshotVersion(0)

	BeginningOfTime = EventTimestamp(0)
	EndOfTime       = EventTimestamp(math.MaxInt64)
)

func CreateEvent(
	producer ProducerID,
	subject SubjectID,
	schemaVersion EventSchemaVersion,
	data EventData,
	store EventStore,
) (Event, error) {
	if data == nil {
		return Event{}, merr.CreateFailedInvocation("CreateEvent", ErrEventDataIsNil)
	}

	latestEvent, err := store.Latest(subject)
	if err != nil {
		return Event{}, merr.CreateFailedInvocation("CreateEvent", err)
	}

	return createEvent(producer, subject, latestEvent.Version+1, schemaVersion, data)
}

func CreateEventBatch(
	producer ProducerID,
	subject SubjectID,
	schemaVersion EventSchemaVersion,
	data []EventData,
	store EventStore,
) ([]Event, error) {
	if len(data) == 0 {
		return []Event{}, merr.CreateFailedInvocation("CreateEventBatch", ErrNoEventData)
	}

	latestEvent, err := store.Latest(subject)
	// TODO: Differentiate between errors, it might be the error is caused
	//   by the absence of events and the current event version is the initial
	//   however it could also be a connectivity issue.
	if err != nil {
		return nil, err
	}
	nextEventVersion := latestEvent.Version + 1

	events := make([]Event, len(data))
	for idx, elem := range data {
		event, err := createEvent(
			producer,
			subject,
			nextEventVersion,
			schemaVersion,
			elem,
		)
		if err != nil {
			return []Event{}, merr.CreateFailedInvocation("CreateEventBatch", err)
		}

		events[idx] = event
	}

	return events, nil
}

func createEvent(
	producer ProducerID,
	subject SubjectID,
	version EventVersion,
	schemaVersion EventSchemaVersion,
	data EventData,
) (Event, error) {
	if data == nil {
		return Event{}, merr.CreateFailedInvocation("createEvent", ErrEventDataIsNil)
	}

	return Event{
		Id:              EventId(uuid.New().String()),
		Producer:        producer,
		Subject:         subject,
		Version:         version,
		SchemaVersion:   schemaVersion,
		SnapshotVersion: InitialSnapshotVersion,
		Name:            CreateEventName(data),
		Timestamp:       EventTimestamp(time.Now().Unix()),
		Data:            data,
	}, nil
}

func RecreateEvent(
	id EventId,
	producer ProducerID,
	subject SubjectID,
	version EventVersion,
	name EventName,
	timestamp EventTimestamp,
	data EventData,
) Event {
	return Event{
		Id:        id,
		Producer:  producer,
		Subject:   subject,
		Version:   version,
		Name:      name,
		Timestamp: timestamp,
		Data:      data,
	}
}

func CreateEventName(data EventData) EventName {
	eventType := reflect.TypeOf(data).String()
	eventTypeParts := strings.Split(eventType, ".")
	eventName := eventTypeParts[len(eventTypeParts)-1]
	return EventName(eventName)
}

func (event Event) Marshall() ([]byte, error) {
	return json.Marshal(event)
}

func Unmarshal(data []byte) (Event, error) {
	var event Event
	if err := json.Unmarshal(data, &event); err != nil {
		return Event{}, err
	}
	return event, nil
}

func (event Event) Unmarshal(data interface{}) error {
	bytes, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, data)
}
