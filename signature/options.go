package signature

import "reflect"

type options struct {
	funcName           string
	argsNameHook       func(index int, argType reflect.Type, isVariadic bool) (string, bool)
	argsTypeHook       func(index int, argType reflect.Type, isVariadic bool) (string, bool)
	returnArgsNameHook func(index int, argType reflect.Type) (string, bool)
	returnArgsTypeHook func(index int, argType reflect.Type) (string, bool)
	returnSign         bool

	typeHook func(argType reflect.Type, isVariadic bool) (string, bool)
}

// Option to execute the template.
type Option func(options *options)

// WithFuncName to set the function name.
func WithFuncName(funcName string) Option {
	return func(options *options) {
		options.funcName = funcName
	}
}

// WithArgsNameHook to set the args name hook.
// If the hook returns false, the default name will be used.
func WithArgsNameHook(fn func(index int, argType reflect.Type, isVariadic bool) (string, bool)) Option {
	return func(options *options) {
		options.argsNameHook = fn
	}
}

// WithArgsTypeHook to set the args type hook.
// If the hook returns false, the default type will be used.
func WithArgsTypeHook(fn func(index int, argType reflect.Type, isVariadic bool) (string, bool)) Option {
	return func(options *options) {
		options.argsTypeHook = fn
	}
}

// WithReturnArgsNameHook to set the return args name hook.
// If the hook returns false, the default name will be used.
func WithReturnArgsNameHook(fn func(index int, argType reflect.Type) (string, bool)) Option {
	return func(options *options) {
		options.returnArgsNameHook = fn
	}
}

// WithReturnArgsTypeHook to set the return args type hook.
// If the hook returns false, the default type will be used.
func WithReturnArgsTypeHook(fn func(index int, argType reflect.Type) (string, bool)) Option {
	return func(options *options) {
		options.returnArgsTypeHook = fn
	}
}

// WithReturn to set the return sign.
func WithReturn(returnSign bool) Option {
	return func(options *options) {
		options.returnSign = returnSign
	}
}

func WithTypeHook(fn func(argType reflect.Type, isVariadic bool) (string, bool)) Option {
	return func(options *options) {
		options.typeHook = fn
	}
}
