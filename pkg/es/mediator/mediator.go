package mediator

import "github.com/hywmongous/example-service/pkg/es"

type Subscription func(subject es.SubjectID, data es.Data)
type Connector chan ConnectorResult

type ConnectorResult struct {
	Subject es.SubjectID
	Data    es.Data
}

type Mediator interface {
	ListenTo(topic es.Topic, receiver Subscription)
	Listen(receiver Subscription)
	ChannelTo(topic es.Topic)
	Channel()
	Publish(topic es.Topic, data es.Data)
}

type defaultMediator struct {
	connectors          map[es.Topic][]Connector
	universalConnectors []Connector
	receivers           map[es.Topic][]Subscription
	universalReceivers  []Subscription
	Mediator
}

// Singleton for the mediator used when one is not defined
// This makes it possible to write "mediator.Publish(...)" and so one
var Default defaultMediator = createDefaultMediator()

func createDefaultMediator() defaultMediator {
	return defaultMediator{
		connectors:          map[es.Topic][]Connector{},
		universalConnectors: make([]Connector, 0),
		receivers:           map[es.Topic][]Subscription{},
		universalReceivers:  make([]Subscription, 0),
	}
}

func createConnector() Connector {
	return make(Connector)
}

func ListenTo(topic es.Topic, receiver Subscription) {
	if _, ok := Default.receivers[topic]; !ok {
		Default.receivers[topic] = make([]Subscription, 0)
	}
	Default.receivers[topic] = append(Default.receivers[topic], receiver)
}

func Listen(receiver Subscription) {
	Default.universalReceivers = append(Default.universalReceivers, receiver)
}

func ChannelTo(topic es.Topic) Connector {
	channel := createConnector()
	if _, ok := Default.connectors[topic]; !ok {
		Default.connectors[topic] = make([]Connector, 0)
	}
	Default.connectors[topic] = append(Default.connectors[topic], channel)
	return channel
}

func Channel() Connector {
	channel := createConnector()
	Default.universalConnectors = append(Default.universalConnectors, channel)
	return channel
}

func Publish(subject es.SubjectID, data es.Data) {
	connectorResult := ConnectorResult{
		Subject: subject,
		Data:    data,
	}

	topic := es.CreateTopicForData(data)
	for _, receiver := range Default.receivers[topic] {
		receiver(subject, data)
	}
	for _, connector := range Default.connectors[topic] {
		connector <- connectorResult
	}
	for _, receiver := range Default.universalReceivers {
		receiver(subject, data)
	}
	for _, connector := range Default.universalConnectors {
		connector <- connectorResult
	}
}
