package validators

import "reflect"

type Validatable interface {
	string | []string
}

type validatehandler struct {
	strfn   func(string) error
	slicefn func([]string) error
}

func (vh *validatehandler) validate(v any) error {
	// Get the reflect.Value of the interface
	val := reflect.ValueOf(v)

	// Check if the value is a slice, array, string, map, or channel
	switch val.Kind() {
	case reflect.String:
		return vh.strfn(val.String())
	case reflect.Slice:
		// cast to string slice
		strSlice := any(v).([]string)
		return vh.slicefn(strSlice)
	default:
		panic("unsupported type")
	}
}
