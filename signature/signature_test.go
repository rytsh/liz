package signature

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFunc(t *testing.T) {
	type args struct {
		fn   any
		opts []Option
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "own func",
			args: args{
				fn: Func,
				opts: []Option{
					WithFuncName("Func"),
					WithArgsNameHook(func(index int, argType reflect.Type, isVariadic bool) (string, bool) {
						return fmt.Sprintf("arg%d_%t", index, isVariadic), true
					}),
					WithReturn(true),
				},
			},
			want: "Func(arg0_false interface {}, arg1_true ...signature.Option) string",
		},
		{
			name: "Println",
			args: args{
				fn: fmt.Println,
				opts: []Option{
					WithFuncName("Println"),
					WithArgsNameHook(func(index int, argType reflect.Type, isVariadic bool) (string, bool) {
						return "a", true
					}),
					WithReturn(true),
					WithReturnArgsNameHook(func(index int, argType reflect.Type) (string, bool) {
						switch index {
						case 0:
							return "n", true
						case 1:
							return "err", true
						}

						return "", false
					}),
					WithTypeHook(func(argType reflect.Type, isVariadic bool) (string, bool) {
						v := argType.String()
						if v == "interface {}" {
							if isVariadic {
								return "...any", true
							}

							return "any", true
						}

						return "", false
					}),
				},
			},
			want: "Println(a ...any) (n int, err error)",
		},
		{
			name: "test func",
			args: args{
				fn: func(a int, b string, c ...interface{}) (d int, e error) {
					return
				},
				opts: []Option{
					WithArgsNameHook(func(index int, argType reflect.Type, isVariadic bool) (string, bool) {
						switch index {
						case 0:
							return "a", true
						case 1:
							return "b", true
						case 2:
							return "c", true
						}

						return "", false
					}),
					WithReturn(true),
					WithReturnArgsNameHook(func(index int, argType reflect.Type) (string, bool) {
						switch index {
						case 0:
							return "d", true
						case 1:
							return "e", true
						}

						return "", false
					}),
					WithTypeHook(func(argType reflect.Type, isVariadic bool) (string, bool) {
						// check argType is empty interface
						if argType.String() == "interface {}" {
							if isVariadic {
								return "...interface{}", true
							}

							return "interface{}", true
						}

						return "", false
					}),
				},
			},
			want: "(a int, b string, c ...interface{}) (d int, e error)",
		},
		{
			name: "test func without return args",
			args: args{
				fn: func(a int, b string, c ...interface{}) (d int, e error) {
					return
				},
				opts: []Option{
					WithArgsNameHook(func(index int, argType reflect.Type, isVariadic bool) (string, bool) {
						switch index {
						case 0:
							return "a", true
						case 1:
							return "b", true
						case 2:
							return "c", true
						}

						return "", false
					}),
					WithReturn(true),
					WithTypeHook(func(argType reflect.Type, isVariadic bool) (string, bool) {
						if argType.String() == "interface {}" {
							if isVariadic {
								return "...interface{}", true
							}

							return "interface{}", true
						}

						return "", false
					}),
				},
			},
			want: "(a int, b string, c ...interface{}) (int, error)",
		},
		{
			name: "test func without return",
			args: args{
				fn: func(a int, b string, c ...interface{}) (d int, e error) {
					return
				},
				opts: []Option{
					WithArgsNameHook(func(index int, argType reflect.Type, isVariadic bool) (string, bool) {
						switch index {
						case 0:
							return "a", true
						case 1:
							return "b", true
						case 2:
							return "c", true
						}

						return "", false
					}),
					WithTypeHook(func(argType reflect.Type, isVariadic bool) (string, bool) {
						if argType.String() == "interface {}" {
							if isVariadic {
								return "...interface{}", true
							}

							return "interface{}", true
						}

						return "", false
					}),
				},
			},
			want: "(a int, b string, c ...interface{})",
		},
		{
			name: "test func without arg types",
			args: args{
				fn: func(a int, b string, c ...interface{}) (d int, e error) {
					return
				},
				opts: []Option{
					WithFuncName("Func"),
					WithArgsNameHook(func(index int, argType reflect.Type, isVariadic bool) (string, bool) {
						switch index {
						case 0:
							return "a", true
						case 1:
							return "b", true
						case 2:
							return "c", true
						}

						return "", false
					}),
					WithArgsTypeHook(func(index int, argType reflect.Type, isVariadic bool) (string, bool) {
						return "", true
					}),
					WithTypeHook(func(argType reflect.Type, isVariadic bool) (string, bool) {
						if argType.String() == "interface {}" {
							if isVariadic {
								return "...interface{}", true
							}

							return "interface{}", true
						}

						return "", false
					}),
				},
			},
			want: "Func(a, b, c)",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Func(tt.args.fn, tt.args.opts...); got != tt.want {
				t.Errorf("%v, want %v", got, tt.want)
			}
		})
	}
}
