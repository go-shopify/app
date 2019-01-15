package shopify

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestParseScope(t *testing.T) {
	testCases := []struct {
		S        string
		Expected Scope
	}{
		{
			"",
			Scope{},
		},
		{
			",",
			Scope{},
		},
		{
			"a",
			Scope{"a"},
		},
		{
			"a,b,,c",
			Scope{"a", "b", "c"},
		},
		{
			" ,, ,a,  ,,b,,c, ,, ",
			Scope{"a", "b", "c"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.S, func(t *testing.T) {
			value, err := ParseScope(testCase.S)

			if err != nil {
				t.Fatalf("expected no error but got: %s", err)
			}

			if !reflect.DeepEqual(testCase.Expected, value) {
				t.Errorf("expected: %#v\ngot: %#v", testCase.Expected, value)
			}
		})
	}
}

func TestScopeJSON(t *testing.T) {
	value := struct {
		X Scope `json:"x"`
	}{}

	data := []byte(`{"x": "a,b"}`)

	if err := json.Unmarshal(data, &value); err != nil {
		t.Fatalf("expected no error but got: %s", err)
	}

	ref := Scope{"a", "b"}

	if !reflect.DeepEqual(value.X, ref) {
		t.Errorf("values differ (%#v != %#v)", value.X, ref)
	}

	data, err := json.Marshal(value)

	if err != nil {
		t.Fatalf("expected no error but got: %s", err)
	}

	expected := `{"x":"a,b"}`

	if string(data) != expected {
		t.Errorf("expected: %s\ngot: %s", expected, string(data))
	}
}
