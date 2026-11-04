package converter

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
)

type Item struct {
	Key, Value string
}

// CDBItems flattens m into items keyed by "/"-delimited paths below path ("" for the root).
//
// If childLists is set, every non-empty map also gets an onlineconf child list, the way the
// onlineconf server adds them for a module with child_lists enabled: the key is the map's
// path followed by "/" (just "/" for the root), the value is a JSON array of the map's keys
// sorted bytewise. Keys must then be valid onlineconf names, non-empty and free of "/",
// otherwise they would collide with or contradict the lists.
func CDBItems(m map[string]any, path string, childLists bool) ([]Item, error) {
	items := make([]Item, 0)

	if childLists && len(m) > 0 {
		item, err := childList(m, path)
		if err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	for k, v := range m {
		p := strings.Join([]string{path, k}, "/")

		switch reflect.ValueOf(v).Kind() {
		case reflect.Map:
			res, err := CDBItems(v.(map[string]any), p, childLists)
			if err != nil {
				return nil, err
			}

			items = append(items, res...)
		case reflect.Int:
			items = append(items, Item{Key: p, Value: fmt.Sprintf("s%d", v)})
		case reflect.Float64:
			strValue := strings.TrimRight(fmt.Sprintf("s%f", v), "0")
			items = append(items, Item{Key: p, Value: strValue})
		case reflect.Slice:
			j, err := json.Marshal(v)
			if err != nil {
				return nil, err
			}
			items = append(items, Item{Key: p, Value: fmt.Sprintf("j%s", j)})
		default:
			var m map[string]any
			vStr, _ := v.(string)

			err := json.Unmarshal([]byte(vStr), &m)
			if err != nil {
				items = append(items, Item{Key: p, Value: "s" + vStr})
				continue
			}
			items = append(items, Item{Key: p, Value: "j" + vStr})
		}
	}

	return items, nil
}

// childList builds the child list item of the non-empty map m at path, matching
// deepCreateChildList in the onlineconf admin (admin/go/resolver/tree.go).
func childList(m map[string]any, path string) (Item, error) {
	names := make([]string, 0, len(m))
	for name := range m {
		if name == "" || strings.Contains(name, "/") {
			return Item{}, fmt.Errorf("key %q in %s/: child lists need onlineconf names, which are non-empty and contain no \"/\"", name, path)
		}
		names = append(names, name)
	}
	sort.Strings(names)

	list, err := json.Marshal(names)
	if err != nil {
		return Item{}, err
	}

	return Item{Key: path + "/", Value: "j" + string(list)}, nil
}
