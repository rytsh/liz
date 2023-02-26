package templatex

import (
	"bytes"
	"fmt"
	"sort"
	textTemplate "text/template"

	"github.com/rytsh/liz/utils/templatex/store"
)

type Template struct {
	template       *textTemplate.Template
	templateParsed *textTemplate.Template
	funcs          map[string]interface{}
}

func New(opts ...store.Option) *Template {
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

func (t *Template) setFunctions(opts ...store.Option) {
	optsNew := make([]store.Option, 0, len(opts)+1)
	optsNew = append(optsNew, store.WithFnValue(t))
	optsNew = append(optsNew, opts...)

	t.funcs = store.New(optsNew...).Funcs()
	t.template.Funcs(t.funcs)
}

func (t *Template) AddFunc(funcMap textTemplate.FuncMap) {
	for k, v := range funcMap {
		t.funcs[k] = v
	}

	t.template.Funcs(funcMap)
}

func (t *Template) ParseGlob(pattern string) error {
	if pattern == "" {
		return nil
	}

	tpl, err := t.template.ParseGlob(pattern)
	if err != nil {
		return fmt.Errorf("Parse error: %w", err)
	}

	t.template = tpl

	return nil
}

// Parse content and set new template to parsed.
func (t *Template) Parse(content string) error {
	tpl, err := t.template.Clone()
	if err != nil {
		return fmt.Errorf("execute clone error: %w", err)
	}

	// Execute the template and write the output to the buffer
	tpl, err = tpl.Parse(content)
	if err != nil {
		return fmt.Errorf("Parse error: %w", err)
	}

	t.templateParsed = tpl

	return nil
}

func (t *Template) Execute(opts ...Option) error {
	o := &options{
		writer: &bytes.Buffer{},
	}
	for _, opt := range opts {
		opt(o)
	}

	err := t.execute(o)
	if err != nil {
		return err
	}

	return nil
}

func (t *Template) ExecuteBuffer(opts ...Option) ([]byte, error) {
	output := &bytes.Buffer{}
	o := &options{}
	for _, opt := range opts {
		opt(o)
	}

	o.writer = output

	err := t.execute(o)
	if err != nil {
		return nil, err
	}

	return output.Bytes(), nil
}

func (t *Template) execute(o *options) error {
	tpl := t.template
	if o.parsed && t.templateParsed != nil {
		tpl = t.templateParsed
	}

	tpl, err := tpl.Clone()
	if err != nil {
		return fmt.Errorf("execute clone error: %w", err)
	}

	parsedTpl := tpl
	if !o.parsed {
		parsedTpl, err = tpl.Parse(o.content)
		if err != nil {
			return fmt.Errorf("execute parse error: %w", err)
		}
	}

	// Execute the template and write the output to the buffer
	if o.template != "" {
		err = parsedTpl.ExecuteTemplate(o.writer, o.template, o.values)
	} else {
		err = parsedTpl.Execute(o.writer, o.values)
	}

	if err != nil {
		return fmt.Errorf("Execute error: %w", err)
	}

	return nil
}
