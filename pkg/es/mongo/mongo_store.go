package mongo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/cockroachdb/errors"

	"github.com/hywmongous/example-service/pkg/es"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	mongoConnectionAction func(context context.Context, collection *mongo.Collection) error
)

type stage struct {
	events      []es.Event
	hasSnapshot bool
	snapshot    es.Snapshot
}

type MongoEventStore struct {
	stages []stage
}

const (
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

	snapshotIdKey            = "snapshot.id"
	snapshotProducerKey      = "snapshot.producer"
	snapshotSubjectKey       = "snapshot.subject"
	snapshotversionKey       = "snapshot.version"
	snapshotSchemaVersionKey = "snapshot.schemaversion"
	snapshotNameKey          = "snapshot.name"
	snapshotTimestampKey     = "snapshot.timestamp"
	snapshotDataKey          = "snapshot.data"

	mongoLessThan    = "$lt"
	mongoGreaterThan = "$gt"

	mongoAscending  = 1
	mongoDescending = -1
)

var (
	ErrDatabaseNotFound   = errors.New("eventstore database could not be found")
	ErrCollectionNotFound = errors.New("events collection could not be found")
	ErrInsertion          = errors.New("failed inserting one or more documents into collection")
	ErrEventNotFound      = errors.New("event not found")
)

func CreateMongoEventStore() MongoEventStore {
	store := MongoEventStore{
		stages: make([]stage, 1),
	}
	store.stages[0] = createStage()
	return store
}

func createStage() stage {
	return stage{
		events: make([]es.Event, 0),
	}
}

func (store *MongoEventStore) addStage() {
	store.stages = append(store.stages, createStage())
}

func (store *MongoEventStore) stage(event es.Event) {
	last := len(store.stages) - 1
	store.stages[last].events = append(store.stages[last].events, event)
}

func (store *MongoEventStore) clearStage() {
	store.stages = make([]stage, 0)
	store.addStage()
}

