package logit

import "sync"

type Tags struct {
	mutex   sync.RWMutex
	entries map[string]any
}

func (t *Tags) add(name string, value any) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.entries[name] = value
}

func (t *Tags) copy() *Tags {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	c := Tags{
		mutex:   sync.RWMutex{},
		entries: map[string]any{},
	}
	for k, v := range t.entries {
		c.entries[k] = v
	}
	return &c
}

func (t *Tags) Retrieve() map[string]any {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	tm := map[string]any{}
	for k, v := range t.entries {
		tm[k] = v
	}
	return tm
}
