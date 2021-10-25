package mongo

import (
	"context"
	"encoding/json"
	"log"
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

type EventStore struct {
	stage            es.Stage
	insertionHistory map[string][]interface{}
}

const (
	timeoutDuration = 10 * time.Second

	databaseName        = "eventstore"
	eventsCollection    = "events"
	snapshotsCollection = "snapshots"
)

const (
	documentIDKey = "_id"

	// The commented constants are kept to display document structure.
	// eventIdKey              = "event.id".
	eventProducerKey = "event.producer"
	eventSubjectKey  = "event.subject"
	eventVersionKey  = "event.version"
	// eventSchemaVersionKey   = "event.schemaversion".
	eventSnapShotVersionKey = "event.snapshotversion"
	// eventNameKey            = "event.name".
	eventTimestampKey = "event.timestamp"
	// eventDataKey            = "event.data".

	// snapshotIdKey            = "snapshot.id"
	// snapshotProducerKey      = "snapshot.producer".
	snapshotSubjectKey = "snapshot.subject"
	snapshotVersionKey = "snapshot.version"
	// snapshotSchemaVersionKey = "snapshot.schemaversion"
	// snapshotNameKey          = "snapshot.name"
	// snapshotTimestampKey     = "snapshot.timestamp"
	// snapshotDataKey          = "snapshot.data".

	mongoLessThan    = "$lt"
	mongoGreaterThan = "$gt"
	mongoIn          = "$in"

	mongoAscending  = 1
	mongoDescending = -1
)

var (
	ErrDatabaseNotFound                          = errors.New("eventstore database could not be found")
	ErrCollectionNotFound                        = errors.New("events collection could not be found")
	ErrInsertion                                 = errors.New("failed inserting one or more documents into collection")
	ErrEventNotFound                             = errors.New("event not found")
	ErrMissingEventKey                           = errors.New("document does not have event key")
	ErrMissingSnapshotKey                        = errors.New("document does not have snapshot key")
	ErrStageOutOfSync                            = errors.New("stage is out of sync with remote")
	ErrRollbackFailed                            = errors.New("rollback deletions failed")
	ErrEventCreationFailedOnLoad                 = errors.New("event could not be created and loaded")
	ErrEventBatchCreationFailedOnSend            = errors.New("event batch could not be created and sent")
	ErrCouldNotFindSnapshots                     = errors.New("snapshots could not be found in database")
	ErrCouldNotFindEvents                        = errors.New("events could not be found in database")
	ErrCouldNotConnectToEventStore               = errors.New("could not connect to event store")
	ErrEventCouldNotBeDecoded                    = errors.New("snapshot could not be decoded")
	ErrCustomEventDecoder                        = errors.New("custom snapshot decoder could not decode snapshot")
	ErrEventFailedJSONMarshalling                = errors.New("event could not be json marshalled")
	ErrEventFailedJSONUnmarshalling              = errors.New("event could not be json unmarshalled")
	ErrSnapshotCouldNotBeDecoded                 = errors.New("snapshot could not be decoded")
	ErrCustomSnapshotDecoder                     = errors.New("custom snapshot decoder could not decode snapshot")
	ErrSnapshotFailedJSONMarshalling             = errors.New("event could not be json marshalled")
	ErrSnapshotFailedJSONUnmarshalling           = errors.New("event could not be json unmarshalled")
	ErrMongoDocumentDeletionFailed               = errors.New("deleting documents failed")
	ErrMongoDocumentInsertionFailed              = errors.New("inserting document failed")
	ErrMongoClientCouldNotBeCreated              = errors.New("creating new mongo client failed")
	ErrMongoClientCouldNotConnect                = errors.New("mongo client could not connect")
	ErrMongoClientCouldNotCreateSession          = errors.New("mongo client could not create session")
	ErrMongoClientCouldNotConnectionToCollection = errors.New("mongo client could not connect to collection")
	ErrMongoClientCouldNotPerformAction          = errors.New("mongo client could not perform action")
	ErrMongoClientCouldNotDisconnect             = errors.New("mongo client failed disconnecting")
)

func CreateMongoEventStore() *EventStore {
	return &EventStore{
		stage:            es.CreateStage(),
		insertionHistory: make(map[string][]interface{}),
	}
}

func (store *EventStore) Stage() es.Stage {
	return store.stage
}

func (store *EventStore) collection(client *mongo.Client, collectionName string) (*mongo.Collection, error) {
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

func (store *EventStore) connect(action mongoConnectionAction, collectionName string) error {
	options := options.Client()
	uri := options.ApplyURI("mongodb://root:root@ia_mongo:27017")

	// Client construction
	client, err := mongo.NewClient(uri)
	if err != nil {
		return errors.Wrap(err, ErrMongoClientCouldNotBeCreated.Error())
	}

	// Create the context
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	// Construct the connected client
	if err = client.Connect(ctx); err != nil {
		return errors.Wrap(err, ErrMongoClientCouldNotConnect.Error())
	}

	// create session
	session, err := client.StartSession()
	if err != nil {
		return errors.Wrap(err, ErrMongoClientCouldNotCreateSession.Error())
	}

	defer session.EndSession(ctx)

	// begin transaction to ensure rollbacks on errors
	// if err = session.StartTransaction(); err != nil {
	// 	return err
	// }

	// Connect to the collection
	collection, err := store.collection(client, collectionName)
	if err != nil {
		return errors.Wrap(err, ErrMongoClientCouldNotConnectionToCollection.Error())
	}

	// Do action encapsulated in the transaction (session)
	err = mongo.WithSession(ctx, session, func(sc mongo.SessionContext) error {
		if err = action(sc, collection); err != nil {
			// sc.AbortTransaction(sc)
			return errors.Wrap(err, ErrMongoClientCouldNotPerformAction.Error())
		}
		// if err = sc.CommitTransaction(sc); err != nil {
		// 	return err
		// }
		return nil
	})

	return errors.Wrap(client.Disconnect(ctx), ErrMongoClientCouldNotDisconnect.Error())
}

func (store *EventStore) findOneEvent(filter interface{}, options ...*options.FindOneOptions) (es.Event, error) {
	var resultantEvent es.Event

	action := func(ctx context.Context, collection *mongo.Collection) error {
		result := collection.FindOne(ctx, filter, options...)
		if result.Err() != nil {
			return errors.Wrap(result.Err(), ErrCouldNotFindEvents.Error())
		}

		if err := decodeEvent(result, &resultantEvent); err != nil {
			return errors.Wrap(err, ErrEventCouldNotBeDecoded.Error())
		}

		return nil
	}

	return resultantEvent, store.connect(action, eventsCollection)
}

func (store *EventStore) findAllEvents(filter interface{}, options ...*options.FindOptions) ([]es.Event, error) {
	var events []es.Event

	action := func(ctx context.Context, collection *mongo.Collection) error {
		cursor, err := collection.Find(ctx, filter, options...)
		if err != nil {
			return errors.Wrap(err, ErrCouldNotFindEvents.Error())
		}
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var event es.Event

			err := decodeEvent(cursor, &event)
			if err != nil {
				return errors.Wrap(err, ErrEventCouldNotBeDecoded.Error())
			}

			events = append(events, event)
		}

		return nil
	}

	return events, store.connect(action, eventsCollection)
}

func (store *EventStore) addToInsertionHistory(collectionName string, insertionIDs ...interface{}) {
	if _, found := store.insertionHistory[collectionName]; !found {
		store.insertionHistory[collectionName] = make([]interface{}, 0)
	}

	store.insertionHistory[collectionName] = append(store.insertionHistory[collectionName], insertionIDs...)
}

func (store *EventStore) clearInsertionHistory() {
	store.insertionHistory = make(map[string][]interface{})
}

func (store *EventStore) insertManyDocuments(documents []interface{}, collectionName string) error {
	action := func(ctx context.Context, collection *mongo.Collection) error {
		results, err := collection.InsertMany(ctx, documents)
		if err == nil {
			store.addToInsertionHistory(collectionName, results.InsertedIDs...)
		}

		return errors.Wrap(err, ErrMongoDocumentInsertionFailed.Error())
	}

	return store.connect(action, collectionName)
}

func (store *EventStore) insertDocument(document interface{}, collectionName string) error {
	action := func(ctx context.Context, collection *mongo.Collection) error {
		result, err := collection.InsertOne(ctx, document)
		if err == nil {
			store.addToInsertionHistory(collectionName, result.InsertedID)
		}

		return errors.Wrap(err, ErrMongoDocumentInsertionFailed.Error())
	}

	return store.connect(action, collectionName)
}

func (store *EventStore) deleteManyDocument(
	filter interface{},
	collectionName string,
	options ...*options.DeleteOptions,
) error {
	action := func(ctx context.Context, collection *mongo.Collection) error {
		_, err := collection.DeleteMany(ctx, filter, options...)

		return errors.Wrap(err, ErrMongoDocumentDeletionFailed.Error())
	}

	return store.connect(action, collectionName)
}

// func (store *MongoEventStore) deleteDocument(
//	filter interface{},
//	collectionName string,
//	options ...*options.DeleteOptions
// ) error {
// 	action := func(ctx context.Context, collection *mongo.Collection) error {
// 		_, err := collection.DeleteOne(ctx, filter, options...)
// 		return err
// 	}
// 	return store.connect(action, collectionName)
// }

func (store *EventStore) rollbackInsertions() error {
	for collectionName, ids := range store.insertionHistory {
		filter := bson.D{
			{Key: documentIDKey, Value: bson.D{
				{Key: mongoIn, Value: ids},
			}},
		}
		options := options.Delete()

		if err := store.deleteManyDocument(filter, collectionName, options); err != nil {
			return errors.Wrap(err, "rollback deletion of documents failed")
		}
	}

	return nil
}

func decodeEvent(
	decoder interface{ Decode(interface{}) error },
	value interface{},
) error {
	var document bson.M
	if err := decoder.Decode(&document); err != nil {
		return errors.Wrap(err, ErrCustomEventDecoder.Error())
	}

	eventDocument, ok := document["event"].(bson.M)
	if !ok {
		return ErrMissingEventKey
	}

	if err := unmarshalDocument(eventDocument, value); err != nil {
		return errors.Wrap(err, ErrEventFailedJSONUnmarshalling.Error())
	}

	return nil
}

func decodeSnapshot(
	decoder interface{ Decode(interface{}) error },
	value interface{},
) error {
	var document bson.M
	if err := decoder.Decode(&document); err != nil {
		return errors.Wrap(err, ErrCustomSnapshotDecoder.Error())
	}

	snapshotDocument, ok := document["snapshot"].(bson.M)
	if !ok {
		return ErrMissingSnapshotKey
	}

	if err := unmarshalDocument(snapshotDocument, value); err != nil {
		return errors.Wrap(err, ErrSnapshotFailedJSONUnmarshalling.Error())
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
		return errors.Wrap(err, ErrEventFailedJSONMarshalling.Error())
	}

	if err = json.Unmarshal(obj, &value); err != nil {
		return errors.Wrap(err, ErrEventFailedJSONUnmarshalling.Error())
	}

	return nil
}

func (store *EventStore) findOneSnapshot(
	filter interface{},
	options ...*options.FindOneOptions,
) (es.Snapshot, error) {
	var resultantSnapshot es.Snapshot

	action := func(ctx context.Context, collection *mongo.Collection) error {
		result := collection.FindOne(ctx, filter, options...)
		if result.Err() != nil {
			return errors.Wrap(result.Err(), ErrCouldNotFindSnapshots.Error())
		}

		if err := decodeSnapshot(result, &resultantSnapshot); err != nil {
			return errors.Wrap(err, ErrSnapshotCouldNotBeDecoded.Error())
		}

		return nil
	}
	err := store.connect(action, snapshotsCollection)

	return resultantSnapshot, errors.Wrap(err, ErrCouldNotConnectToEventStore.Error())
}

func (store *EventStore) Send(producer es.ProducerID, subject es.SubjectID, data []es.Data) ([]es.Event, error) {
	events, err := es.CreateEventBatch(producer, subject, es.Version(1), data, store)
	if err != nil {
		return nil, errors.Wrap(err, ErrEventBatchCreationFailedOnSend.Error())
	}

	return events, store.sendEvents(events)
}

func (store *EventStore) sendEvents(events []es.Event) error {
	if len(events) == 0 {
		return nil
	}

	documents := marshallEventDocuments(events)

	return store.insertManyDocuments(documents, eventsCollection)
}

func (store *EventStore) sendSnapshot(snapshot es.Snapshot) error {
	document := marshallSnapshotDocument(snapshot)

	return store.insertDocument(document, snapshotsCollection)
}

func (store *EventStore) Load(producer es.ProducerID, subject es.SubjectID, data es.Data) error {
	event, err := es.CreateEvent(producer, subject, es.Version(1), data, store)
	if err != nil {
		return errors.Wrap(err, ErrEventCreationFailedOnLoad.Error())
	}

	store.stage.AddEvent(event)

	return nil
}

func (store *EventStore) Clear() {
	for _, subject := range store.stage.Subjects() {
		store.stage.Clear(subject)
	}
}

func (store *EventStore) isStageInSync(subject es.SubjectID) bool {
	// Check whether the first staged event is
	// the next remote event in the remote store
	if store.stage.IsEmpty(subject) {
		return true
	}

	// We ignore the "found" bool return value
	// because we have just made the check "isStageEmpty"
	firstStagedEvent, _ := store.stage.FirstEvent(subject)

	latestRemoteEvent, err := store.latestRemoteEvent(subject)
	if errors.Is(err, es.ErrNoEvents) {
		return true
	}

	return firstStagedEvent.Version == latestRemoteEvent.Version+1
}

func (store *EventStore) shipSubject(subject es.SubjectID) error {
	if !store.isStageInSync(subject) {
		return ErrStageOutOfSync
	}

	stages := store.stage.EventStages(subject)
	for _, stage := range stages {
		if err := store.sendEvents(stage.Events()); err != nil {
			return errors.Wrap(err, "shipping the events failed")
		}

		if stage.Snapshot() != nil {
			if err := store.sendSnapshot(*stage.Snapshot()); err != nil {
				return errors.Wrap(err, "shipping the snapshot failed")
			}
		}
	}

	store.stage.Clear(subject)

	if err := recover(); err != nil {
		log.Println("Mongo Store paniced", err, ".")

		return nil
	}

	return nil
}

func (store *EventStore) Ship(ctx context.Context) error {
	// Clearing the insertion history is not "defered"
	//   the reason for this is: When an error cocured when
	//   shipping then it might be rollbacking itself failed.
	//   If we have cleared even though it failed then it would
	//   be impossible to ever return to the desired state
	//   where all insertions have been deleted because we would
	//   not know which documents to delete to acquire this.
	// As of now this is a defered call because I (Andreas) believe
	//   it causes an panic the other way
	defer store.clearInsertionHistory()

	subjects := store.stage.Subjects()
	for _, subject := range subjects {
		err := store.shipSubject(subject)
		if err != nil {
			log.Println("Shipping subject", subject, "failed")
			log.Println("Rollback issued because", err)
			// If an error is encountered of any count then rollback.
			//   we do so even if we lsot connection. Because in the
			//   mean time it is possible connection has been established
			// FIXED: Spike tests makes this rollback cause a panic
			//   This occurred because i called "Error()" on "rollbackErr"
			//   even when "rollbackErr" is nil causing a null dereference error
			rollbackErr := store.rollbackInsertions()
			if rollbackErr != nil {
				return errors.Wrap(err, rollbackErr.Error())
			}

			return errors.Wrap(err, "rollback successful")
		}
	}

	return nil
}

func (store *EventStore) Snapshot(producer es.ProducerID, subject es.SubjectID, data es.Data) error {
	snapshot, err := es.CreateSnapshot(producer, subject, es.Version(1), data, store)
	if err != nil {
		return errors.Wrap(err, "Snapshot creation failed")
	}

	store.stage.AddSnapshot(snapshot)

	return nil
}

func (store *EventStore) Concerning(subject es.SubjectID) ([]es.Event, error) {
	filter := bson.D{{Key: eventSubjectKey, Value: subject}}
	options := options.Find()
	options.SetSort(bson.D{{Key: snapshotVersionKey, Value: mongoAscending}})

	events, err := store.findAllEvents(filter, options)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (store *EventStore) By(producer es.ProducerID) ([]es.Event, error) {
	filter := bson.D{{Key: eventProducerKey, Value: producer}}
	options := options.Find()
	options.SetSort(bson.D{{Key: eventVersionKey, Value: mongoAscending}})

	events, err := store.findAllEvents(filter, options)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (store *EventStore) Between(subject es.SubjectID, from es.Version, to es.Version) ([]es.Event, error) {
	filter := bson.D{
		{Key: eventSubjectKey, Value: subject},
		{Key: eventVersionKey, Value: bson.D{
			{Key: mongoLessThan, Value: to},
			{Key: mongoGreaterThan, Value: from},
		}},
	}
	options := options.Find()
	options.SetSort(bson.D{{Key: eventVersionKey, Value: mongoAscending}})

	events, err := store.findAllEvents(filter, options)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (store *EventStore) With(subject es.SubjectID, snapshot es.Version) ([]es.Event, error) {
	filter := bson.D{
		{Key: eventSubjectKey, Value: subject},
		{Key: eventSnapShotVersionKey, Value: snapshot},
	}
	options := options.Find()
	options.SetSort(bson.D{{Key: eventVersionKey, Value: mongoAscending}})

	events, err := store.findAllEvents(filter, options)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (store *EventStore) After(subject es.SubjectID, pointInTime es.Timestamp) ([]es.Event, error) {
	return store.Temporal(subject, pointInTime, es.EndOfTime)
}

func (store *EventStore) Before(subject es.SubjectID, pointInTime es.Timestamp) ([]es.Event, error) {
	return store.Temporal(subject, pointInTime, es.BeginningOfTime)
}

func (store *EventStore) Temporal(subject es.SubjectID, from es.Timestamp, to es.Timestamp) ([]es.Event, error) {
	filter := bson.D{
		{Key: eventSubjectKey, Value: subject},
		{Key: eventTimestampKey, Value: bson.D{
			{Key: mongoLessThan, Value: to},
			{Key: mongoGreaterThan, Value: from},
		}},
	}
	options := options.Find()
	options.SetSort(bson.D{{Key: eventVersionKey, Value: mongoAscending}})

	events, err := store.findAllEvents(filter, options)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (store *EventStore) latestRemoteEvent(subject es.SubjectID) (es.Event, error) {
	filter := bson.D{{Key: eventSubjectKey, Value: subject}}
	options := options.FindOne()
	options.SetSort(bson.D{{Key: eventVersionKey, Value: mongoDescending}})

	event, err := store.findOneEvent(filter, options)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return es.Event{}, es.ErrNoEvents
	}

	return event, errors.Wrap(err, "latest remote event encountered an error")
}

func (store *EventStore) LatestEvent(subject es.SubjectID) (es.Event, error) {
	if latestStagedEvent, found := store.stage.LatestEvent(subject); found {
		return latestStagedEvent, nil
	}

	return store.latestRemoteEvent(subject)
}

func (store *EventStore) latestRemoteSnapshot(subject es.SubjectID) (es.Snapshot, error) {
	filter := bson.D{
		{Key: snapshotSubjectKey, Value: subject},
	}
	options := options.FindOne()
	options.SetSort(bson.D{{Key: snapshotVersionKey, Value: mongoDescending}})

	snapshot, err := store.findOneSnapshot(filter, options)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return es.Snapshot{}, es.ErrNoSnapshots
	}

	return snapshot, errors.Wrap(err, "latest remote snapshot encountered an error")
}

func (store *EventStore) LatestSnapshot(subject es.SubjectID) (es.Snapshot, error) {
	if latestStagedSnapshot, found := store.stage.LatestSnapshot(subject); found {
		return latestStagedSnapshot, nil
	}

	return store.latestRemoteSnapshot(subject)
}
