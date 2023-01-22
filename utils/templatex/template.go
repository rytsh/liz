package templatex

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	textTemplate "text/template"

	"github.com/rytsh/liz/utils/templatex/functions"
)

type Template struct {
	template *textTemplate.Template
	funcs    map[string]interface{}
}

func New(opts ...functions.Option) *Template {
	tpl := &Template{
		template: textTemplate.New("txt"),
	}

	tpl.setFunctions(opts...)

	return tpl
}

func (t *Template) ListFunctions() {
	funcs := FuncInfos(t.funcs)
	sort.Slice(funcs, func(i, j int) bool {
		return funcs[i].Name < funcs[j].Name
	})

	for _, v := range funcs {
		fmt.Println(v.Description)
	}
}

func (t *Template) Reset() {
	t.template = textTemplate.New("txt")
}

// SetDelims sets the template delimiters to the specified strings
// and returns the template to allow chaining.
func (t *Template) SetDelims(left, right string) *Template {
	if left == "" {
		left = "{{"
	}

	if right == "" {
		right = "}}"
	}

	t.template.Delims(left, right)

	return t
}

func (t *Template) setFunctions(opts ...functions.Option) {
	t.funcs = functions.New(opts...).Funcs()
	t.template.Funcs(t.funcs)
}

func (t *Template) ParseGlob(pattern string) (*Template, error) {
	if pattern == "" {
		return t, nil
	}

	tpl, err := t.template.ParseGlob(pattern)
	if err != nil {
		return nil, fmt.Errorf("Parse error: %w", err)
	}

	t.template = tpl

	return t, nil
}

func (t *Template) ExecuteBytes(v any, content string) ([]byte, error) {
	output, err := t.execute(v, content)
	if err != nil {
		return output.Bytes(), err
	}

	return output.Bytes(), nil
}

func (t *Template) Execute(v any, content string) (string, error) {
	output, err := t.execute(v, content)
	if err != nil {
		return "", err
	}

	return output.String(), nil
}

func (t *Template) execute(v any, content string) (*bytes.Buffer, error) {
	var b bytes.Buffer

	tpl, err := t.template.Clone()
	if err != nil {
		return nil, fmt.Errorf("execute clone error: %w", err)
	}
	// Execute the template and write the output to the buffer
	if err := textTemplate.Must(tpl.Parse(content)).Execute(&b, v); err != nil {
		return nil, fmt.Errorf("Execute error: %w", err)
	}

	return &b, nil
}

func (t *Template) ExecuteContent(writer io.Writer, v any, content []byte) error {
	tpl, err := t.template.Clone()
	if err != nil {
		return fmt.Errorf("ExecuteContent clone error: %w", err)
	}

	// Execute the template and write the output to the buffer
	if err := textTemplate.Must(tpl.Parse(string(content))).Execute(writer, v); err != nil {
		return fmt.Errorf("ExecuteContent error: %w", err)
	}

	return nil
}
