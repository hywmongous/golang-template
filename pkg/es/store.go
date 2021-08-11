package es

type EventStore interface {
	Stock(producer ProducerID, subject SubjectID, data []EventData) error
	Retrieve(subject SubjectID, callback Callback) ([]Event, error)

	LatestEvent(subject SubjectID, dataType EventDataType) (Event, error)
}

type (
	Callback func(event Event) error
)
