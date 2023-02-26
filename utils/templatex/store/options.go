package store

type options struct {
	disableFuncs []string
	addFuncs     map[string]interface{}
	fnValue      interface{}
}

type Option func(*options)

func WithDisableFuncs(funcs ...string) Option {
	return func(o *options) {
		o.disableFuncs = append(o.disableFuncs, funcs...)
	}
}

func WithAddFuncs(funcs map[string]interface{}) Option {
	return func(o *options) {
		if o.addFuncs == nil {
			o.addFuncs = make(map[string]interface{}, len(funcs))
		}

		for k, v := range funcs {
			o.addFuncs[k] = v
		}
	}
}

func WithAddFunc(key string, f interface{}) Option {
	return func(o *options) {
		if o.addFuncs == nil {
			o.addFuncs = make(map[string]interface{}, 1)
		}

		o.addFuncs[key] = f
	}
}

func WithFnValue[T any](fn T) Option {
	return func(o *options) {
		o.fnValue = fn
	}
}

func WithAddFuncsTpl[T any](fn func(T) map[string]interface{}) Option {
	return func(o *options) {
		funcs := fn(o.fnValue.(T))

		if o.addFuncs == nil {
			o.addFuncs = make(map[string]interface{}, len(funcs))
		}

		for k, v := range funcs {
			o.addFuncs[k] = v
		}
	}
}

func WithAddFuncTpl[T any](key string, f func(T) interface{}) Option {
	return func(o *options) {
		if o.addFuncs == nil {
			o.addFuncs = make(map[string]interface{}, 1)
		}

		o.addFuncs[key] = f(o.fnValue.(T))
	}
}
