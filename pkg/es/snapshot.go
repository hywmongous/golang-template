package es

import (
	"time"

	"github.com/cockroachdb/errors"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
)

type Snapshot struct {
	Id            Ident
	Producer      ProducerID
	Subject       SubjectID
	Version       Version
	SchemaVersion Version
	Name          Title
	Timestamp     Timestamp
	Data          Data
}

var (
	ErrSnapshotDataIsNil = errors.New("data cannot be nil")
)

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
		Id:            Ident(uuid.New().String()),
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
	Id Ident,
	Producer ProducerID,
	Subject SubjectID,
	Version Version,
	SchemaVersion Version,
	Name Title,
	Timestamp Timestamp,
	Data Data,
) Snapshot {
	return Snapshot{
		Id:            Id,
		Producer:      Producer,
		Subject:       Subject,
		Version:       Version,
		SchemaVersion: SchemaVersion,
		Name:          Name,
		Timestamp:     Timestamp,
		Data:          Data,
	}
}

func nextSnapshotVersion(subject SubjectID, store EventStore) (Version, error) {
	latestSnapshot, err := store.LatestSnapshot(subject)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return InitialSnapshotVersion, nil
	} else if err != nil {
		return InitialSnapshotVersion, errors.Wrap(err, "could not retrieve the latests snapshot")
	}
	return latestSnapshot.Version + 1, nil
}
