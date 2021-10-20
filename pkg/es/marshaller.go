package es

import (
	"encoding/json"

	"github.com/cockroachdb/errors"
)

var (
	ErrDataCouldNotBeMarshalledAsEvent   = errors.New("event could not be json marshalled")
	ErrDataCouldNotBeUnmarshalledAsEvent = errors.New("data byte array could not be json unmarshalled to event")
	ErrMarshallTypeConversionFailed      = errors.New("json marhsalling between types for conversion failed")
)

func (event Event) Marshall() ([]byte, error) {
	data, err := json.Marshal(event)

	return data, errors.Wrap(err, ErrDataCouldNotBeMarshalledAsEvent.Error())
}

func UnmarshalEvent(data []byte) (Event, error) {
	var event Event
	if err := json.Unmarshal(data, &event); err != nil {
		return Event{}, errors.Wrap(err, ErrDataCouldNotBeUnmarshalledAsEvent.Error())
	}

	return event, nil
}

func (event Event) Unmarshal(receiver Data) error {
	return unmarshal(event.Data, receiver)
}

func (snapshot Snapshot) Unmarshal(receiver Data) error {
	return unmarshal(snapshot.Data, receiver)
}

func unmarshal(from interface{}, to interface{}) error {
	bytes, err := json.Marshal(from)
	if err != nil {
		return errors.Wrap(err, ErrMarshallTypeConversionFailed.Error())
	}

	return errors.Wrap(json.Unmarshal(bytes, to), ErrMarshallTypeConversionFailed.Error())
}
