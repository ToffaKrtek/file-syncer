package events

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/ToffaKrtek/file-syncer/internal/config"
)

type eventsData struct {
	Events     map[string]Event    `json:"events"`
	Listeners  map[string]Listener `json:"listeners"`
	Subscribes map[string][]string `json:"subscribes"`
}

var (
	eventsRepoFile = "./eventsRepo.json"
	eData          *eventsData
	locked         int
	mu             sync.Mutex
)

func eventRepository() *eventsData {
	mu.Lock()
	defer mu.Unlock()
	if eData == nil {
		var err error
		_, err = loadEventRepo()
		if err != nil {
			panic("Ошибка загрузки репозитория событий")
		}
	}
	locked++
	return eData
}

func closeRepo() {
	mu.Lock()
	defer mu.Unlock()
	locked--
	if locked == 0 {
		eData = nil
	}
}

func loadEventRepo() (*eventsData, error) {
	if _, err := os.Stat(eventsRepoFile); os.IsNotExist(err) {
		edata := &eventsData{
			Events:     make(map[string]Event),
			Listeners:  make(map[string]Listener),
			Subscribes: make(map[string][]string),
		}
		saveEventRepo(edata)
	}
	data, err := os.ReadFile(eventsRepoFile)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, eData); err != nil {
		return nil, err
	}
	return eData, nil
}

func saveEventRepo(edata *eventsData) error {
	data, err := json.MarshalIndent(edata, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(eventsRepoFile, data, 0644); err != nil {
		return err
	}
	return nil
}

func FindEvent(eventName string) (Event, bool) {
	defer closeRepo()
	if event, got := eventRepository().Events[eventName]; got {
		return event, got
	}
	return Event{}, false
}

// TODO::
// func (ed eventsData) AddListener(syncItemName string, eventNames []string) {}

func GetListeners(itemName string) ([]string, bool) {
	defer closeRepo()
	if listeners, got := eventRepository().Subscribes[itemName]; got {
		return listeners, got
	}
	return []string{}, false
}

func FindListener(listenerName string) (Listener, bool) {
	defer closeRepo()
	if listener, got := eventRepository().Listeners[listenerName]; got {
		return listener, got
	}
	return Listener{}, false
}

func Trigger(item config.Item) {
	defer closeRepo()
	if listeners, got := GetListeners(item.Name); got {
		for _, listenerName := range listeners {
			if listener, ok := FindListener(listenerName); ok {
				listener.Trigger()
			}
		}
	}
}

type Listener struct {
	SyncItemName string   `json:"item_name"`
	EventNames   []string `json:"event_names"`
}

func (l Listener) Trigger() {
	for _, eventName := range l.EventNames {
		if event, ok := FindEvent(eventName); ok {
			event.Run()
		}
	}
}

type Event struct {
	Job  Job    `json:"job"`
	Name string `json:"name"`
}

func (e Event) Run() {
	// TODO:: run job
}

type Job interface {
	Run(config.Item) error
}
