package races

import (
	"context"
	"sync"
	"time"

	"github.com/TBPixel/tww-rando-twitch-bot/internal/config"
	"github.com/TBPixel/tww-rando-twitch-bot/internal/racetime"
)

// Monitor actively pings the races API, maintaining a list
// of current races
type Monitor struct {
	config config.Racetime
	races  []racetime.RaceData
	mut    sync.Mutex
}

// NewMonitor creates a new racetime monitor
func NewMonitor(config config.Racetime) *Monitor {
	return &Monitor{
		config: config,
		races:  []racetime.RaceData{},
		mut:    sync.Mutex{},
	}
}

// Listen for new races, updating the local race list of races
func (m *Monitor) Listen(ctx context.Context) error {
	for {
		select {
		case <-time.After(time.Minute):
			res, err := racetime.CategoryDetail(m.config, m.config.Category)
			if err != nil {
				return err
			}

			m.mut.Lock()
			m.races = res.CurrentRaces
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
