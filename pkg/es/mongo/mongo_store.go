package mongo

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/hywmongous/example-service/pkg/es"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	mongoConnectionAction func(context context.Context, collection *mongo.Collection) error
)

type MongoEventStore struct {
	commit []es.Event
}

const (
	commitCapacity  = 1 // We append each time an event is added
	timeoutDuration = 10 * time.Second

	databaseName        = "eventstore"
	eventsCollection    = "events"
	snapshotsCollection = "snapshots"
)

const (
	eventIdKey              = "event.id"
	eventProducerKey        = "event.producer"
	eventSubjectKey         = "event.subject"
	eventVersionKey         = "event.version"
	eventSchemaVersionKey   = "event.schemaversion"
	eventSnapShotVersionKey = "event.snapshotversion"
	eventNameKey            = "event.name"
	eventTimestampKey       = "event.timestamp"
	eventDataKey            = "event.data"

	mongoLessThan    = "$lt"
	mongoGreaterThan = "$gt"
)

var (
	ErrDatabaseNotFound   = errors.New("eventstore database could not be found")
	ErrCollectionNotFound = errors.New("events collection could not be found")
	ErrInsertion          = errors.New("failed inserting one or more documents into collection")
	ErrEventNotFound      = errors.New("event not found")
)

func CreateMongoEventStore() MongoEventStore {
	return MongoEventStore{
		commit: make([]es.Event, 0, commitCapacity),
	}
}

func (store *MongoEventStore) Commit() []es.Event {
	return store.commit
}

func (store *MongoEventStore) stage(event es.Event) {
	store.commit = append(store.commit, event)
}

func (store *MongoEventStore) clearStage() {
	store.commit = make([]es.Event, 0, commitCapacity)
}

func (store *MongoEventStore) unstage(lookup es.Ident) (es.Event, error) {
	for idx, event := range store.commit {
		if event.Id == lookup {
			store.commit = append(store.commit[:idx], store.commit[idx+1:]...)
			return event, nil
		}
	}
	return es.Event{}, ErrEventNotFound
}

func (store *MongoEventStore) collection(client *mongo.Client, collectionName string) (*mongo.Collection, error) {
	// Establish database connection
	database := client.Database(databaseName)
	if database == nil {
		return nil, ErrDatabaseNotFound
	}

	// Establish collection connection
	collection := database.Collection(collectionName)
	if collection == nil {
		return nil, ErrCollectionNotFound
	}

	return collection, nil
}

func (store *MongoEventStore) connect(action mongoConnectionAction, collectionName string) error {
	options := options.Client()
	uri := options.ApplyURI("mongodb://root:root@ia_mongo:27017")

	// Client construction
	client, err := mongo.NewClient(uri)
	if err != nil {
		return err
	}

	// Create the context
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	// Construct the connected client
	if err = client.Connect(ctx); err != nil {
		return err
	}
	defer client.Disconnect(ctx)

	// Connect to the collection
	collection, err := store.collection(client, collectionName)
	if err != nil {
		return err
	}

	return action(ctx, collection)
}

func (store *MongoEventStore) insertManyDocuments(documents []interface{}, collectionName string) error {
	action := func(ctx context.Context, collection *mongo.Collection) error {
		_, err := collection.InsertMany(ctx, documents)
		return err
	}

	return store.connect(action, collectionName)
}

func (store *MongoEventStore) insertDocument(document interface{}, collectionName string) error {
	action := func(ctx context.Context, collection *mongo.Collection) error {
		_, err := collection.InsertOne(ctx, document)
		return err
	}

	return store.connect(action, collectionName)
}

func decode(
	decoder interface{ Decode(interface{}) error },
	value interface{},
) error {
	var document bson.M
	if err := decoder.Decode(&document); err != nil {
		return err
	}

	if err := unmarshalDocument(document, value); err != nil {
		return err
	}

	return nil
}

func marshallEventDocument(event es.Event) interface{} {
	return bson.D{{
		Key:   "event",
		Value: event,
	}}
}

func marshallEventDocuments(events []es.Event) []interface{} {
	documents := make([]interface{}, len(events))
	for idx, event := range events {
		documents[idx] = marshallEventDocument(event)
	}
	return documents
}

func marshallSnapshotDocument(snapshot es.Snapshot) interface{} {
	return bson.D{{
		Key:   "snapshot",
		Value: snapshot,
	}}
}

func unmarshalDocument(document bson.M, value interface{}) error {
	eventDocument, ok := document["event"].(bson.M)
	if !ok {
		return errors.New("document with key 'event' could not be converted to a bson map")
	}

	// The following is a JSON work around golang mongodb
	// driver does not support decoding of interface{}.
	// This caused issues with the evenData within
	// the event itself. This is however supported by
	// the json package, so for Marshalling we use json
	obj, err := json.Marshal(eventDocument)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(obj, &value); err != nil {
		return err
	}

	return nil
}

func (store *MongoEventStore) findOne(collectionName string, filter interface{}, options ...*options.FindOneOptions) (es.Event, error) {
	var resultantEvent es.Event
	action := func(ctx context.Context, collection *mongo.Collection) error {
		result := collection.FindOne(ctx, filter, options...)
		if result == nil {
			return result.Err()
		}

		if conErr := decode(result, &resultantEvent); conErr != nil {
			return conErr
		}

		return nil
	}
	return resultantEvent, store.connect(action, collectionName)
}

