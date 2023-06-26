package signature

import (
	"reflect"
	"strings"
)

func getOpts(opts []Option) *options {
	opt := &options{}

	for _, o := range opts {
		o(opt)
	}

	return opt
}

func Func(fn any, opts ...Option) string {
	opt := getOpts(opts)

	funcName := opt.funcName

	// get function signature with reflect as string
	reflectFn := reflect.ValueOf(fn)
	reflectFnType := reflectFn.Type()

	isVariadic := reflectFnType.IsVariadic()
	reflectFnParams := reflectFnType.NumIn()
	reflectFnResults := reflectFnType.NumOut()

	// get function description
	var description strings.Builder
	description.WriteString(funcName + "(")
	if reflectFnParams > 0 {
		for i := 0; i < reflectFnParams; i++ {
			argumentType := reflectFnType.In(i)
			variadic := isVariadic && i == reflectFnParams-1
			if variadic {
				argumentType = argumentType.Elem()
			}

			space := argName(i, argumentType, variadic, &description, opt)
			funcType(false, i, argumentType, variadic, &description, opt, space)

			if i < reflectFnParams-1 {
				description.WriteString(", ")
			}
		}
	}

	description.WriteString(")")

	if opt.returnSign && reflectFnResults > 0 {
		if reflectFnResults > 1 {
			description.WriteString(" (")
		} else if reflectFnResults == 1 {
			description.WriteString(" ")
		}

		for i := 0; i < reflectFnResults; i++ {
			argumentType := reflectFnType.Out(i)

			space := argNameReturn(i, argumentType, false, &description, opt)
			funcType(true, i, argumentType, false, &description, opt, space)

			if i < reflectFnResults-1 {
				description.WriteString(", ")
			}
		}

		if reflectFnResults > 1 {
			description.WriteString(")")
		}
	}

	return description.String()
}

func funcType(returnType bool, i int, argumentType reflect.Type, isVariadic bool, description *strings.Builder, opt *options, space bool) {
	var argumentTypeStr string

	changed := false
	if returnType {
		if opt.returnArgsTypeHook != nil {
			if v, ok := opt.returnArgsTypeHook(i, argumentType); ok {
				argumentTypeStr = v
				changed = true
			}
		}
	} else if opt.argsTypeHook != nil {
		if v, ok := opt.argsTypeHook(i, argumentType, isVariadic); ok {
			argumentTypeStr = v
			changed = true
		}
	}
	if !changed && opt.typeHook != nil {
		if v, ok := opt.typeHook(argumentType, isVariadic); ok {
			argumentTypeStr = v
			changed = true
		}
	}
	if !changed {
		argumentTypeStr = argumentType.String()
		if space {
			description.WriteString(" ")
		}
		if isVariadic {
			description.WriteString("...")
		}
	} else {
		if space && argumentTypeStr != "" {
			description.WriteString(" ")
		}
	}

	description.WriteString(argumentTypeStr)
}

func argName(i int, argumentType reflect.Type, isVariadic bool, description *strings.Builder, opt *options) bool {
	space := false
	if opt.argsNameHook != nil {
		if argName, ok := opt.argsNameHook(i, argumentType, isVariadic); ok {
			space = true
			description.WriteString(argName)
		}
	}

	return space
}

func argNameReturn(i int, argumentType reflect.Type, isVariadic bool, description *strings.Builder, opt *options) bool {
	space := false
	if opt.returnArgsNameHook != nil {
		if argName, ok := opt.returnArgsNameHook(i, argumentType); ok {
			space = true
			description.WriteString(argName)
		}
	}

	return space
}
