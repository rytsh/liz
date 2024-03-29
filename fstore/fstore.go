package fstore

import (
	"github.com/Masterminds/sprig/v3"
	"github.com/rytsh/liz/fstore/generic"
	"github.com/rytsh/liz/templatex"
	"github.com/rytsh/liz/templatex/store"

	_ "github.com/rytsh/liz/fstore/funcs"
)

func FuncMap(opts ...Option) map[string]interface{} {
	opt := optionRun(opts...)
	return funcX(opt)(opt.templatex)
}

func FuncMapTpl(opts ...Option) func(t *templatex.Template) map[string]interface{} {
	return funcX(optionRun(opts...))
}

func optionRun(opts ...Option) options {
	opt := options{}
	for _, o := range opts {
		o(&opt)
	}

	return opt
}

func funcX(o options) func(t *templatex.Template) map[string]interface{} {
	return func(t *templatex.Template) map[string]interface{} {
		disabled := make(map[string]struct{}, len(o.disableFuncs))
		for _, d := range o.disableFuncs {
			disabled[d] = struct{}{}
		}

		storeHolder := store.Holder{}
		storeHolder.AddFuncs(sprig.GenericFuncMap())

		// custom functions
		generic.CallReg.
			AddArgument("trust", o.trust).
			AddArgument("log", nil).
			AddArgument("template", t).
			AddArgument("workDir", o.workDir)

		for name, v := range definedFuncMaps() {
			if _, ok := disabled[name]; ok {
				continue
			}

			if !checkSpecificFunc(o, name) {
				continue
			}

			storeHolder.AddFuncs(v)
		}

		for _, fName := range generic.CallReg.GetFunctionNames() {
			if _, ok := disabled[fName]; ok {
				continue
			}

			if !checkSpecificFunc(o, fName) {
				continue
			}

			storeHolder.AddFunc(fName, generic.GetFunc(fName))
		}

		return storeHolder.Funcs()
	}
}

func checkSpecificFunc(o options, name string) bool {
	if len(o.specificFunc) > 0 {
		foundName := false
		for _, s := range o.specificFunc {
			if s == name {
				foundName = true
				break
			}
		}

		return foundName
	}

	return true
}
