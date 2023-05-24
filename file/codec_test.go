package file

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

func TestCodec_Decode(t *testing.T) {
	type args struct {
		r io.Reader
		v any
	}
	tests := []struct {
		name    string
		c       Codec
		args    args
		wantErr bool
		want    any
	}{
		{
			name: "raw decode",
			c:    RAW{},
			args: args{
				r: bytes.NewReader([]byte("foo")),
				v: &[]byte{},
			},
			wantErr: false,
			want:    []byte("foo"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.Decode(tt.args.r, tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("Codec.Decode() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			rawType := reflect.ValueOf(tt.args.v)

			if rawType.Kind() == reflect.Ptr {
				value := rawType.Elem().Interface()
				if !reflect.DeepEqual(value, tt.want) {
					t.Errorf("Codec.Decode() = %v, want %v", value, tt.want)
				}
			}

			if reflect.DeepEqual(tt.args.v, tt.want) {
				t.Errorf("Codec.Decode() = %v, want %v", tt.args.v, tt.want)
			}
		})
	}
}
