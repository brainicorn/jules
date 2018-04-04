package jules

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	aryRegex = regexp.MustCompile(`^([^.\s\[\]]*)(\[(\d*)\]|$)`)
)

func valueFromMapByDotPath(key string, data map[string]interface{}) (interface{}, bool) {
	pathParts := strings.Split(key, ".")
	parent := data
	for i, k := range pathParts {
		key, aryIndex := getKeyAndIndex(k)
		if val, ok := parent[key]; ok {
			if i == len(pathParts)-1 && aryIndex < 0 {
				return val, true
			}

			switch val.(type) {
			case map[string]interface{}:
				parent = val.(map[string]interface{})
			case []map[string]interface{}:
				parent = val.([]map[string]interface{})[aryIndex]
			default:
				if reflect.ValueOf(val).Kind().String() == reflect.Slice.String() {
					return reflect.ValueOf(val).Index(aryIndex).Interface(), true
				}
			}
		} else {
			break
		}
	}

	return nil, false
}

func deleteFromMapByDotPath(key string, data map[string]interface{}) bool {
	pathParts := strings.Split(key, ".")
	parent := data
	for i, k := range pathParts {
		key, aryIndex := getKeyAndIndex(k)
		if val, ok := parent[key]; ok {
			if i == len(pathParts)-1 && aryIndex < 0 {
				delete(parent, key)
				return true
			}

			switch val.(type) {
			case map[string]interface{}:
				parent = val.(map[string]interface{})
			case []map[string]interface{}:
				parent = val.([]map[string]interface{})[aryIndex]
			default:
				break
			}
		} else {
			break
		}
	}

	return false
}

func getKeyAndIndex(varname string) (string, int) {
	parts := aryRegex.FindStringSubmatch(varname)
	key := varname
	idx := -1

	if len(parts) == 4 {
		key = parts[1]

		if len(parts[3]) > 0 {
			idx, _ = strconv.Atoi(parts[3])
		}
	}

	return key, idx
}
