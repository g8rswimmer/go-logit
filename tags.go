package logit

import "sync"

type Tags struct {
	mutex sync.RWMutex
	tag   map[string]any
}

func (t *Tags) add(name string, value any) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.tag[name] = value
}

func (t *Tags) copy() *Tags {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	c := Tags{
		mutex: sync.RWMutex{},
		tag:   map[string]any{},
	}
	for k, v := range t.tag {
		c.tag[k] = v
	}
	return &c
}

func (t *Tags) Retrieve() map[string]any {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	tm := map[string]any{}
	for k, v := range t.tag {
		tm[k] = v
	}
	return tm
}
