package jules

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type JulesMatcher struct {
	Rules JuleSet
}

func NewMatcher(jrules json.RawMessage) (*JulesMatcher, error) {
	var err error
	var juleset JuleSet

	juleset, err = validateJuleSet(jrules)

	if err == nil && len(juleset) < 1 {
		err = fmt.Errorf("No rules found")
	}

	return &JulesMatcher{Rules: juleset}, err
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
		// loop over each object and grab our starting point
		for _, obj := range objects {
			root, rootFound := findRoot(obj, rootPath)

			if !rootFound {
				err = fmt.Errorf("root object not found at path '%s'", rootPath)
			}

			if err == nil {
				// we have a root, let's see if we can match everything
				allMatched = applyRules(root, m.Rules)
			}
		}
	}

	return allMatched, err
}

func findRoot(obj map[string]interface{}, rootPath string) (map[string]interface{}, bool) {
	root := obj

	if len(rootPath) > 0 {
		if foundObj, wasFound := valueFromMapByDotPath(rootPath, obj); wasFound {
			if rootObj, ok := foundObj.(map[string]interface{}); ok {
				root = rootObj
			}
		} else {
			root = nil
		}
	}

	return root, (root != nil)
}

func applyRules(root map[string]interface{}, rules JuleSet) bool {
	allMatched := true

	for _, rule := range rules {
		matched, _ := testCompositeOrCondition(root, rule.Conditions)
		if !matched {
			allMatched = false
		}
	}

	return allMatched
}

func testCompositeOrCondition(root map[string]interface{}, corc CompositeOrCondition) (bool, []string) {
	var matched bool
	var fails []string

	switch corc.(type) {
	case Composite:
		matched, fails = testComposite(root, corc.(Composite))

	case Condition:
		matched, fails = testCondition(root, corc.(Condition))
	}

	return matched, fails
}

func testComposite(root map[string]interface{}, c Composite) (bool, []string) {
	switch c.Match {
	case "all":
		for _, corc := range c.Conditions {
			b, _ := testCompositeOrCondition(root, corc)

			if !b {
				return false, []string{}
			}
		}
		return true, []string{}
	case "any":
		for _, corc := range c.Conditions {
			b, _ := testCompositeOrCondition(root, corc)

			if b {
				return true, []string{}
			}
		}
		return false, []string{}
	}
	return false, []string{}
}

func testCondition(root map[string]interface{}, c Condition) (bool, []string) {
	var passed bool

	switch c.Operation {
	case "exists":
		_,passed = valueFromMapByDotPath(c.Path, root)
		break
	case "notexists":
		_,found := valueFromMapByDotPath(c.Path, root)
		passed = !found
		break
	default:
		if valToCheck, foundVal := valueFromMapByDotPath(c.Path, root); foundVal {
			if compFunc, gotComp := comparators[c.Operation]; gotComp {
				passed = compFunc(valToCheck, c.Value)
			}
		}
	}

	return passed, []string{}
}
