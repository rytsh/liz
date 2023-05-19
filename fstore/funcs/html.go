package funcs

import (
	"html"

	"github.com/rytsh/liz/fstore/generic"
)

func init() {
	generic.CallReg.AddFunction("html", generic.ReturnWithFn(Html{}))
}

type Html struct{}

func (Html) EscapeString(v string) string {
	return html.EscapeString(v)
}

func (Html) UnescapeString(v string) string {
	return html.UnescapeString(v)
}
