package mongo

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	merr "github.com/hywmongous/example-service/pkg/errors"
	"github.com/hywmongous/example-service/pkg/es"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	mongoAction func(context context.Context, collection *mongo.Collection) error
)

type MongoEventStore struct {
	commit []es.Event
}

const (
	initialCommitLength = 0 // We append each time an event is added
	timeoutDuration     = 10 * time.Second
	databaseName        = "eventstore"
	collectionName      = "events"
)

const (
	EventIdKey              = "event.id"
	EventProducerKey        = "event.producer"
	EventSubjectKey         = "event.subject"
	EventVersionKey         = "event.version"
	EventSchemaVersionKey   = "event.schemaversion"
	EventSnapShotVersionKey = "event.snapshotversion"
	EventNameKey            = "event.name"
	EventTimestampKey       = "event.timestamp"
	EventDataKey            = "event.data"

	MongoLessThan    = "$lt"
	MongoGreaterThan = "$gt"
)

var (
	ErrDatabaseNotFound   = errors.New("eventstore database could not be found")
	ErrCollectionNotFound = errors.New("events collection could not be found")
	ErrInsertion          = errors.New("failed inserting one or more documents into collection")
	ErrEventNotFound      = errors.New("event not found")
)

func CreateMongoEventStore() MongoEventStore {
	return MongoEventStore{
		commit: make([]es.Event, initialCommitLength),
	}
}

func (store MongoEventStore) collection(client *mongo.Client) (*mongo.Collection, error) {
	// Establish database connection
	database := client.Database(databaseName)
	if database == nil {
		return nil, merr.CreateFailedStructInvocation("MongoEventStore", "collection", ErrDatabaseNotFound)
	}

	// Establish collection connection
	collection := database.Collection(collectionName)
	if collection == nil {
		return nil, merr.CreateFailedStructInvocation("MongoEventStore", "collection", ErrCollectionNotFound)
	}

	return collection, nil
}

func (store MongoEventStore) connect(action mongoAction) error {
	options := options.Client()
	uri := options.ApplyURI("mongodb://root:root@ia_mongo:27017")

	// Client construction
	client, err := mongo.NewClient(uri)
	if err != nil {
		return merr.CreateFailedStructInvocation("MongoEventStore", "connect", err)
	}

	// Create the client
	if err != nil {
		return merr.CreateFailedStructInvocation("MongoEventStore", "connect", err)
	}

	// Create the context
	ctx, cancel := context.WithTimeout(context.Background(), timeoutDuration)
	defer cancel()

	// Construct the connected client
	err = client.Connect(ctx)
	if err != nil {
		return merr.CreateFailedStructInvocation("MongoEventStore", "connect", err)
	}
	defer client.Disconnect(ctx)

	// Connect to the collection
	collection, err := store.collection(client)
	if err != nil {
		return merr.CreateFailedStructInvocation("MongoEventStore", "connect", err)
	}

	return action(ctx, collection)
}

func (store MongoEventStore) insertManyDocuments(documents []interface{}) error {
	action := func(ctx context.Context, collection *mongo.Collection) error {
		_, err := collection.InsertMany(ctx, documents)
		if err != nil {
			return merr.CreateFailedStructInvocation("MongoEventStore", "action", err)
		}
		return nil
	}

	return store.connect(action)
}

func (store MongoEventStore) findOne(filter interface{}, options ...*options.FindOneOptions) (es.Event, error) {
	var event es.Event
	var result *mongo.SingleResult
	action := func(ctx context.Context, collection *mongo.Collection) error {
		result = collection.FindOne(ctx, filter, options...)
		if result == nil {
			return merr.CreateFailedStructInvocation("MongoEventStore", "findOne.action", result.Err())
		}

		var document bson.M
		if err := result.Decode(&document); err != nil {
			return merr.CreateFailedStructInvocation("MongoEventStore", "findOne.action", err)
		}

		unmarshaledEvent, err := unmarshalDocument(document)
		if err != nil {
			return merr.CreateFailedStructInvocation("MongoEventStore", "findOne.action", err)
		}
		event = unmarshaledEvent

		return nil
	}
	return event, store.connect(action)
}

func (store MongoEventStore) findAll(filter interface{}, options ...*options.FindOptions) ([]es.Event, error) {
	var events []es.Event
	var cursor *mongo.Cursor
	action := func(ctx context.Context, collection *mongo.Collection) error {
		result, err := collection.Find(ctx, filter, options...)
		if err != nil {
			return merr.CreateFailedStructInvocation("MongoEventStore", "findAll.action", err)
		}
		cursor = result
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var document bson.M
			if err := cursor.Decode(&document); err != nil {
				return merr.CreateFailedStructInvocation("MongoEventStore", "For.loop", err)
			}

			event, err := unmarshalDocument(document)
			if err != nil {
				return merr.CreateFailedStructInvocation("MongoEventStore", "For.loop", err)
			}

			events = append(events, event)
		}

		return nil
	}

	if err := store.connect(action); err != nil {
		return nil, err
	}
	return events, nil
}

func marshallDocuments(events []es.Event) []interface{} {
	documents := make([]interface{}, len(events))
	for idx, event := range events {
		documents[idx] = bson.D{
			bson.E{
				Key:   "event",
				Value: event,
			},
		}
	}
	return documents
}

