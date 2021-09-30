package configstore

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	pa "github.com/duncanvanzyl/prometheus-announcer"

	"github.com/hashicorp/go-hclog"
)

var logger = hclog.Default()

func SetLogger(l hclog.Logger) {
	logger = l.Named("cs")
}

type target struct {
	host string
	ct   pa.ConfigType
	t    time.Time
}

func (t target) String() string {
	return fmt.Sprintf("Host: %q, Type: %v, Time: %s", t.host, t.ct, t.t.Format("15:04:05"))
}

type staticConfig struct {
	Targets []string          `json:"targets"`
	Labels  map[string]string `json:"labels,omitempty"`
}

type Config interface {
	AddTarget(string, pa.ConfigType) error
	RemoveTarget(string)
	JSON() ([]byte, error)
	Run(context.Context)
}

type ConfigStore struct {
	mu       sync.RWMutex
	configs  []target
	lifetime time.Duration
	interval time.Duration
}

func New(lifetime, interval time.Duration) *ConfigStore {
	return &ConfigStore{
		lifetime: lifetime,
		interval: interval,
	}
}

func (cs *ConfigStore) hasTarget(host string, ct pa.ConfigType) *target {
	for i := range cs.configs {
		t := &cs.configs[i]
		if t.ct == ct && t.host == host {
			return t
		}
	}
	return nil
}

func (cs *ConfigStore) AddTarget(host string, ct pa.ConfigType) error {
	logger.Debug("request add target", "host", host, "type", ct)

	cs.mu.Lock()
	defer cs.mu.Unlock()

	if ct >= 2 { // 2 is a magic number since there are 2 configuration types defined
		return fmt.Errorf("invalid config type")
	}

	if t := cs.hasTarget(host, ct); t != nil {
		logger.Debug("target already exists", "host", host, "type", ct)
		t.t = time.Now()
		return nil
	}

	// TODO: check host format
	logger.Info("adding target", "host", host, "type", ct)
	cs.configs = append(cs.configs, target{host: host, ct: ct, t: time.Now()})
	return nil
}

func (cs *ConfigStore) RemoveTarget(string) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
}

// TODO: Rethink this
func (cs *ConfigStore) JSON() ([]byte, error) {
	cs.mu.RLock()
	defer cs.mu.RUnlock()

	m := make(map[pa.ConfigType][]target)
	for _, v := range cs.configs {
		m[v.ct] = append(m[v.ct], v)
	}

	j := []staticConfig{}
	for i, v := range m {
		j = append(j, staticConfig{})
		for _, t := range v {
			j[i].Targets = append(j[i].Targets, t.host)
		}
	}

	return json.Marshal(j)
}

func (cs *ConfigStore) purge(t time.Time) {
	logger.Debug("checking purges", "time", t, "configs", cs.configs)
	cs.mu.Lock()
	defer cs.mu.Unlock()

	for i, tar := range cs.configs {
		if t.Sub(tar.t) > cs.lifetime {
			logger.Info("purging target", "target", tar)
			if len(cs.configs) == 1 {
				cs.configs = []target{}
				return
			}

			cs.configs[i] = cs.configs[len(cs.configs)-1]
			cs.configs = cs.configs[:len(cs.configs)-1]
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