func (store *MongoEventStore) findAll(collectionName string, filter interface{}, options ...*options.FindOptions) ([]es.Event, error) {
	var events []es.Event
	action := func(ctx context.Context, collection *mongo.Collection) error {
		cursor, err := collection.Find(ctx, filter, options...)
		if err != nil {
			return err
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var event es.Event
			conErr := decode(cursor, &event)
			if conErr != nil {
				return conErr
			}
			events = append(events, event)
		}

		return nil
	}

	return events, store.connect(action, collectionName)
}

func (store *MongoEventStore) Send(producer es.ProducerID, subject es.SubjectID, data []es.Data) ([]es.Event, error) {
	events, err := es.CreateEventBatch(producer, subject, es.Version(1), data, store)
	if err != nil {
		return nil, err
	}
	return events, store.sendEvents(events)
}

func (store *MongoEventStore) sendEvents(events []es.Event) error {
	documents := marshallEventDocuments(events)
	return store.insertManyDocuments(documents, eventsCollection)
}

func (store *MongoEventStore) sendSnapshot(snapshot es.Snapshot) error {
	document := marshallSnapshotDocument(snapshot)
	return store.insertDocument(document, snapshotsCollection)
}

func (store *MongoEventStore) Load(producer es.ProducerID, subject es.SubjectID, data es.Data) (es.Event, error) {
	event, err := es.CreateEvent(producer, subject, es.Version(1), data, store)
	if err != nil {
		return event, err
	}
	store.stage(event)
	return event, nil
}

func (store *MongoEventStore) Unload(lookup es.Ident) (es.Event, error) {
	return store.unstage(lookup)
}

func (store *MongoEventStore) Ship() ([]es.Event, error) {
	if err := store.sendEvents(store.commit); err != nil {
		return nil, err
	}
	shipment := store.commit
	store.clearStage()
	return shipment, nil
}

func (store *MongoEventStore) Snapshot(producer es.ProducerID, subject es.SubjectID, data es.Data) (es.Snapshot, error) {
	snapshot, err := es.CreateSnapshot(producer, subject, es.Version(1), data, store)
	if err != nil {
		return es.Snapshot{}, err
	}
	return snapshot, store.sendSnapshot(snapshot)
}

func (store *MongoEventStore) Concerning(subject es.SubjectID) ([]es.Event, error) {
	filter := bson.D{{Key: eventSubjectKey, Value: subject}}
	options := options.Find()
	options.SetSort(bson.D{{Key: eventVersionKey, Value: -1}})

	events, err := store.findAll(eventsCollection, filter, options)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (store *MongoEventStore) By(producer es.ProducerID) ([]es.Event, error) {
	filter := bson.D{{Key: eventProducerKey, Value: producer}}
	options := options.Find()
	options.SetSort(bson.D{{Key: eventVersionKey, Value: -1}})

	events, err := store.findAll(eventsCollection, filter, options)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (store *MongoEventStore) Between(subject es.SubjectID, from es.Version, to es.Version) ([]es.Event, error) {
	filter := bson.D{
		{Key: eventSubjectKey, Value: subject},
		{Key: eventVersionKey, Value: bson.D{
			{Key: mongoLessThan, Value: to},
			{Key: mongoGreaterThan, Value: from},
		}},
	}
	options := options.Find()
	options.SetSort(bson.D{{Key: eventVersionKey, Value: -1}})

	events, err := store.findAll(eventsCollection, filter, options)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (store *MongoEventStore) With(subject es.SubjectID, snapshot es.Version) ([]es.Event, error) {
	filter := bson.D{
		{Key: eventSubjectKey, Value: subject},
		{Key: eventSnapShotVersionKey, Value: snapshot},
	}
	options := options.Find()
	options.SetSort(bson.D{{Key: eventVersionKey, Value: -1}})

	events, err := store.findAll(eventsCollection, filter, options)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (store *MongoEventStore) After(subject es.SubjectID, pointInTime es.Timestamp) ([]es.Event, error) {
	return store.Temporal(subject, pointInTime, es.EndOfTime)
}

func (store *MongoEventStore) Before(subject es.SubjectID, pointInTime es.Timestamp) ([]es.Event, error) {
	return store.Temporal(subject, pointInTime, es.BeginningOfTime)
}

func (store *MongoEventStore) Temporal(subject es.SubjectID, from es.Timestamp, to es.Timestamp) ([]es.Event, error) {
	filter := bson.D{
		{Key: eventSubjectKey, Value: subject},
		{Key: eventTimestampKey, Value: bson.D{
			{Key: mongoLessThan, Value: to},
			{Key: mongoGreaterThan, Value: from},
		}},
	}
	options := options.Find()
	options.SetSort(bson.D{{Key: eventVersionKey, Value: -1}})

	events, err := store.findAll(eventsCollection, filter, options)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (store *MongoEventStore) LatestEvent(subject es.SubjectID) (es.Event, error) {
	filter := bson.D{{Key: eventSubjectKey, Value: subject}}
	options := options.FindOne()
	options.SetSort(bson.D{{Key: eventVersionKey, Value: -1}})

	event, err := store.findOne(eventsCollection, filter, options)
	if err != nil {
		return es.Event{}, err
	}
	return event, nil
}

func (store *MongoEventStore) LatestSnapshot(subject es.SubjectID) (es.Snapshot, error) {
	return es.Snapshot{}, nil
}
