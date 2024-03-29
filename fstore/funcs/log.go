package funcs

import (
	"github.com/rs/zerolog/log"
	"github.com/rytsh/liz/fstore/generic"
)

func init() {
	generic.CallReg.
		AddFunction("log", generic.ReturnWithFn(Log{})).
		AddFunction("nothing", generic.ReturnWithFn(Nothing))
}

type Log struct{}

func (Log) Debug(v any) any {
	log.Debug().Msgf("%v\n", v)
	return v
}

func (Log) Info(v any) any {
	log.Info().Msgf("%v\n", v)
	return v
}

func (Log) Warn(v any) any {
	log.Warn().Msgf("%v\n", v)
	return v
}

func (Log) Error(v any) any {
	log.Error().Msgf("%v\n", v)
	return v
}

func (Log) Fatal(v any) any {
	log.Fatal().Msgf("%v\n", v)
	return v
}

func (Log) Panic(v any) any {
	log.Panic().Msgf("%v\n", v)
	return v
}

func Nothing(v any) string {
	return ""
}