func (store *MongoEventStore) unstage(lookup es.Ident) (es.Event, error) {
	for stageIdx, stage := range store.stages {
		for eventIdx, event := range stage.events {
			if event.Id == lookup {
				store.stages[stageIdx].events = append(stage.events[:eventIdx], stage.events[eventIdx+1:]...)
				return event, nil
			}
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

	// create session
	session, err := client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	// begin transaction to ensure rollbacks on errors
	// if err = session.StartTransaction(); err != nil {
	// 	return err
	// }

	// Connect to the collection
	collection, err := store.collection(client, collectionName)
	if err != nil {
		return err
	}

	// Do action encapsulated in the transaction (session)
	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err = action(sc, collection); err != nil {
			// sc.AbortTransaction(sc)
			return err
		}
		// if err = sc.CommitTransaction(sc); err != nil {
		// 	log.Println("Step 9")
		// 	return err
		// }
		return nil
	})

	return err
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

func decodeEvent(
	decoder interface{ Decode(interface{}) error },
	value interface{},
) error {
	var document bson.M
	if err := decoder.Decode(&document); err != nil {
		return err
	}

	eventDocument, ok := document["event"].(bson.M)
	if !ok {
		return errors.New("document does not have event key")
	}

	if err := unmarshalDocument(eventDocument, value); err != nil {
		return err
	}

	return nil
}

func decodeSnapshot(
	decoder interface{ Decode(interface{}) error },
	value interface{},
) error {
	var document bson.M
	if err := decoder.Decode(&document); err != nil {
		return err
	}

	snapshotDocument, ok := document["snapshot"].(bson.M)
	if !ok {
		return errors.New("document does not have snapshot key")
	}

	if err := unmarshalDocument(snapshotDocument, value); err != nil {
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
	// The following is a JSON work around golang mongodb
	// driver does not support decoding of interface{}.
	// This caused issues with the evenData within
	// the event itself. This is however supported by
	// the json package, so for Marshalling we use json
	obj, err := json.Marshal(document)
	if err != nil {
		return err
	}

	if err = json.Unmarshal(obj, &value); err != nil {
		return err
	}

	return nil
}

func (store *MongoEventStore) findOneEvent(filter interface{}, options ...*options.FindOneOptions) (es.Event, error) {
	var resultantEvent es.Event
	action := func(ctx context.Context, collection *mongo.Collection) error {
		result := collection.FindOne(ctx, filter, options...)
		if result == nil {
			return result.Err()
		}

		if conErr := decodeEvent(result, &resultantEvent); conErr != nil {
			return conErr
		}

		return nil
	}
	return resultantEvent, store.connect(action, eventsCollection)
}

func (store *MongoEventStore) findAllEvents(filter interface{}, options ...*options.FindOptions) ([]es.Event, error) {
	var events []es.Event
	action := func(ctx context.Context, collection *mongo.Collection) error {
		cursor, err := collection.Find(ctx, filter, options...)
		if err != nil {
			return err
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var event es.Event
			conErr := decodeEvent(cursor, &event)
			if conErr != nil {
				return conErr
			}
			events = append(events, event)
		}

		return nil
	}

	return events, store.connect(action, eventsCollection)
}

func (store *MongoEventStore) findOneSnapshot(filter interface{}, options ...*options.FindOneOptions) (es.Snapshot, error) {
	var resultantSnapshot es.Snapshot
	action := func(ctx context.Context, collection *mongo.Collection) error {
		result := collection.FindOne(ctx, filter, options...)
		if result == nil {
			return result.Err()
		}

		if conErr := decodeSnapshot(result, &resultantSnapshot); conErr != nil {
			return conErr
		}

		return nil
	}
	return resultantSnapshot, store.connect(action, snapshotsCollection)
}

func (store *MongoEventStore) Send(producer es.ProducerID, subject es.SubjectID, data []es.Data) ([]es.Event, error) {
	events, err := es.CreateEventBatch(producer, subject, es.Version(1), data, store)
	if err != nil {
		return nil, err
	}
	return events, store.sendEvents(events)
}

func (store *MongoEventStore) sendEvents(events []es.Event) error {
	if len(events) == 0 {
		return nil
	}
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

func (store *MongoEventStore) Clear() error {
	store.clearStage()
	return nil
}

func (store *MongoEventStore) Ship() ([]es.Event, error) {
	events := make([]es.Event, 1, 10)
	for _, stage := range store.stages {
		if err := store.sendEvents(stage.events); err != nil {
			return events, errors.Wrap(err, "shipping the events failed")
		}
		events = append(events, stage.events...)

		if stage.hasSnapshot {
			if err := store.sendSnapshot(stage.snapshot); err != nil {
				return events, errors.Wrap(err, "shipping the snapshot failed")
			}
		}
	}
	store.clearStage()
	return events, nil
}

func (store *MongoEventStore) Snapshot(producer es.ProducerID, subject es.SubjectID, data es.Data) (es.Snapshot, error) {
	snapshot, err := es.CreateSnapshot(producer, subject, es.Version(1), data, store)
	if err != nil {
		return es.Snapshot{}, errors.Wrap(err, "Snapshot creation failed")
	}
	store.stages[len(store.stages)-1].snapshot = snapshot
	store.stages[len(store.stages)-1].hasSnapshot = true
	store.addStage() // We create a stage per snapshot
	return snapshot, nil
}

func (store *MongoEventStore) Concerning(subject es.SubjectID) ([]es.Event, error) {
	filter := bson.D{{Key: eventSubjectKey, Value: subject}}
	options := options.Find()
	options.SetSort(bson.D{{Key: eventVersionKey, Value: mongoDescending}})

	events, err := store.findAllEvents(filter, options)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (store *MongoEventStore) By(producer es.ProducerID) ([]es.Event, error) {
	filter := bson.D{{Key: eventProducerKey, Value: producer}}
	options := options.Find()
	options.SetSort(bson.D{{Key: eventVersionKey, Value: mongoDescending}})

	events, err := store.findAllEvents(filter, options)
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
	options.SetSort(bson.D{{Key: eventVersionKey, Value: mongoDescending}})

	events, err := store.findAllEvents(filter, options)
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
	options.SetSort(bson.D{{Key: eventVersionKey, Value: mongoDescending}})

	events, err := store.findAllEvents(filter, options)
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
	options.SetSort(bson.D{{Key: eventVersionKey, Value: mongoDescending}})

	events, err := store.findAllEvents(filter, options)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (store *MongoEventStore) LatestEvent(subject es.SubjectID) (es.Event, error) {
	filter := bson.D{{Key: eventSubjectKey, Value: subject}}
	options := options.FindOne()
	options.SetSort(bson.D{{Key: eventVersionKey, Value: mongoDescending}})

	event, err := store.findOneEvent(filter, options)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return es.Event{}, es.ErrNoEvents
	}

	return event, errors.Wrap(err, "latest event encountered an error")
}

func (store *MongoEventStore) LatestSnapshot(subject es.SubjectID) (es.Snapshot, error) {
	filter := bson.D{
		{Key: snapshotSubjectKey, Value: subject},
	}
	options := options.FindOne()
	options.SetSort(bson.D{{Key: eventSnapShotVersionKey, Value: mongoDescending}})

	snapshot, err := store.findOneSnapshot(filter, options)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return es.Snapshot{}, es.ErrNoSnapshots
	}

	return snapshot, errors.Wrap(err, "latest snapshot encountered an error")
}
