package converter

import (
	"reflect"
	"testing"
)

func TestCDBItems(t *testing.T) {
	testCases := []struct {
		name     string
		inputMap map[string]any
		expected []Item
	}{
		{
			name: "int",
			inputMap: map[string]any{
				"int": 123,
			},
			expected: []Item{
				{Key: "/int", Value: "s123"},
			},
		},
		{
			name: "float",
			inputMap: map[string]any{
				"float": 1.25,
			},
			expected: []Item{
				{Key: "/float", Value: "s1.25"},
			},
		},
		{
			name: "string",
			inputMap: map[string]any{
				"string": "localhost",
			},
			expected: []Item{
				{Key: "/string", Value: "slocalhost"},
			},
		},
		{
			name: "unicode",
			inputMap: map[string]any{
				"string": "привет",
			},
			expected: []Item{
				{Key: "/string", Value: "sпривет"},
			},
		},
		{
			name: "json slice",
			inputMap: map[string]any{
				"list": []any{1, "string"},
			},
			expected: []Item{
				{Key: "/list", Value: `j[1,"string"]`},
			},
		},
		{
			name: "json string",
			inputMap: map[string]any{
				"json": `{"key": "value"}`,
			},
			expected: []Item{
				{Key: "/json", Value: `j{"key": "value"}`},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := CDBItems(tc.inputMap, "")
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}
