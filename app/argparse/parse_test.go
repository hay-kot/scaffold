package argparse

import (
	"reflect"
	"testing"
)

func Test_Parse_StringValues(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    map[string]any
		wantErr bool
	}{
		{
			name: "simple_string_without_type",
			args: []string{"name=MyProject"},
			want: map[string]any{
				"name": "MyProject",
			},
		},
		{
			name: "multiple_strings_without_type",
			args: []string{"name=MyProject", "description=My awesome project"},
			want: map[string]any{
				"name":        "MyProject",
				"description": "My awesome project",
			},
		},
		{
			name: "explicit_string_type",
			args: []string{"name:string=MyProject"},
			want: map[string]any{
				"name": "MyProject",
			},
		},
		{
			name: "string_shorthand",
			args: []string{"title:str=Hello World"},
			want: map[string]any{
				"title": "Hello World",
			},
		},
		{
			name: "empty_string",
			args: []string{"name:string="},
			want: map[string]any{
				"name": "",
			},
		},
		{
			name: "string_with_equals_sign",
			args: []string{"formula=a=b+c"},
			want: map[string]any{
				"formula": "a=b+c",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Parse_IntegerValues(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    map[string]any
		wantErr bool
	}{
		{
			name: "positive_integer",
			args: []string{"count:int=5"},
			want: map[string]any{
				"count": 5,
			},
		},
		{
			name: "negative_integer",
			args: []string{"offset:int=-10"},
			want: map[string]any{
				"offset": -10,
			},
		},
		{
			name: "zero",
			args: []string{"start:int=0"},
			want: map[string]any{
				"start": 0,
			},
		},
		{
			name: "large_integer",
			args: []string{"port:int=65535"},
			want: map[string]any{
				"port": 65535,
			},
		},
		{
			name:    "invalid_integer",
			args:    []string{"count:int=abc"},
			wantErr: true,
		},
		{
			name:    "float_as_integer",
			args:    []string{"count:int=3.14"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Parse_FloatValues(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    map[string]any
		wantErr bool
	}{
		{
			name: "simple_float",
			args: []string{"price:float=19.99"},
			want: map[string]any{
				"price": 19.99,
			},
		},
		{
			name: "float64_explicit",
			args: []string{"rate:float64=0.5"},
			want: map[string]any{
				"rate": 0.5,
			},
		},
		{
			name: "float32_explicit",
			args: []string{"value:float32=3.14"},
			want: map[string]any{
				"value": float32(3.14),
			},
		},
		{
			name: "negative_float",
			args: []string{"temp:float=-273.15"},
			want: map[string]any{
				"temp": -273.15,
			},
		},
		{
			name: "integer_as_float",
			args: []string{"whole:float=42"},
			want: map[string]any{
				"whole": 42.0,
			},
		},
		{
			name:    "invalid_float",
			args:    []string{"price:float=not-a-number"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Parse_BooleanValues(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    map[string]any
		wantErr bool
	}{
		{
			name: "true_value",
			args: []string{"enabled:bool=true"},
			want: map[string]any{
				"enabled": true,
			},
		},
		{
			name: "false_value",
			args: []string{"disabled:bool=false"},
			want: map[string]any{
				"disabled": false,
			},
		},
		{
			name: "numeric_true",
			args: []string{"yes:bool=1"},
			want: map[string]any{
				"yes": true,
			},
		},
		{
			name: "numeric_false",
			args: []string{"no:bool=0"},
			want: map[string]any{
				"no": false,
			},
		},
		{
			name: "uppercase_true",
			args: []string{"flag:bool=TRUE"},
			want: map[string]any{
				"flag": true,
			},
		},
		{
			name:    "invalid_boolean",
			args:    []string{"maybe:bool=maybe"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Parse_StringSliceValues(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    map[string]any
		wantErr bool
	}{
		{
			name: "simple_string_slice",
			args: []string{"langs:[]string=Python,JavaScript,Go"},
			want: map[string]any{
				"langs": []string{"Python", "JavaScript", "Go"},
			},
		},
		{
			name: "string_slice_shorthand",
			args: []string{"tags:[]str=dev,test,prod"},
			want: map[string]any{
				"tags": []string{"dev", "test", "prod"},
			},
		},
		{
			name: "single_element",
			args: []string{"envs:[]string=production"},
			want: map[string]any{
				"envs": []string{"production"},
			},
		},
		{
			name: "empty_slice",
			args: []string{"items:[]string="},
			want: map[string]any{
				"items": []string{},
			},
		},
		{
			name: "slice_with_spaces",
			args: []string{"phrases:[]string=hello world,foo bar"},
			want: map[string]any{
				"phrases": []string{"hello world", "foo bar"},
			},
		},
		{
			name: "escaped_commas",
			args: []string{`items:[]string=hello\,world,foo\,bar,normal`},
			want: map[string]any{
				"items": []string{"hello,world", "foo,bar", "normal"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Parse_IntSliceValues(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    map[string]any
		wantErr bool
	}{
		{
			name: "simple_int_slice",
			args: []string{"ports:[]int=8080,8081,8082"},
			want: map[string]any{
				"ports": []int{8080, 8081, 8082},
			},
		},
		{
			name: "negative_numbers",
			args: []string{"temps:[]int=-10,0,10,20"},
			want: map[string]any{
				"temps": []int{-10, 0, 10, 20},
			},
		},
		{
			name: "single_int",
			args: []string{"ids:[]int=42"},
			want: map[string]any{
				"ids": []int{42},
			},
		},
		{
			name: "empty_int_slice",
			args: []string{"nums:[]int="},
			want: map[string]any{
				"nums": []int{},
			},
		},
		{
			name:    "invalid_int_in_slice",
			args:    []string{"ports:[]int=8080,not-a-number,8082"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Parse_FloatSliceValues(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    map[string]any
		wantErr bool
	}{
		{
			name: "simple_float_slice",
			args: []string{"rates:[]float=0.1,0.5,0.9"},
			want: map[string]any{
				"rates": []float64{0.1, 0.5, 0.9},
			},
		},
		{
			name: "float64_slice",
			args: []string{"values:[]float64=1.1,2.2,3.3"},
			want: map[string]any{
				"values": []float64{1.1, 2.2, 3.3},
			},
		},
		{
			name: "float32_slice",
			args: []string{"scores:[]float32=0.5,1.5,2.5"},
			want: map[string]any{
				"scores": []float32{0.5, 1.5, 2.5},
			},
		},
		{
			name: "mixed_int_float",
			args: []string{"mixed:[]float=1,2.5,3"},
			want: map[string]any{
				"mixed": []float64{1.0, 2.5, 3.0},
			},
		},
		{
			name:    "invalid_float_in_slice",
			args:    []string{"rates:[]float=0.1,invalid,0.9"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Parse_BoolSliceValues(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    map[string]any
		wantErr bool
	}{
		{
			name: "simple_bool_slice",
			args: []string{"flags:[]bool=true,false,true"},
			want: map[string]any{
				"flags": []bool{true, false, true},
			},
		},
		{
			name: "numeric_bools",
			args: []string{"bits:[]bool=1,0,1,0"},
			want: map[string]any{
				"bits": []bool{true, false, true, false},
			},
		},
		{
			name: "mixed_formats",
			args: []string{"mixed:[]bool=true,1,false,0"},
			want: map[string]any{
				"mixed": []bool{true, true, false, false},
			},
		},
		{
			name:    "invalid_bool_in_slice",
			args:    []string{"flags:[]bool=true,maybe,false"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Parse_JSONValues(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    map[string]any
		wantErr bool
	}{
		{
			name: "json_object",
			args: []string{`config:json={"debug":true,"port":8080,"name":"app"}`},
			want: map[string]any{
				"config": map[string]interface{}{
					"debug": true,
					"port":  float64(8080),
					"name":  "app",
				},
			},
		},
		{
			name: "json_array",
			args: []string{`data:json=[1,2,3,4,5]`},
			want: map[string]any{
				"data": []interface{}{float64(1), float64(2), float64(3), float64(4), float64(5)},
			},
		},
		{
			name: "json_string",
			args: []string{`message:json="Hello, World!"`},
			want: map[string]any{
				"message": "Hello, World!",
			},
		},
		{
			name: "json_number",
			args: []string{`count:json=42`},
			want: map[string]any{
				"count": float64(42),
			},
		},
		{
			name: "json_boolean",
			args: []string{`active:json=true`},
			want: map[string]any{
				"active": true,
			},
		},
		{
			name: "json_null",
			args: []string{`value:json=null`},
			want: map[string]any{
				"value": nil,
			},
		},
		{
			name: "nested_json",
			args: []string{`user:json={"name":"John","age":30,"tags":["admin","user"]}`},
			want: map[string]any{
				"user": map[string]interface{}{
					"name": "John",
					"age":  float64(30),
					"tags": []interface{}{"admin", "user"},
				},
			},
		},
		{
			name: "json_with_equals",
			args: []string{`query:json={"filter":"status=active"}`},
			want: map[string]any{
				"query": map[string]interface{}{
					"filter": "status=active",
				},
			},
		},
		{
			name:    "invalid_json",
			args:    []string{`config:json={invalid json}`},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Parse_MixedTypes(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    map[string]any
		wantErr bool
	}{
		{
			name: "all_types_together",
			args: []string{
				"name=MyProject",
				"port:int=8080",
				"debug:bool=true",
				"version:float=1.5",
				"langs:[]string=Go,Python",
				"ports:[]int=8080,8081",
				"rates:[]float=0.1,0.2",
				"flags:[]bool=true,false",
				`config:json={"theme":"dark"}`,
			},
			want: map[string]any{
				"name":    "MyProject",
				"port":    8080,
				"debug":   true,
				"version": 1.5,
				"langs":   []string{"Go", "Python"},
				"ports":   []int{8080, 8081},
				"rates":   []float64{0.1, 0.2},
				"flags":   []bool{true, false},
				"config":  map[string]interface{}{"theme": "dark"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Parse_EdgeCases(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		want    map[string]any
		wantErr bool
	}{
		{
			name:    "missing_equals",
			args:    []string{"name"},
			wantErr: true,
		},
		{
			name:    "empty_key",
			args:    []string{"=value"},
			wantErr: true,
		},
		{
			name:    "unknown_type",
			args:    []string{"value:unknown=something"},
			wantErr: true,
		},
		{
			name: "no_arguments",
			args: []string{},
			want: map[string]any{},
		},
		{
			name:    "invalid_slice_type",
			args:    []string{"data:[]unknown=a,b,c"},
			wantErr: true,
		},
		{
			name:    "empty_type_hint",
			args:    []string{"key:=value"},
			wantErr: true,
		},
		{
			name:    "empty_key_with_type",
			args:    []string{":string=value"},
			wantErr: true,
		},
		{
			name: "multiple_equals_in_value",
			args: []string{"equation=a=b=c=d"},
			want: map[string]any{
				"equation": "a=b=c=d",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_splitEscaped(t *testing.T) {
	tests := []struct {
		name  string
		input string
		sep   rune
		want  []string
	}{
		{
			name:  "no_escapes",
			input: "a,b,c",
			sep:   ',',
			want:  []string{"a", "b", "c"},
		},
		{
			name:  "escaped_separator",
			input: `a\,b,c\,d`,
			sep:   ',',
			want:  []string{"a,b", "c,d"},
		},
		{
			name:  "mixed_escapes",
			input: `normal,has\,comma,another`,
			sep:   ',',
			want:  []string{"normal", "has,comma", "another"},
		},
		{
			name:  "trailing_escape",
			input: `a,b\,`,
			sep:   ',',
			want:  []string{"a", "b,"},
		},
		{
			name:  "empty_elements",
			input: "a,,b",
			sep:   ',',
			want:  []string{"a", "", "b"},
		},
		{
			name:  "only_separator",
			input: ",",
			sep:   ',',
			want:  []string{"", ""},
		},
		{
			name:  "empty_string",
			input: "",
			sep:   ',',
			want:  []string{""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitEscaped(tt.input, tt.sep)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitEscaped() = %v, want %v", got, tt.want)
			}
		})
	}
}
