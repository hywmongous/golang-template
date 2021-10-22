package mediator

import (
	"github.com/hywmongous/example-service/pkg/es"
)

type (
	Subscription func(subject es.SubjectID, data es.Data)
	Connector    chan ConnectorResult
)

type ConnectorResult struct {
	Subject es.SubjectID
	Data    es.Data
}

type Mediator struct {
	connectors          map[es.Topic][]Connector
	universalConnectors []Connector
	receivers           map[es.Topic][]Subscription
	universalReceivers  []Subscription
}

func Create() *Mediator {
	return &Mediator{
		connectors:          make(map[es.Topic][]Connector),
		universalConnectors: make([]Connector, 0),
		receivers:           make(map[es.Topic][]Subscription),
		universalReceivers:  make([]Subscription, 0),
	}
}

func createConnector() Connector {
	return make(Connector)
}

func (mediator *Mediator) ListenTo(topic es.Topic, receiver Subscription) {
	if _, ok := mediator.receivers[topic]; !ok {
		mediator.receivers[topic] = make([]Subscription, 0)
	}

	mediator.receivers[topic] = append(mediator.receivers[topic], receiver)
}

func (mediator *Mediator) Listen(receiver Subscription) {
	mediator.universalReceivers = append(mediator.universalReceivers, receiver)
}

func (mediator *Mediator) ChannelTo(topic es.Topic) Connector {
	channel := createConnector()

	if _, ok := mediator.connectors[topic]; !ok {
		mediator.connectors[topic] = make([]Connector, 0)
	}

	mediator.connectors[topic] = append(mediator.connectors[topic], channel)

	return channel
}

func (mediator *Mediator) Channel() Connector {
	channel := createConnector()
	mediator.universalConnectors = append(mediator.universalConnectors, channel)

	return channel
}

func (mediator *Mediator) Publish(subject es.SubjectID, data es.Data) {
	connectorResult := ConnectorResult{
		Subject: subject,
		Data:    data,
	}

	topic := es.CreateTopicForData(data)

	for _, receiver := range mediator.receivers[topic] {
		receiver(subject, data)
	}

	for _, connector := range mediator.connectors[topic] {
		connector <- connectorResult
	}

	for _, receiver := range mediator.universalReceivers {
		receiver(subject, data)
	}

	for _, connector := range mediator.universalConnectors {
		connector <- connectorResult
	}
}
