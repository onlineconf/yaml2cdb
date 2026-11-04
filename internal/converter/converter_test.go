package converter

import (
	"reflect"
	"sort"
	"testing"
)

func TestCDBItems(t *testing.T) {
	testCases := []struct {
		name       string
		inputMap   map[string]any
		childLists bool
		expected   []Item
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
		{
			name: "nested map",
			inputMap: map[string]any{
				"db": map[string]any{"host": "localhost", "port": 5432},
			},
			expected: []Item{
				{Key: "/db/host", Value: "slocalhost"},
				{Key: "/db/port", Value: "s5432"},
			},
		},
		{
			name: "child lists",
			inputMap: map[string]any{
				"b": map[string]any{
					"y": "1",
					"x": map[string]any{"z": 2},
				},
				"a":     "s",
				"empty": map[string]any{},
				"list":  []any{"q"},
			},
			childLists: true,
			expected: []Item{
				{Key: "/", Value: `j["a","b","empty","list"]`},
				{Key: "/a", Value: "ss"},
				{Key: "/b/", Value: `j["x","y"]`},
				{Key: "/b/x/", Value: `j["z"]`},
				{Key: "/b/x/z", Value: "s2"},
				{Key: "/b/y", Value: "s1"},
				{Key: "/list", Value: `j["q"]`},
			},
		},
		{
			// bytewise order and HTML-safe escaping, as json.Marshal in the onlineconf admin
			name: "child list order and escaping",
			inputMap: map[string]any{
				"é":   1,
				"a":   1,
				"B":   1,
				"<x>": 1,
			},
			childLists: true,
			expected: []Item{
				{Key: "/", Value: `j["\u003cx\u003e","B","a","é"]`},
				{Key: "/<x>", Value: "s1"},
				{Key: "/B", Value: "s1"},
				{Key: "/a", Value: "s1"},
				{Key: "/é", Value: "s1"},
			},
		},
		{
			name:       "child lists of empty root",
			inputMap:   map[string]any{},
			childLists: true,
			expected:   []Item{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := CDBItems(tc.inputMap, "", tc.childLists)
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}

			sort.Slice(result, func(i, j int) bool { return result[i].Key < result[j].Key })
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("Expected %v, but got %v", tc.expected, result)
			}
		})
	}
}

func TestCDBItemsChildListsInvalidKeys(t *testing.T) {
	testCases := []struct {
		name     string
		inputMap map[string]any
	}{
		{name: "empty key", inputMap: map[string]any{"": "x"}},
		{name: "key with slash", inputMap: map[string]any{"a/b": "x"}},
		{name: "nested empty key", inputMap: map[string]any{"a": map[string]any{"": "x"}}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := CDBItems(tc.inputMap, "", true); err == nil {
				t.Error("Expected an error with child lists")
			}

			if _, err := CDBItems(tc.inputMap, "", false); err != nil {
				t.Errorf("Unexpected error without child lists: %v", err)
			}
		})
	}
}
