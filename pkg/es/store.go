package es

type EventStore interface {
	// Immediately sends an Event to the warehouse
	Send(producer ProducerID, subject SubjectID, data []EventData) ([]Event, error)
	// The same as "begin commit"
	Load(producer ProducerID, subject SubjectID, data EventData) (Event, error)
	// The same as removing the current EventData in the commit
	Unload(eventId EventId) (Event, error)
	// Ships the EventData to the Database
	Ship() ([]Event, error)

	// Requests the Events for a specific "Subject"
	// The events are in sorted order with ascending versions
	Concerning(subject SubjectID) ([]Event, error)
	// Requests the Events created by a specific "Producer"
	// The result is in order with ascending versions
	By(producer ProducerID) ([]Event, error)

	// Requests the Events between a specific range of version for the given "Subject"
	// from and to is inclusive meaning:
	// from=1 and to=1 returns Event with version 1
	// from=1 and to=2 returns Events with respective versions 1 and 2
	// The result is in order with ascending versions
	Between(subject SubjectID, from EventVersion, to EventVersion) ([]Event, error)

	// Requests the events for a given a snapshot
	With(subject SubjectID, snapshot SnapshotVersion) ([]Event, error)

	// Requests all events after a point in time
	After(subject SubjectID, pointInTime EventTimestamp) ([]Event, error)
	// Requests the events within a temporal range
	Temporal(subject SubjectID, from EventTimestamp, to EventTimestamp) ([]Event, error)
	// Requests all events before a point in time
	Before(subject SubjectID, pointInTime EventTimestamp) ([]Event, error)

	// Returns the latest event shipped to the database for a given subject
	// This is not temporal based but version based.
	Latest(subject SubjectID) (Event, error)
}
