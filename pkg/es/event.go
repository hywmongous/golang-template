package es

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	merr "github.com/hywmongous/example-service/pkg/errors"
)

type (
	EventId            string
	ProducerID         string
	SubjectID          string
	EventVersion       uint // TODO: Use this!
	EventSchemaVersion uint
	EventName          string
	EventTimestamp     int64
	EventData          interface{}
	EventDataType      reflect.Type
)

type Event struct {
	Id        EventId
	Producer  ProducerID
	Subject   SubjectID
	Version   EventVersion
	Name      EventName
	Timestamp EventTimestamp
	Data      EventData
}

var (
	EventType = reflect.TypeOf((*Event)(nil))

	ErrEventDataIsNil = errors.New("data cannot be nil")
	ErrNoEventData    = errors.New("event data array is length 0")
)

func CreateEvent(
	producer ProducerID,
	subject SubjectID,
	data EventData,
	store EventStore,
) (Event, error) {
	if data == nil {
		return Event{}, merr.CreateFailedInvocation("CreateEvent", ErrEventDataIsNil)
	}

	latestEvent, err := store.LatestEvent(subject, EventType)
	if err != nil {
		return Event{}, merr.CreateFailedInvocation("CreateEvent", err)
	}

	return createEvent(producer, subject, latestEvent.Version+1, data)
}

func CreateEventBatch(
	producer ProducerID,
	subject SubjectID,
	data []EventData,
	store EventStore,
) ([]Event, error) {
	if len(data) == 0 {
		return []Event{}, merr.CreateFailedInvocation("CreateEventBatch", ErrNoEventData)
	}

	latestEvent, err := store.LatestEvent(subject, EventType)
	if err != nil {
		return []Event{}, merr.CreateFailedInvocation("CreateEventBatch", err)
	}

	events := make([]Event, len(data))
	for idx, elem := range data {
		event, err := createEvent(
			producer,
			subject,
			EventVersion(int(latestEvent.Version)+idx+1),
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
	data EventData,
) (Event, error) {
	if data == nil {
		return Event{}, merr.CreateFailedInvocation("createEvent", ErrEventDataIsNil)
	}

	return Event{
		Id:        EventId(uuid.New().String()),
		Producer:  producer,
		Subject:   subject,
		Version:   version,
		Name:      CreateEventName(data),
		Timestamp: EventTimestamp(time.Now().Unix()),
		Data:      data,
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

func (event Event) Unmarshal(data interface{}) error {
	bytes, _ := json.Marshal(event.Data)
	json.Unmarshal(bytes, data)
	return nil
}
