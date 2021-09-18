package es

type EventStage struct {
	events   []Event
	snapshot *Snapshot
}

type Stage struct {
	subjects map[SubjectID][]EventStage
}

func CreateStage() Stage {
	return Stage{
		subjects: map[SubjectID][]EventStage{},
	}
}

func CreateEventStage() EventStage {
	return EventStage{
		events:   make([]Event, 0),
		snapshot: nil,
	}
}

func (eventStage EventStage) Events() []Event {
	return eventStage.events
}

func (eventStage EventStage) Snapshot() *Snapshot {
	return eventStage.snapshot
}

func (stage Stage) Events() []Event {
	events := make([]Event, 0, 32)
	for _, subject := range stage.Subjects() {
		for _, eventStage := range stage.EventStages(subject) {
			events = append(events, eventStage.events...)
		}
	}
	return events
}

func (stage *Stage) Subjects() []SubjectID {
	subjects := make([]SubjectID, len(stage.subjects))
	idx := 0
	for subject := range stage.subjects {
		subjects[idx] = subject
		idx++
	}
	return subjects
}

func (stage *Stage) Clear(subject SubjectID) {
	if _, found := stage.subjects[subject]; !found {
		return
	}
	stage.subjects[subject] = make([]EventStage, 1)
}

func (stage *Stage) IsEmpty(subject SubjectID) bool {
	if _, found := stage.subjects[subject]; !found {
		return true
	}
	return len(stage.subjects[subject][0].events) == 0
}

func (stage *Stage) addEventStage(subject SubjectID) {
	stage.subjects[subject] = append(
		stage.subjects[subject],
		CreateEventStage(),
	)
}

func (stage *Stage) EventStages(subject SubjectID) []EventStage {
	_, found := stage.subjects[subject]
	if !found {
		stage.subjects[subject] = make([]EventStage, 0)
		stage.addEventStage(subject)
	}
	return stage.subjects[subject]
}

func (stage *Stage) HasSubject(subject SubjectID) bool {
	_, found := stage.subjects[subject]
	return found
}

func (stage *Stage) FirstEvent(subject SubjectID) (Event, bool) {
	if !stage.HasSubject(subject) {
		return Event{}, false
	}

	firstStage := stage.subjects[subject][0]
	if len(firstStage.events) > 0 {
		return firstStage.events[0], true
	}
	return Event{}, false
}

func (stage *Stage) LatestEvent(subject SubjectID) (Event, bool) {
	if !stage.HasSubject(subject) {
		return Event{}, false
	}

	eventStages := stage.subjects[subject]

	found := false
	var latestEvent Event
	for _, eventStage := range eventStages {
		if len(eventStage.events) > 0 {
			latestEvent = eventStage.events[len(eventStage.events)-1]
			found = true
		}
	}
	return latestEvent, found
}

func (stage *Stage) LatestSnapshot(subject SubjectID) (Snapshot, bool) {
	eventStages := stage.EventStages(subject)
	if len(eventStages) > 1 {
		return *eventStages[len(eventStages)-2].snapshot, true
	}
	return Snapshot{}, false
}

func (stage *Stage) AddEvent(event Event) {
	eventStages := stage.EventStages(event.Subject)
	eventStages[len(eventStages)-1].events = append(eventStages[len(eventStages)-1].events, event)
}

func (stage *Stage) AddSnapshot(snapshot Snapshot) {
	eventStages := stage.EventStages(snapshot.Subject)
	eventStages[len(eventStages)-1].snapshot = &snapshot
	stage.subjects[snapshot.Subject] = append(eventStages, CreateEventStage())
}
