package es

type EventStore interface {
	Stock(producer ProducerID, subject SubjectID, data []EventData) ([]Event, error)
	Retrieve(subject SubjectID) ([]Event, error)

	Latest(subject SubjectID) (Event, error)
	CurrentEventVersion(subject SubjectID) EventVersion
}
