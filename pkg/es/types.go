package es

import (
	"math"
	"reflect"
	"strings"
)

const (
	InitialEventVersion       = Version(0)
	InitialEventSchemaVersion = Version(0)
	InitialSnapshotVersion    = Version(0)

	BeginningOfTime = Timestamp(0)
	EndOfTime       = Timestamp(math.MaxInt64)
)

type (
	Title      string
	Ident      string
	ProducerID string
	SubjectID  string
	Version    uint
	Timestamp  int64
	Data       interface{}
)

func CreateTitleForData(data Data) Title {
	eventType := reflect.TypeOf(data).String()
	eventTypeParts := strings.Split(eventType, ".")
	eventName := eventTypeParts[len(eventTypeParts)-1]

	return Title(eventName)
}

func CreateTopicForData(data Data) Topic {
	return Topic(CreateTitleForData(data))
}
