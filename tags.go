package logit

import "sync"

type tags struct {
	mutex sync.RWMutex
	tag   map[string]any
}

func (t *tags) add(name string, value any) {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	t.tag[name] = value
}

func (t *tags) copy() *tags {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	c := tags{
		mutex: sync.RWMutex{},
		tag:   map[string]any{},
	}
	for k, v := range t.tag {
		c.tag[k] = v
	}
	return &c
}

func (t *tags) retrieve() map[string]any {
	t.mutex.RLock()
	defer t.mutex.RUnlock()
	tm := map[string]any{}
	for k, v := range t.tag {
		tm[k] = v
	}
	return tm
}
