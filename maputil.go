package jules

import (
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	aryRegex = regexp.MustCompile(`^([^.\s\[\]]*)(\[(\d*|[-])\]|$)`)
)

func lastParentAndChildToken(key string, data map[string]interface{}) (interface{}, string) {
	var parent interface{}
	var parentPath string
	var lastToken string

	parent = data

	lastSep := strings.LastIndexAny(key, ".[")

	if lastSep < 0 {
		return parent, key
	}

	parentPath = key[0:lastSep]
	lastToken = key[lastSep+1:]
	if strings.HasSuffix(lastToken, "]") {
		lastToken = lastToken[:len(lastToken)-1]
	}

	pathParts := strings.Split(parentPath, ".")
Out:
	for _, k := range pathParts {
		key, aryIndex := getKeyAndIndex(k)

		switch reflect.TypeOf(parent).Kind() {
		case reflect.Map:
			if obj, found := parent.(map[string]interface{})[key]; found {

				if aryIndex > -1 && reflect.TypeOf(obj).Kind().String() == reflect.Slice.String() {
					parent = obj.([]interface{})[aryIndex]
				} else {
					parent = obj
				}
			} else {
				parent = nil
				break Out
			}
		case reflect.Slice:
			parent = parent.([]interface{})[aryIndex]
		}
	}

	return parent, lastToken
}

func getFromMapByDotPath(key string, data map[string]interface{}) (interface{}, bool) {
	var val interface{}

	parent, childToken := lastParentAndChildToken(key, data)

	if parent != nil {
		switch reflect.TypeOf(parent).Kind() {
		case reflect.Map:
			val = parent.(map[string]interface{})[childToken]
		case reflect.Slice:
			ps := parent.([]interface{})
			idx, _ := strconv.Atoi(childToken)

			if idx < len(ps) {
				val = parent.([]interface{})[idx]
			}
		}
	}

	return val, (val != nil)
}

func addToMapByDotPath(key string, data map[string]interface{}, newval interface{}) bool {
	var added bool

	parent, childToken := lastParentAndChildToken(key, data)

	if parent != nil {
		switch reflect.TypeOf(parent).Kind() {
		case reflect.Map:
			parent.(map[string]interface{})[childToken] = newval
			added = true
		case reflect.Slice:
			ps := parent.([]interface{})
			i, _ := strconv.Atoi(childToken)

			if i == math.MaxInt64 {
				ps = append(ps, newval)
				added = true
			} else if i < len(ps) {
				ps = append(ps, 0)
				copy(ps[i+1:], ps[i:])
				ps[i] = newval
				added = true
			}

			if added {
				lastSep := strings.LastIndex(key, "[")
				parentPath := key[0:lastSep]
				addToMapByDotPath(parentPath, data, ps)
			}
		}
	}

	return added
}

func replaceInMapByDotPath(key string, data map[string]interface{}, newval interface{}) bool {
	var replaced bool

	parent, childToken := lastParentAndChildToken(key, data)

	if parent != nil {
		switch reflect.TypeOf(parent).Kind() {
		case reflect.Map:
			pmap := parent.(map[string]interface{})

			if _, exists := pmap[childToken]; exists {
				pmap[childToken] = newval
				replaced = true
			}

		case reflect.Slice:
			ps := parent.([]interface{})
			i, _ := strconv.Atoi(childToken)

			if i < len(ps) {
				ps[i] = newval
				replaced = true

				lastSep := strings.LastIndex(key, "[")
				parentPath := key[0:lastSep]
				addToMapByDotPath(parentPath, data, ps)
			}
		}
	}

	return replaced
}

func moveInMapByDotPath(from, to string, data map[string]interface{}) bool {
	var moved bool

	if val, found := getFromMapByDotPath(from, data); found {
		deleteFromMapByDotPath(from, data)
		moved = addToMapByDotPath(to, data, val)
	}

	return moved
}

func copyInMapByDotPath(from, to string, data map[string]interface{}) bool {
	var copied bool

	if val, found := getFromMapByDotPath(from, data); found {
		copied = addToMapByDotPath(to, data, val)
	}

	return copied
}

func deleteFromMapByDotPath(key string, data map[string]interface{}) bool {
	var baleeted bool

	parent, childToken := lastParentAndChildToken(key, data)

	if parent != nil {
		switch reflect.TypeOf(parent).Kind() {
		case reflect.Map:
			delete(parent.(map[string]interface{}), childToken)
			baleeted = true
		case reflect.Slice:
			ps := parent.([]interface{})
			i, _ := strconv.Atoi(childToken)
			ps = append(ps[:i], ps[i+1:]...)

			lastSep := strings.LastIndex(key, "[")
			parentPath := key[0:lastSep]
			addToMapByDotPath(parentPath, data, ps)

			baleeted = true
		}
	}

	return baleeted
}

func getKeyAndIndex(varname string) (string, int) {
	parts := aryRegex.FindStringSubmatch(varname)
	key := varname
	idx := -1

	if len(parts) == 4 {
		key = parts[1]

		if len(parts[3]) > 0 {
			if parts[3] == "-" {
				idx = math.MaxInt64
			} else {
				idx, _ = strconv.Atoi(parts[3])
			}
		}
	}

	return key, idx
}