func unmarshalDocument(document bson.M) (es.Event, error) {
	eventDocument, ok := document["event"].(bson.M)
	if !ok {
		return es.Event{}, merr.ErrTechnical
	}

	// The following is a JSON work around golang mongodb
	// driver does not support decoding of interface{}.
	// This caused issues with the evenData within
	// the event itself. This is however supported by
	// the json package, so for Marshalling we use json
	obj, err := json.Marshal(eventDocument)
	if err != nil {
		return es.Event{}, merr.CreateFailedInvocation("unmarshalDocument", err)
	}

	var event es.Event
	err = json.Unmarshal(obj, &event)
	if err != nil {
		return es.Event{}, merr.CreateFailedInvocation("unmarshalDocument", err)
	}

	return event, nil
}

func (store MongoEventStore) Send(producer es.ProducerID, subject es.SubjectID, data []es.EventData) ([]es.Event, error) {
	events, err := es.CreateEventBatch(producer, subject, es.EventSchemaVersion(1), data, store)
	if err != nil {
		return nil, merr.CreateFailedStructInvocation("EventStore", "Send", err)
	}
	return events, store.SendEvents(events)
}

func (store MongoEventStore) SendEvents(events []es.Event) error {
	documents := marshallDocuments(events)
	return store.insertManyDocuments(documents)
}

func (store MongoEventStore) Load(producer es.ProducerID, subject es.SubjectID, data es.EventData) (es.Event, error) {
	event, err := es.CreateEvent(producer, subject, es.EventSchemaVersion(1), data, store)
	if err != nil {
		return es.Event{}, merr.CreateFailedStructInvocation("EventStore", "Load", err)
	}
	(&store).commit = append(store.commit, event)
	return event, nil
}

func (store MongoEventStore) Unload(eventId es.EventId) (es.Event, error) {
	for idx, event := range store.commit {
		if event.Id == eventId {
			(&store).commit = append(store.commit[:idx], store.commit[idx+1:]...)
			return event, nil
		}
	}
	return es.Event{}, ErrEventNotFound
}

func (store MongoEventStore) Ship() ([]es.Event, error) {
	return store.commit, store.SendEvents(store.commit)
}

func (store MongoEventStore) Concerning(subject es.SubjectID) ([]es.Event, error) {
	filter := bson.D{{Key: EventSubjectKey, Value: subject}}
	options := options.Find()
	options.SetSort(bson.D{{Key: EventVersionKey, Value: -1}})

	events, err := store.findAll(filter, options)
	if err != nil {
		return nil, merr.CreateFailedStructInvocation("MongoEventStore", "Concerning", err)
	}

	return events, nil
}

func (store MongoEventStore) By(producer es.ProducerID) ([]es.Event, error) {
	filter := bson.D{{Key: EventProducerKey, Value: producer}}
	options := options.Find()
	options.SetSort(bson.D{{Key: EventVersionKey, Value: -1}})

	events, err := store.findAll(filter, options)
	if err != nil {
		return nil, merr.CreateFailedStructInvocation("MongoEventStore", "By", err)
	}

	return events, nil
}

func (store MongoEventStore) Between(subject es.SubjectID, from es.EventVersion, to es.EventVersion) ([]es.Event, error) {
	filter := bson.D{
		{Key: EventSubjectKey, Value: subject},
		{Key: EventVersionKey, Value: bson.D{
			{Key: MongoLessThan, Value: to},
			{Key: MongoGreaterThan, Value: from},
		}},
	}
	options := options.Find()
	options.SetSort(bson.D{{Key: EventVersionKey, Value: -1}})

	events, err := store.findAll(filter, options)
	if err != nil {
		return nil, merr.CreateFailedStructInvocation("MongoEventStore", "By", err)
	}

	return events, nil
}

func (store MongoEventStore) With(subject es.SubjectID, snapshot es.SnapshotVersion) ([]es.Event, error) {
	filter := bson.D{
		{Key: EventSubjectKey, Value: subject},
		{Key: EventSnapShotVersionKey, Value: snapshot},
	}
	options := options.Find()
	options.SetSort(bson.D{{Key: EventVersionKey, Value: -1}})

	events, err := store.findAll(filter, options)
	if err != nil {
		return nil, merr.CreateFailedStructInvocation("MongoEventStore", "In", err)
	}

	return events, nil
}

func (store MongoEventStore) After(subject es.SubjectID, pointInTime es.EventTimestamp) ([]es.Event, error) {
	return store.Temporal(subject, pointInTime, es.EndOfTime)
}

func (store MongoEventStore) Before(subject es.SubjectID, pointInTime es.EventTimestamp) ([]es.Event, error) {
	return store.Temporal(subject, pointInTime, es.BeginningOfTime)
}

func (store MongoEventStore) Temporal(subject es.SubjectID, from es.EventTimestamp, to es.EventTimestamp) ([]es.Event, error) {
	filter := bson.D{
		{Key: EventSubjectKey, Value: subject},
		{Key: EventTimestampKey, Value: bson.D{
			{Key: MongoLessThan, Value: to},
			{Key: MongoGreaterThan, Value: from},
		}},
	}
	options := options.Find()
	options.SetSort(bson.D{{Key: EventVersionKey, Value: -1}})

	events, err := store.findAll(filter, options)
	if err != nil {
		return nil, merr.CreateFailedStructInvocation("MongoEventStore", "Temporal", err)
	}

	return events, nil
}

func (store MongoEventStore) Latest(subject es.SubjectID) (es.Event, error) {
	filter := bson.D{{Key: EventSubjectKey, Value: subject}}
	options := options.FindOne()
	options.SetSort(bson.D{{Key: EventVersionKey, Value: -1}})

	event, err := store.findOne(filter, options)
	if err != nil {
		return es.Event{}, merr.CreateFailedStructInvocation("MongoEventStore", "Latest", err)
	}
	return event, nil
}
