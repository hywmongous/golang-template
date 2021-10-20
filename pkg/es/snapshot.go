package es

import (
	"time"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
)

type Snapshot struct {
	ID            Ident
	Producer      ProducerID
	Subject       SubjectID
	Version       Version
	SchemaVersion Version
	Name          Title
	Timestamp     Timestamp
	Data          Data
}

var ErrSnapshotDataIsNil = errors.New("data cannot be nil")

func CreateSnapshot(
	producer ProducerID,
	subject SubjectID,
	schemaVersion Version,
	data Data,
	store EventStore,
) (Snapshot, error) {
	if data == nil {
		return Snapshot{}, ErrSnapshotDataIsNil
	}

	nextSnapshotVersion, err := nextSnapshotVersion(subject, store)
	if err != nil {
		return Snapshot{}, err
	}

	return Snapshot{
		ID:            Ident(uuid.New().String()),
		Producer:      producer,
		Subject:       subject,
		Version:       nextSnapshotVersion,
		SchemaVersion: schemaVersion,
		Name:          CreateTitleForData(data),
		Timestamp:     Timestamp(time.Now().Unix()),
		Data:          data,
	}, nil
}

func RecreateSnapshot(
	id Ident,
	producer ProducerID,
	subject SubjectID,
	version Version,
	schemaVersion Version,
	name Title,
	timestamp Timestamp,
	data Data,
) Snapshot {
	return Snapshot{
		ID:            id,
		Producer:      producer,
		Subject:       subject,
		Version:       version,
		SchemaVersion: schemaVersion,
		Name:          name,
		Timestamp:     timestamp,
		Data:          data,
	}
}

func EmptySnapshot() Snapshot {
	return RecreateSnapshot(
		Ident(""),
		ProducerID(""),
		SubjectID(""),
		Version(0),
		Version(0),
		Title(""),
		Timestamp(0),
		nil,
	)
}

func nextSnapshotVersion(subject SubjectID, store EventStore) (Version, error) {
	latestSnapshot, err := store.LatestSnapshot(subject)
	if errors.Is(err, ErrNoSnapshots) {
		// We also incremente by one if no snapshots have been made
		// The reason for this is that version = 0 is seen kinda like
		// a snapshot or epoch of all events. meaning the first events
		// will use snapshot version 0
		return InitialSnapshotVersion + 1, nil
	} else if err != nil {
		return InitialSnapshotVersion, errors.Wrap(err, "could not retrieve the latests snapshot")
	}

	return latestSnapshot.Version + 1, nil
}
