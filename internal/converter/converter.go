package converter

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

type Item struct {
	Key, Value string
}

func CDBItems(m map[string]any, path string) ([]Item, error) {
	items := make([]Item, 0)
	for k, v := range m {
		p := strings.Join([]string{path, k}, "/")

		switch reflect.ValueOf(v).Kind() {
		case reflect.Map:
			res, err := CDBItems(v.(map[string]any), p)
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
