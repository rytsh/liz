package store

import (
	"sync"
)

type Holder struct {
	funcMap map[string]interface{}
	mutex   sync.RWMutex
}

func New(opts ...Option) *Holder {
	return new(Holder).initializeFuncs(opts...)
}

func (h *Holder) setFuncMap() {
	if h.funcMap == nil {
		h.funcMap = make(map[string]interface{})
	}
}

func (h *Holder) initializeFuncs(opts ...Option) *Holder {
	h.setFuncMap()

	option := &options{}
	for _, opt := range opts {
		opt(option)
	}

	h.AddFuncs(
		// Add additonal functions here
		option.addFuncs,
	)

	for _, f := range option.disableFuncs {
		delete(h.funcMap, f)
	}

	return h
}

func (h *Holder) Funcs() map[string]interface{} {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	h.setFuncMap()

	return h.funcMap
}

func (h *Holder) AddFuncs(funcs ...map[string]interface{}) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.setFuncMap()

	for _, f := range funcs {
		for k, v := range f {
			h.funcMap[k] = v
		}
	}
}

func (h *Holder) AddFunc(name string, fn interface{}) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.setFuncMap()

	h.funcMap[name] = fn
}
