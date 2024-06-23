package validators

import (
	"testing"
)

func Test_NotZero_String(t *testing.T) {
	v := "test"
	err := NotZero(v)
	if err != nil {
		t.Error("NotZero failed")
	}

	v = ""
	err = NotZero(v)
	if err == nil {
		t.Error("NotZero failed")
	}
}

func Test_MinLength(t *testing.T) {
	// String
	v := "test"
	err := MinLength[string](1)(v)
	if err != nil {
		t.Error("MinLength failed")
	}

	// Zero
	v = ""
	err = MinLength[string](1)(v)
	if err == nil {
		t.Error("MinLength failed")
	}

	// Array Strings
	vs := []string{"test"}
	err = MinLength[[]string](1)(vs)
	if err != nil {
		t.Error("MinLength failed")
	}

	// Zero
	vs = []string{}
	err = MinLength[[]string](1)(vs)
	if err == nil {
		t.Error("MinLength failed")
	}
}

func Test_MaxLength(t *testing.T) {
	// String
	v := "test"
	err := MaxLength[string](4)(v)
	if err != nil {
		t.Error("MaxLength failed" + err.Error())
	}

	// Zero
	v = ""
	err = MaxLength[string](4)(v)
	if err != nil {
		t.Errorf("MaxLength('%s') failed", v)
	}

	// Array Strings
	vs := []string{"1", "2", "3", "4", "5"}
	err = MaxLength[[]string](4)(vs)
	if err == nil {
		t.Error("MaxLength failed")
	}

	// Zero
	vs = []string{}
	err = MaxLength[[]string](4)(vs)
	if err != nil {
		t.Error("MaxLength failed")
	}
}

func Test_Combine(t *testing.T) {
	combined := Combine(NotZero, MinLength[string](4))

	v := "test"
	err := combined(v)
	if err != nil {
		t.Error("Combine failed")
	}

	v = "t"
	err = combined(v)
	if err == nil {
		t.Error("Combine failed")
	}
}
