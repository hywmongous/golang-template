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
	mongoAction   func(context context.Context, collection *mongo.Collection) error
	mongoCallback func(cursor *mongo.Cursor) error
)

type MongoEventStore struct{}

const (
	timeoutDuration = 10 * time.Second
	databaseName    = "eventstore"
	collectionName  = "events"
)

var (
	ErrDatabaseNotFound   = errors.New("eventstore database could not be found")
	ErrCollectionNotFound = errors.New("events collection could not be found")
	ErrInsertion          = errors.New("failed inserting one or more documents into collection")
)

func CreateMongoEventStore() MongoEventStore {
	return MongoEventStore{}
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

func (store MongoEventStore) findOne(filter interface{}, options ...*options.FindOneOptions) (*mongo.SingleResult, error) {
	var result *mongo.SingleResult
	action := func(ctx context.Context, collection *mongo.Collection) error {
		result = collection.FindOne(ctx, filter, options...)
		if result == nil {
			return merr.CreateFailedStructInvocation("MongoEventStore", "findOne.action", result.Err())
		}
		return nil
	}
	return result, store.connect(action)
}

func (store MongoEventStore) findAll(filter interface{}, factory mongoCallback, options ...*options.FindOptions) error {
	var cursor *mongo.Cursor
	action := func(ctx context.Context, collection *mongo.Collection) error {
		result, err := collection.Find(ctx, filter, options...)
		if err != nil {
			return merr.CreateFailedStructInvocation("MongoEventStore", "findAll.action", err)
		}
		cursor = result
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			if err = factory(cursor); err != nil {
				return merr.CreateFailedStructInvocation("MongoEventStore", "findAll.action", err)
			}
		}

		return nil
	}
	return store.connect(action)
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

func (store MongoEventStore) Stock(producer es.ProducerID, subject es.SubjectID, data []es.EventData) ([]es.Event, error) {
	events, err := es.CreateEventBatch(producer, subject, es.EventSchemaVersion(1), data, store)
	if err != nil {
		return nil, merr.CreateFailedStructInvocation("EventStore", "Send", err)
	}

	documents := marshallDocuments(events)
	return events, store.insertManyDocuments(documents)
}

func (store MongoEventStore) Retrieve(subject es.SubjectID) ([]es.Event, error) {
	filter := bson.D{{Key: "event.subject", Value: subject}}
	options := options.Find()
	options.SetSort(bson.D{{Key: "event.version", Value: 1}})
	var events []es.Event

	loop := func(cursor *mongo.Cursor) error {
		var document bson.M
		if err := cursor.Decode(&document); err != nil {
			return merr.CreateFailedStructInvocation("MongoEventStore", "Retrieve.loop", err)
		}

		event, err := unmarshalDocument(document)
		if err != nil {
			return merr.CreateFailedStructInvocation("MongoEventStore", "Retrieve.loop", err)
		}

		events = append(events, event)

		return nil
	}

	err := store.findAll(filter, loop, options)
	if err != nil {
		return nil, merr.CreateFailedStructInvocation("MongoEventStore", "Retrieve", err)
	}

	return events, nil
}

func (store MongoEventStore) CurrentEventVersion(subject es.SubjectID) es.EventVersion {
	// TODO: There must be a query where we check whether
	//   the first event has been published if so then
	//   we retrieve the latest and return it's version
	latest, err := store.Latest(subject)
	if err != nil {
		return es.InitialEventVersion
	}
	return latest.Version
}

func (store MongoEventStore) Latest(subject es.SubjectID) (es.Event, error) {
	filter := bson.D{{Key: "event.subject", Value: subject}}
	options := options.FindOne()
	options.SetSort(bson.D{{Key: "event.version", Value: -1}})

	result, err := store.findOne(filter, options)
	if err != nil {
		return es.Event{}, merr.CreateFailedStructInvocation("MongoEventStore", "Latest", err)
	}

	var document bson.M
	if err = result.Decode(&document); err != nil {
		return es.Event{}, merr.CreateFailedStructInvocation("MongoEventStore", "Latest", err)
	}

	event, err := unmarshalDocument(document)
	if err != nil {
		return es.Event{}, merr.CreateFailedStructInvocation("MongoEventStore", "Latest", err)
	}

	return event, nil
}
