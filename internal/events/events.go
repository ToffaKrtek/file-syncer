package events

import "github.com/ToffaKrtek/file-syncer/internal/config"

type EventsData struct {
	Events     map[string]Event    `json:"events"`
	Listeners  map[string]Listener `json:"listeners"`
	Subscribes map[string][]string `json:"subscribes"`
}

func (ed EventsData) FindEvent(eventName string) (Event, bool) {
	if event, got := ed.Events[eventName]; got {
		return event, got
	}
	return Event{}, false
}

// TODO::
// func (ed EventsData) AddListener(syncItemName string, eventNames []string) {}

func (ed EventsData) GetListeners(itemName string) ([]string, bool) {
	if listeners, got := ed.Subscribes[itemName]; got {
		return listeners, got
	}
	return []string{}, false
}

func (ed EventsData) FindListener(listenerName string) (Listener, bool) {
	if listener, got := ed.Listeners[listenerName]; got {
		return listener, got
	}
	return Listener{}, false
}

func (ed EventsData) Trigger(item config.Item) {
	if listeners, got := ed.GetListeners(item.Name); got {
		for _, listenerName := range listeners {
			if listener, ok := ed.FindListener(listenerName); ok {
				listener.Trigger()
			}
		}
	}
}

type Listener struct {
	SyncItemName string   `json:"`
	EventNames   []string `json:"event_names"`
}

func (l Listener) Trigger() {
	for _, eventName := range l.EventNames {
		if event, ok := ed.FindEvent(eventName); ok {
			event.Run()
		}
	}
}

type Event struct {
	Job  Job    `json:"job"`
	Name string `json:"name"`
}

type Job interface {
	Run(config.Item) error
}
