package shopify

import (
	"reflect"
	"testing"
)

func TestParseScopes(t *testing.T) {
	testCases := []struct {
		S        string
		Expected Scopes
	}{
		{
			"",
			Scopes{},
		},
		{
			",",
			Scopes{},
		},
		{
			"a",
			Scopes{"a"},
		},
		{
			"a,b,,c",
			Scopes{"a", "b", "c"},
		},
		{
			" ,, ,a,  ,,b,,c, ,, ",
			Scopes{"a", "b", "c"},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.S, func(t *testing.T) {
			value, err := ParseScopes(testCase.S)

			if err != nil {
				t.Fatalf("expected no error but got: %s", err)
			}

			if !reflect.DeepEqual(testCase.Expected, value) {
				t.Errorf("expected: %#v\ngot: %#v", testCase.Expected, value)
			}
		})
	}
}
