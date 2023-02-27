package fstore

import (
	"github.com/Masterminds/sprig/v3"
	"github.com/rytsh/liz/utils/fstore/generic"
	"github.com/rytsh/liz/utils/templatex"
	"github.com/rytsh/liz/utils/templatex/store"

	_ "github.com/rytsh/liz/utils/fstore/funcs"
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
			AddArgument("template", t).
			AddArgument("workDir", o.workDir)

		for _, fName := range generic.CallReg.GetFunctionNames() {
			if _, ok := disabled[fName]; ok {
				continue
			}

			storeHolder.AddFunc(fName, generic.GetFunc(fName))
		}

		return storeHolder.Funcs()
	}
}
