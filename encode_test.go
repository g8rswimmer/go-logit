package logit

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testAttributes struct {
	Attr1         string `logit:"attribute_one"`
	Attr2         int    `logit:",omitempty"`
	Attr3         []string
	notGoingToLog string
}

func Test_encodeMap(t *testing.T) {
	type args struct {
		input any
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]any
		wantErr bool
	}{
		{
			name: "Testing map encoding",
			args: args{
				input: map[string]any{
					"key1": "value1",
					"key2": 123,
					"key3": []any{"1", "2"},
					"key4": map[string]any{
						"sub_key1": "sub value1",
						"sub_key2": 3455,
					},
					"key5": testAttributes{
						Attr1:         "struct attribute",
						Attr2:         99887,
						Attr3:         []string{"elem1", "elem2"},
						notGoingToLog: "this should not be logged",
					},
					"key6": []testAttributes{
						{
							Attr1:         "struct attribute array 1",
							Attr2:         9988455557,
							Attr3:         []string{"elem1_1", "elem2_1"},
							notGoingToLog: "this should not be logged",
						},
						{
							Attr1:         "struct attribute array 2",
							Attr2:         998444487,
							Attr3:         []string{"elem1_2", "elem2_2"},
							notGoingToLog: "this should not be logged",
						},
					},
				},
			},
			want: map[string]any{
				"key1": "value1",
				"key2": 123,
				"key3": []any{"1", "2"},
				"key4": map[string]any{
					"sub_key1": "sub value1",
					"sub_key2": 3455,
				},
				"key5": map[string]any{
					"attribute_one": "struct attribute",
					"attr3":         []any{"elem1", "elem2"},
				},
				"key6": []any{
					map[string]any{
						"attribute_one": "struct attribute array 1",
						"attr3":         []any{"elem1_1", "elem2_1"},
					},
					map[string]any{
						"attribute_one": "struct attribute array 2",
						"attr3":         []any{"elem1_2", "elem2_2"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Invalid nested map key kind",
			args: args{
				input: map[string]any{
					"key1": "value1",
					"key2": map[int]any{
						1: "should be an error",
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Invalid input type",
			args: args{
				input: struct{ Field string }{"test"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Non-string keys in map",
			args: args{
				input: map[int]any{1: "value1"},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Nested structure with unsupported type",
			args: args{
				input: map[string]any{"key1": time.Date(2025, time.May, 20, 12, 33, 15, 8, time.UTC)},
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "Circular reference detection",
			args: args{
				input: map[string]any{"key1": nil},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := reflect.ValueOf(tt.args.input)
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}

			got, err := encodeMap(val)
			if (err != nil) != tt.wantErr {
				t.Errorf("encodeMap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_encodeArray(t *testing.T) {
	type args struct {
		input any
	}
	tests := []struct {
		name    string
		args    args
		want    []any
		wantErr bool
	}{
		{
			name: "Succes, raw array",
			args: args{
				input: []int{8, 6, 2, 3, 4},
			},
			want:    []any{8, 6, 2, 3, 4},
			wantErr: false,
		},
		{
			name: "Succes, raw array double",
			args: args{
				input: [][]int{{8, 6}, {2, 3}},
			},
			want:    []any{[]any{8, 6}, []any{2, 3}},
			wantErr: false,
		},
		{
			name: "Succes, map array",
			args: args{
				input: []map[string]any{
					{
						"key1": "something one",
						"key2": 2,
						"key4": true,
					},
					{
						"key1": "something two",
						"key2": 7,
						"key4": true,
					},
					{
						"key1": "something three",
						"key2": 100,
						"key4": false,
					},
				},
			},
			want: []any{
				map[string]any{
					"key1": "something one",
					"key2": 2,
					"key4": true,
				},
				map[string]any{
					"key1": "something two",
					"key2": 7,
					"key4": true,
				},
				map[string]any{
					"key1": "something three",
					"key2": 100,
					"key4": false,
				},
			},
			wantErr: false,
		},
		{
			name: "Succes, struct array",
			args: args{
				input: []testAttributes{
					{
						Attr1:         "struct attribute one",
						Attr2:         99887,
						Attr3:         []string{"elem1", "elem2"},
						notGoingToLog: "this should not be logged",
					},
					{
						Attr1:         "struct attribute two",
						Attr2:         1000,
						Attr3:         []string{"elem1_", "elem2"},
						notGoingToLog: "this should not be logged",
					},
					{
						Attr1:         "struct attribute three",
						Attr2:         5,
						Attr3:         []string{"elem1", "elem2__"},
						notGoingToLog: "this should not be logged",
					},
				},
			},
			want: []any{
				map[string]any{
					"attribute_one": "struct attribute one",
					"attr3":         []any{"elem1", "elem2"},
				},
				map[string]any{
					"attribute_one": "struct attribute two",
					"attr3":         []any{"elem1_", "elem2"},
				},
				map[string]any{
					"attribute_one": "struct attribute three",
					"attr3":         []any{"elem1", "elem2__"},
				},
			},
			wantErr: false,
		},
		{
			name: "Fail: not an array",
			args: args{
				input: "not an array",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := reflect.ValueOf(tt.args.input)
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}

			got, err := encodeArray(val)
			if (err != nil) != tt.wantErr {
				t.Errorf("encodeArray() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_encodeStruct(t *testing.T) {
	type testStruct struct {
		SomeMap    map[string]any `logit:"some_map"`
		SomeArray  []int
		SomeStruct testAttributes
	}
	type args struct {
		input any
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]any
		wantErr bool
	}{
		{
			name: "Success",
			args: args{
				input: testStruct{
					SomeMap: map[string]any{
						"some_map_key":   1,
						"some_map_key_2": "some map string",
					},
					SomeArray: []int{9, 4, 1},
					SomeStruct: testAttributes{
						Attr1:         "struct attribute one",
						Attr2:         99887,
						Attr3:         []string{"elem1", "elem2"},
						notGoingToLog: "this should not be logged",
					},
				},
			},
			want: map[string]any{
				"some_map": map[string]any{
					"some_map_key":   1,
					"some_map_key_2": "some map string",
				},
				"somearray": []any{9, 4, 1},
				"somestruct": map[string]any{
					"attribute_one": "struct attribute one",
					"attr3":         []any{"elem1", "elem2"},
				},
			},
			wantErr: false,
		},
		{
			name: "Fail: Not a struct",
			args: args{
				input: 1,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			val := reflect.ValueOf(tt.args.input)
			if val.Kind() == reflect.Ptr {
				val = val.Elem()
			}

			got, err := encodeStruct(val)
			if (err != nil) != tt.wantErr {
				t.Errorf("encodeStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_encode(t *testing.T) {
	type testStruct struct {
		SomeMap    map[string]any `logit:"some_map"`
		SomeArray  []int
		SomeStruct testAttributes
	}
	type args struct {
		input any
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr bool
	}{
		{
			name: "Success: Struct",
			args: args{
				input: testStruct{
					SomeMap: map[string]any{
						"some_map_key":   1,
						"some_map_key_2": "some map string",
					},
					SomeArray: []int{9, 4, 1},
					SomeStruct: testAttributes{
						Attr1:         "struct attribute one",
						Attr2:         99887,
						Attr3:         []string{"elem1", "elem2"},
						notGoingToLog: "this should not be logged",
					},
				},
			},
			want: map[string]any{
				"some_map": map[string]any{
					"some_map_key":   1,
					"some_map_key_2": "some map string",
				},
				"somearray": []any{9, 4, 1},
				"somestruct": map[string]any{
					"attribute_one": "struct attribute one",
					"attr3":         []any{"elem1", "elem2"},
				},
			},
			wantErr: false,
		},
		{
			name: "Succes: Map",
			args: args{
				input: map[string]any{
					"key1": "value1",
					"key2": 123,
					"key3": []any{"1", "2"},
					"key4": map[string]any{
						"sub_key1": "sub value1",
						"sub_key2": 3455,
					},
					"key5": testAttributes{
						Attr1:         "struct attribute",
						Attr2:         99887,
						Attr3:         []string{"elem1", "elem2"},
						notGoingToLog: "this should not be logged",
					},
					"key6": []testAttributes{
						{
							Attr1:         "struct attribute array 1",
							Attr2:         9988455557,
							Attr3:         []string{"elem1_1", "elem2_1"},
							notGoingToLog: "this should not be logged",
						},
						{
							Attr1:         "struct attribute array 2",
							Attr2:         998444487,
							Attr3:         []string{"elem1_2", "elem2_2"},
							notGoingToLog: "this should not be logged",
						},
					},
				},
			},
			want: map[string]any{
				"key1": "value1",
				"key2": 123,
				"key3": []any{"1", "2"},
				"key4": map[string]any{
					"sub_key1": "sub value1",
					"sub_key2": 3455,
				},
				"key5": map[string]any{
					"attribute_one": "struct attribute",
					"attr3":         []any{"elem1", "elem2"},
				},
				"key6": []any{
					map[string]any{
						"attribute_one": "struct attribute array 1",
						"attr3":         []any{"elem1_1", "elem2_1"},
					},
					map[string]any{
						"attribute_one": "struct attribute array 2",
						"attr3":         []any{"elem1_2", "elem2_2"},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Fail",
			args: args{
				input: []int{1, 2, 3},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := encode(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
