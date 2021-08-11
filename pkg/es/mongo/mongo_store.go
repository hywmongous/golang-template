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

func CreateMongoEventStore() (MongoEventStore, error) {
	return MongoEventStore{}, nil
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
	// Create the client
	client, err := store.createClient()
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

func (store MongoEventStore) createClient() (*mongo.Client, error) {
	options := options.Client()
	uri := options.ApplyURI("mongodb://root:root@ia_mongo:27017")

	client, err := mongo.NewClient(uri)
	if err != nil {
		return nil, merr.CreateFailedStructInvocation("MongoEventStore", "createClient", err)
	}

	return client, nil
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

func convertEventsToDocuments(events []es.Event) []interface{} {
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
	eventMap := eventDocument
	obj, err := json.Marshal(eventMap)
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

func (store MongoEventStore) Stock(producer es.ProducerID, subject es.SubjectID, data []es.EventData) error {
	events, err := es.CreateEventBatch(producer, subject, data, store)
	if err != nil {
		return merr.CreateFailedStructInvocation("EventStore", "Send", err)
	}

	documents := convertEventsToDocuments(events)
	return store.insertManyDocuments(documents)
}

func (store MongoEventStore) Retrieve(subject es.SubjectID, callback es.Callback) ([]es.Event, error) {
	filter := bson.D{{Key: "event.subject", Value: subject}}
	options := options.Find()
	options.SetSort(bson.D{{Key: "event.version", Value: 1}})

	loop := func(cursor *mongo.Cursor) error {
		var document bson.M
		if err := cursor.Decode(&document); err != nil {
			return merr.CreateFailedStructInvocation("MongoEventStore", "Retrieve.callback", err)
		}

		event, err := unmarshalDocument(document)
		if err != nil {
			return merr.CreateFailedStructInvocation("MongoEventStore", "Retrieve.callback", err)
		}

		if err = callback(event); err != nil {
			return err
		}

		return nil
	}

	err := store.findAll(filter, loop, options)
	if err != nil {
		return nil, merr.CreateFailedStructInvocation("MongoEventStore", "Retrieve", err)
	}

	return nil, merr.CreateNotImplementedYetStruct("MongoEventStore", "Retrieve")
}

func (store MongoEventStore) LatestEvent(subject es.SubjectID, dataType es.EventDataType) (es.Event, error) {
	filter := bson.D{{Key: "event.subject", Value: subject}}
	options := options.FindOne()
	options.SetSort(bson.D{{Key: "event.version", Value: -1}})

	result, err := store.findOne(filter, options)
	if err != nil {
		return es.Event{}, merr.CreateFailedStructInvocation("MongoEventStore", "LatestEvent", err)
	}

	var document bson.M
	if err = result.Decode(&document); err != nil {
		return es.Event{}, merr.CreateFailedStructInvocation("MongoEventStore", "LatestEvent", err)
	}

	event, err := unmarshalDocument(document)
	if err != nil {
		return es.Event{}, merr.CreateFailedStructInvocation("MongoEventStore", "LatestEvent", err)
	}

	return event, nil
}
