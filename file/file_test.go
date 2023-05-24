package file

import (
	"path"
	"reflect"
	"testing"
)

func TestAPI_LoadRaw(t *testing.T) {
	a := New()

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test file",
			args: args{
				path: "testdata/test.txt",
			},
			want: []byte(`Nisi eu cupidatat dolore sint.
Laborum ex eiusmod velit fugiat eu elit ea sunt Lorem est.
`),
		},
		{
			name: "non exist file",
			args: args{
				path: "testdata/nonexistfile",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := a.LoadRaw(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("API.LoadRaw() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("API.LoadRaw() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAPI_Load(t *testing.T) {
	a := New()

	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name: "json load",
			args: args{
				path: "testdata/test.json",
			},
			want: map[string]interface{}{
				"foo": "bar",
				"bar": map[string]interface{}{
					"foo": float64(1234),
				},
			},
		},
		{
			name: "yaml load",
			args: args{
				path: "testdata/test.yml",
			},
			want: map[string]interface{}{
				"foo": "bar",
				"bar": map[string]interface{}{
					"foo": 1234,
				},
			},
		},
		{
			name: "toml load",
			args: args{
				path: "testdata/test.toml",
			},
			want: map[string]interface{}{
				"foo": "bar",
				"bar": map[string]interface{}{
					"foo": int64(1234),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var v interface{}
			err := a.Load(tt.args.path, &v)
			if (err != nil) != tt.wantErr {
				t.Errorf("API.LoadMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(v, tt.want) {
				t.Errorf("API.LoadMap() = %v, want %v", v, tt.want)
			}
		})
	}
}

func TestAPI_SetMap(t *testing.T) {
	a := New()

	type args struct {
		path string
		data map[string]interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		compare []byte
	}{
		{
			name: "json set",
			args: args{
				path: "test.json",
				data: map[string]interface{}{
					"foo": "bar",
					"bar": map[string]interface{}{
						"foo": 1234,
					},
				},
			},
			compare: []byte(`{
  "bar": {
    "foo": 1234
  },
  "foo": "bar"
}
`),
		},
		{
			name: "yaml set",
			args: args{
				path: "test.yml",
				data: map[string]interface{}{
					"foo": "bar",
					"bar": map[string]interface{}{
						"foo": 1234,
					},
				},
			},
			compare: []byte(`bar:
    foo: 1234
foo: bar
`),
		},
		{
			name: "toml set",
			args: args{
				path: "dummy/test.toml",
				data: map[string]interface{}{
					"foo": "bar",
					"bar": map[string]interface{}{
						"foo": 1234,
					},
				},
			},
			compare: []byte(`foo = "bar"

[bar]
  foo = 1234
`),
		},
	}

	tempdir := t.TempDir()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := path.Join(tempdir, tt.args.path)

			if err := a.SetWithCodec(filePath, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("API.SetWithCodec() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			got, err := a.LoadRaw(filePath)
			if err != nil {
				t.Errorf("API.SetMap() error = %v", err)
			}

			if !reflect.DeepEqual(string(got), string(tt.compare)) {
				t.Errorf("API.SetMap() = \n%s\nwant \n%s", got, tt.compare)
			}
		})
	}
}

func TestAPI_Set(t *testing.T) {
	a := New()

	type args struct {
		path string
		data []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		compare []byte
	}{
		{
			name: "raw set",
			args: args{
				path: "dummy/test.json",
				data: []byte(`anyting can be here`),
			},
			compare: []byte(`anyting can be here`),
			wantErr: false,
		},
	}

	tempdir := t.TempDir()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filePath := path.Join(tempdir, tt.args.path)

			if err := a.SetRaw(filePath, tt.args.data, WithFilePerm("0600")); (err != nil) != tt.wantErr {
				t.Errorf("API.SetRaw() error = %v, wantErr %v", err, tt.wantErr)
			}

			got, err := a.LoadRaw(filePath)
			if err != nil {
				t.Errorf("API.Set() error = %v", err)
			}

			if !reflect.DeepEqual(string(got), string(tt.compare)) {
				t.Errorf("API.Set() = %s, want %s", got, tt.compare)
			}
		})
	}
}

func TestAPI_LoadWithCodec(t *testing.T) {
	type args struct {
		path  string
		dst   interface{}
		codec Codec
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		want    interface{}
	}{
		{
			name: "json load",
			args: args{
				path:  "testdata/test.json",
				dst:   map[string]interface{}{},
				codec: nil,
			},
			wantErr: false,
			want: map[string]interface{}{
				"foo": "bar",
				"bar": map[string]interface{}{
					"foo": float64(1234),
				},
			},
		},
		{
			name: "yaml load",
			args: args{
				path:  "testdata/test.yml",
				dst:   map[string]interface{}{},
				codec: nil,
			},
			wantErr: false,
			want: map[string]interface{}{
				"foo": "bar",
				"bar": map[string]interface{}{
					"foo": 1234,
				},
			},
		},
		{
			name: "toml load",
			args: args{
				path:  "testdata/test.toml",
				dst:   map[string]interface{}{},
				codec: nil,
			},
			wantErr: false,
			want: map[string]interface{}{
				"foo": "bar",
				"bar": map[string]interface{}{
					"foo": int64(1234),
				},
			},
		},
		{
			name: "raw load",
			args: args{
				path:  "testdata/test.txt",
				dst:   []byte{},
				codec: nil,
			},
			wantErr: false,
			want: []byte(`Nisi eu cupidatat dolore sint.
Laborum ex eiusmod velit fugiat eu elit ea sunt Lorem est.
`),
		},
	}

	a := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := a.LoadWithCodec(tt.args.path, &tt.args.dst, tt.args.codec); (err != nil) != tt.wantErr {
				t.Errorf("API.LoadWithCodec() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(tt.args.dst, tt.want) {
				t.Errorf("API.LoadWithCodec() = %#v, want %#v", tt.args.dst, tt.want)
			}
		})
	}
}

func TestAPI_LoadWithReturn(t *testing.T) {
	type args struct {
		path  string
		codec Codec
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "json load",
			args: args{
				path:  "testdata/test.json",
				codec: nil,
			},
			wantErr: false,
			want: map[string]interface{}{
				"foo": "bar",
				"bar": map[string]interface{}{
					"foo": float64(1234),
				},
			},
		},
		{
			name: "raw load",
			args: args{
				path:  "testdata/test.txt",
				codec: nil,
			},
			wantErr: false,
			want: []byte(`Nisi eu cupidatat dolore sint.
Laborum ex eiusmod velit fugiat eu elit ea sunt Lorem est.
`),
		},
	}
	a := New()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := a.LoadWithReturn(tt.args.path, tt.args.codec)
			if (err != nil) != tt.wantErr {
				t.Errorf("API.LoadWithReturn() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("API.LoadWithReturn() = %v, want %v", got, tt.want)
			}
		})
	}
}
