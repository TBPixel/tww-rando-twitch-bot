package races

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/TBPixel/tww-rando-twitch-bot/internal/config"
	"github.com/TBPixel/tww-rando-twitch-bot/internal/racetime"
)

// Monitor actively pings the races API, maintaining a list
// of current races
type Monitor struct {
	category      string
	config        config.Racetime
	races         []racetime.RaceData
	mut           sync.Mutex
	listeners     []chan []racetime.RaceData
	listenerMutex sync.Mutex
}

// NewMonitor creates a new racetime monitor
func NewMonitor(config config.Racetime, category string) *Monitor {
	return &Monitor{
		category:      category,
		config:        config,
		races:         []racetime.RaceData{},
		mut:           sync.Mutex{},
		listeners:     []chan []racetime.RaceData{},
		listenerMutex: sync.Mutex{},
	}
}

func (m *Monitor) AddListener() chan []racetime.RaceData {
	m.listenerMutex.Lock()
	defer m.listenerMutex.Unlock()

	listener := make(chan []racetime.RaceData)
	m.listeners = append(m.listeners, listener)

	return listener
}

func (m *Monitor) RemoveListener(listener chan []racetime.RaceData) {
	m.listenerMutex.Lock()
	defer m.listenerMutex.Unlock()

	for i, l := range m.listeners {
		if l != listener {
			continue
		}

		m.listeners = append(m.listeners[:i], m.listeners[i+1:]...)
	}
}

// Listen for new races, updating the local race list of races
func (m *Monitor) Listen(ctx context.Context) error {
	races, err := racetime.CategoryRaces(m.config, m.category)
	if err != nil {
		log.Println(err)
	} else {
		m.mut.Lock()
		m.races = races
		m.emit(races)
		m.mut.Unlock()
	}

	for {
		select {
		case <-time.After(m.config.RaceRefreshInterval):
			races, err := racetime.CategoryRaces(m.config, m.category)
			if err != nil {
				log.Println(err)
				return nil
			}

			m.mut.Lock()
			m.races = races
			m.emit(races)
			m.mut.Unlock()
		case <-ctx.Done():
			return nil
		}
	}
}

// Races returns a list of the current races
func (m *Monitor) Races() []racetime.RaceData {
	m.mut.Lock()
	defer m.mut.Unlock()

	return m.races
}

func (m *Monitor) emit(races []racetime.RaceData) {
	m.listenerMutex.Lock()
	defer m.listenerMutex.Unlock()

	for _, l := range m.listeners {
		l <- races
	}
}
