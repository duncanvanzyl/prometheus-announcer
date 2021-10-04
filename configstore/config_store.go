package configstore

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	sd "github.com/duncanvanzyl/prometheus-announcer/servicediscovery"
)

type config struct {
	targets []string
	labels  map[string]string
	t       time.Time
}

type ConfigStore struct {
	lifetime time.Duration
	interval time.Duration
	mu       sync.RWMutex
	configs  map[string]*config
}

func New(lifetime, interval time.Duration) *ConfigStore {
	return &ConfigStore{
		lifetime: lifetime,
		interval: interval,
		configs:  map[string]*config{},
	}
}

func (cs *ConfigStore) AddTarget(id string, host []string, ls map[string]string) error {
	logger.Debug("request add target", "host", host, "labels", ls)

	cs.mu.Lock()
	defer cs.mu.Unlock()

	cs.configs[id] = &config{
		targets: host,
		labels:  ls,
		t:       time.Now(),
	}

	return nil
}

func (cs *ConfigStore) RemoveTarget(id string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	delete(cs.configs, id)
}

func (cs *ConfigStore) JSON() ([]byte, error) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	j := []sd.Config{}
	for _, conf := range cs.configs {
		c := sd.Config{
			Targets: conf.targets,
			Labels:  conf.labels,
		}
		j = append(j, c)
	}

	return json.Marshal(j)
}

func (cs *ConfigStore) purge(t time.Time) {
	logger.Debug("checking purges", "time", t, "configs", cs.configs)
	cs.mu.Lock()
	defer cs.mu.Unlock()

	for id, conf := range cs.configs {
		if t.Sub(conf.t) > cs.lifetime {
			logger.Info("purging target", "id", id, "target", conf)
			delete(cs.configs, id)
		}
	}
}

func (cs *ConfigStore) Run(ctx context.Context) {
	tick := time.NewTicker(cs.interval)
	defer tick.Stop()

	for {
		select {
		case t := <-tick.C:
			cs.purge(t)
		case <-ctx.Done():
			logger.Info("config store stop on context", "error", ctx.Err())
			return
		}
	}
}
