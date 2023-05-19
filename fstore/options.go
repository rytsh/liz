package fstore

import (
	"github.com/rytsh/liz/templatex"
	"github.com/worldline-go/logz"
)

type options struct {
	disableFuncs []string
	trust        bool
	log          logz.Adapter
	workDir      string
	templatex    *templatex.Template
	specificFunc []string
}

type Option func(options *options)

func WithSpecificFuncs(specificFuncs ...string) Option {
	return func(options *options) {
		if len(options.specificFunc) > 0 {
			options.specificFunc = append(options.specificFunc, specificFuncs...)
		} else {
			options.specificFunc = specificFuncs
		}
	}
}

func WithDisableFuncs(disableFuncs ...string) Option {
	return func(options *options) {
		if len(options.disableFuncs) > 0 {
			options.disableFuncs = append(options.disableFuncs, disableFuncs...)
		} else {
			options.disableFuncs = disableFuncs
		}
	}
}

func WithTrust(trust bool) Option {
	return func(options *options) {
		options.trust = trust
	}
}

func WithWorkDir(workDir string) Option {
	return func(options *options) {
		options.workDir = workDir
	}
}

func WithTemplatex(t *templatex.Template) Option {
	return func(options *options) {
		options.templatex = t
	}
}

func WithLog(log logz.Adapter) Option {
	return func(options *options) {
		options.log = log
	}
}
