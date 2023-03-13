package fstore

import (
	"github.com/rytsh/liz/utils/templatex"
	"github.com/worldline-go/logz"
)

type options struct {
	disableFuncs []string
	trust        bool
	log          logz.Adapter
	workDir      string
	templatex    *templatex.Template
}

type Option func(options *options)

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
