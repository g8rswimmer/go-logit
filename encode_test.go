package logit

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

//	func TestEncode(t *testing.T) {
//		type args struct {
//			input any
//		}
//		tests := []struct {
//			name    string
//			args    args
//			want    any
//			wantErr bool
//		}{}
//		for _, tt := range tests {
//			t.Run(tt.name, func(t *testing.T) {
//				got, err := Encode(tt.args.input)
//				if (err != nil) != tt.wantErr {
//					t.Errorf("Encode() error = %v, wantErr %v", err, tt.wantErr)
//					return
//				}
//				assert.Equal(t, tt.want, got)
//			})
//		}
//	}
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
