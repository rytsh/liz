package functions

import (
	"sync"
	textTemplate "text/template"

	"github.com/Masterminds/sprig/v3"

	"github.com/rytsh/liz/utils/templatex/functions/custom"
	"github.com/rytsh/liz/utils/templatex/functions/hugo"
	"github.com/rytsh/liz/utils/templatex/functions/humanize"
)

type Holder struct {
	funcMap    map[string]interface{}
	mutex      sync.RWMutex
	initialize sync.Once
}

func New(opts ...Option) *Holder {
	return new(Holder).InitializeFuncs(opts...)
}

func (h *Holder) InitializeFuncs(opts ...Option) *Holder {
	h.initialize.Do(func() {
		if h.funcMap == nil {
			h.funcMap = make(map[string]interface{})
		}

		option := &options{}
		for _, opt := range opts {
			opt(option)
		}

		h.AddFuncs(
			sprig.GenericFuncMap(),
			hugo.FuncMapFn(option.workDir)(),
			humanize.FuncMap(),
			custom.FuncMap(),
			// Add additonal functions here
			option.addFuncs,
		)

		for _, f := range option.disableFuncs {
			delete(h.funcMap, f)
		}
	})

	return h
}

func (h *Holder) Funcs() map[string]interface{} {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	return h.funcMap
}

func (h *Holder) TxtFuncs() textTemplate.FuncMap {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	return h.funcMap
}

func (h *Holder) AddFuncs(funcs ...map[string]interface{}) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for _, f := range funcs {
		for k, v := range f {
			h.funcMap[k] = v
		}
	}
}

type options struct {
	workDir      string
	disableFuncs []string
	addFuncs     map[string]interface{}
}

type Option func(*options)

func WorkDir(workDir string) Option {
	return func(o *options) {
		o.workDir = workDir
	}
}

func DisableFuncs(funcs ...string) Option {
	return func(o *options) {
		o.disableFuncs = append(o.disableFuncs, funcs...)
	}
}

func AddFuncs(funcs map[string]interface{}) Option {
	return func(o *options) {
		if o.addFuncs == nil {
			o.addFuncs = make(map[string]interface{}, len(funcs))
		}

		for k, v := range funcs {
			o.addFuncs[k] = v
		}
	}
}

func AddFunc(key string, f interface{}) Option {
	return func(o *options) {
		if o.addFuncs == nil {
			o.addFuncs = make(map[string]interface{}, 1)
		}

		o.addFuncs[key] = f
	}
}
