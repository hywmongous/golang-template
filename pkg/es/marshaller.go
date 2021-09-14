package es

import (
	"encoding/json"

	"github.com/cockroachdb/errors"
)

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

func (event Event) Unmarshal(receiver Data) error {
	return unmarshal(event.Data, receiver)
}

func (snapshot Snapshot) Unmarshal(receiver Data) error {
	return unmarshal(snapshot.Data, receiver)
}

func unmarshal(from interface{}, to interface{}) error {
	bytes, err := json.Marshal(from)
	if err != nil {
		return errors.Wrap(err, "unmarshalling within es failed")
	}
	return json.Unmarshal(bytes, to)
}
