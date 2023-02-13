package templatex

import "io"

type options struct {
	writer   io.Writer
	content  string
	template string
	values   any
	parsed   bool
}

type Option func(options *options)

func WithIO(w io.Writer) Option {
	return func(options *options) {
		options.writer = w
	}
}

func WithContent(content string) Option {
	return func(options *options) {
		options.content = content
	}
}

func WithTemplate(template string) Option {
	return func(options *options) {
		options.template = template
	}
}

func WithData(values any) Option {
	return func(options *options) {
		options.values = values
	}
}

func WithParsed(parsed bool) Option {
	return func(options *options) {
		options.parsed = parsed
	}
}
