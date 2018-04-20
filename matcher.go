package jules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
)

type JulesMatcher struct {
	Rules JuleSet
}

func NewMatcher(jrules json.RawMessage) (*JulesMatcher, error) {
	var err error
	var juleset JuleSet

	juleset, err = validateJuleSet(jrules)

	if err == nil && (len(juleset) < 1 || numRootConditions(juleset) < 1) {
		err = fmt.Errorf("No rules found")
	}

	return &JulesMatcher{Rules: juleset}, err
}

func numRootConditions(juleset JuleSet) int {
	var num = 0

	for _, jule := range juleset {
		if hasConditions(jule.Condition) {
			num++
		}
	}

	return num
}

func hasConditions(c CompositeOrCondition) bool {
	var emptyComposite Composite
	var emptyCondition Condition

	switch c.(type) {
	case Condition:
		return !reflect.DeepEqual(c, emptyCondition)
	case Composite:
		return !reflect.DeepEqual(c, emptyComposite)
	}
	return false

}

func (m *JulesMatcher) Match(payload json.RawMessage) (bool, error) {
	return m.MatchAt(payload, "")
}

func (m *JulesMatcher) MatchAt(payload json.RawMessage, rootPath string) (bool, error) {
	var err error
	var allMatched bool
	var objects []map[string]interface{}

	if bytes.HasPrefix(payload, []byte("[{")) {
		// if it's an array of objects, we can just unmarshal it
		err = json.Unmarshal(payload, &objects)
	} else if bytes.HasPrefix(payload, []byte("{")) {
		// if it's a single object, we unmarshal it and add it to the objects slice manually
		// this is just so we can always use a slice of objects for consistency
		var single map[string]interface{}
		err = json.Unmarshal(payload, &single)
		if err == nil {
			objects = make([]map[string]interface{}, 1)
			objects[0] = single
		}
	} else {
		// if we got here, we don't know what the heck the payload is
		err = fmt.Errorf("payload must be a json object or a json array of objects")
	}

	if err == nil {
	MatchBreak:
		// loop over each object and grab our starting point
		for _, obj := range objects {
			roots, rootFound := findRoots(obj, rootPath)

			if !rootFound {
				err = fmt.Errorf("root object not found at path '%s'", rootPath)
			}

			if err == nil {
				rt := reflect.ValueOf(roots)
				switch rt.Kind() {
				case reflect.Slice:
					for _, root := range roots.([]interface{}) {

						rootMap, isMap := root.(map[string]interface{})
						if isMap {
							allMatched = applyMatchRules(rootMap, m.Rules)
							if !allMatched {
								break MatchBreak
							}
						} else {
							allMatched = false
							break MatchBreak
						}

					}

				case reflect.Map:
					allMatched = applyMatchRules(roots.(map[string]interface{}), m.Rules)

				default:
					allMatched = false
					break MatchBreak
				}
			}
		}

	}

	return allMatched, err
}

func applyMatchRules(root map[string]interface{}, rules JuleSet) bool {
	allMatched := true

	for _, rule := range rules {
		matched, _ := matchCompositeOrCondition(root, rule.Condition)
		if !matched {
			allMatched = false
		}
	}

	return allMatched
}

func matchCompositeOrCondition(root map[string]interface{}, corc CompositeOrCondition) (bool, []string) {
	var matched bool
	var fails []string

	switch corc.(type) {
	case Composite:
		matched, fails = matchComposite(root, corc.(Composite))

	case Condition:
		matched, fails = matchCondition(root, corc.(Condition))
	}

	return matched, fails
}

func matchComposite(root map[string]interface{}, c Composite) (bool, []string) {
	var matched bool
	var fails []string

Out:
	switch c.Match {
	case "all":
		for _, corc := range c.Conditions {
			b, _ := matchCompositeOrCondition(root, corc)

			if !b {
				matched = false
				fails = []string{}
				break Out
			}
		}
		matched = true
		fails = []string{}
		break Out
	case "any":
		for _, corc := range c.Conditions {
			b, _ := matchCompositeOrCondition(root, corc)

			if b {
				matched = true
				fails = []string{}
				break Out
			}
		}
		matched = false
		fails = []string{}
		break Out
	}
	return matched, fails
}

func matchCondition(root map[string]interface{}, c Condition) (bool, []string) {
	var passed bool

	switch c.Operation {
	case "exists":
		_, passed = getFromMapByDotPath(c.Path, root)
		break
	case "notexists":
		_, found := getFromMapByDotPath(c.Path, root)
		passed = !found
		break
	default:
		if valToCheck, foundVal := getFromMapByDotPath(c.Path, root); foundVal {
			if compFunc, gotComp := comparators[c.Operation]; gotComp {
				passed = compFunc(valToCheck, c.Value)
			}
		}
	}

	return passed, []string{}
}
