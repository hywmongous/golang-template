package es

import "encoding/json"

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

func (event Event) Unmarshal(data Data) error {
	bytes, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, data)
}
