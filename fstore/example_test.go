package fstore

import (
	"fmt"

	"github.com/rytsh/liz/templatex"
	"github.com/rytsh/liz/templatex/store"
)

func Example() {
	tpl := templatex.New(store.WithAddFuncsTpl(
		FuncMapTpl(),
	))

	output, err := tpl.ExecuteBuffer(
		templatex.WithContent(
			`{{ $v := codec.JsonDecode (codec.StringToByte .) }}{{ $v.data.name }}`,
		),
		templatex.WithData(`{"data": {"name": "Hatay"}}`),
	)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}

	fmt.Printf("%s", output)
	// Output:
	// Hatay
}
