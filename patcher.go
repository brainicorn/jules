package jules

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type JulesPatcher struct {
	Rules JuleSet
}

func NewPatcher(jrules json.RawMessage) (*JulesPatcher, error) {
	var err error
	var juleset JuleSet

	juleset, err = validateJuleSet(jrules)

	if err == nil && (len(juleset) < 1 || numRootActions(juleset) < 1) {
		err = fmt.Errorf("No actions found")
	}

	return &JulesPatcher{Rules: juleset}, err
}

func numRootActions(juleset JuleSet) int {
	var num = 0

	for _, jule := range juleset {
		num += len(jule.Actions)
	}

	return num
}

func (p *JulesPatcher) Patch(payload json.RawMessage) (json.RawMessage, bool, error) {
	return p.PatchAt(payload, "")
}

func (p *JulesPatcher) PatchAt(payload json.RawMessage, rootPath string) (json.RawMessage, bool, error) {
	var err error
	var anyPatched bool
	var result []byte
	var objects []map[string]interface{}
	var thingToMarshal interface{}

	if bytes.HasPrefix(payload, []byte("[{")) {
		// if it's an array of objects, we can just unmarshal it
		err = json.Unmarshal(payload, &objects)
		thingToMarshal = objects
	} else if bytes.HasPrefix(payload, []byte("{")) {
		// if it's a single object, we unmarshal it and add it to the objects slice manually
		// this is just so we can always use a slice of objects for consistency
		var single map[string]interface{}
		err = json.Unmarshal(payload, &single)
		if err == nil {
			thingToMarshal = single
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
				// we have a root, let's see if we can patch everything
				anyPatched = applyPatchRules(root, p.Rules)
			}

		}
	}

	if err == nil {
		result, err = json.Marshal(thingToMarshal)
	}
	return result, anyPatched, err
}

func applyPatchRules(root map[string]interface{}, rules JuleSet) bool {
	var anyPatched bool

	for _, rule := range rules {
		applyActions := true

		if len(rule.Actions) > 0 && hasConditions(rule.Condition) {
			applyActions, _ = matchCompositeOrCondition(root, rule.Condition)
			fmt.Println(fmt.Sprintf("conditions met? %t", applyActions))
		}

		if applyActions {
			for _, action := range rule.Actions {
				//TODO do something with errors
				fmt.Println(fmt.Sprintf("applying action: '%+v'", action))
				patched, _ := applyPatch(root, action)

				if patched {
					anyPatched = true
				}
			}
		}
	}

	return anyPatched
}

func applyPatch(root map[string]interface{}, action Action) (bool, error) {
	var err error
	var patched bool

	switch action.(type) {
	case ValueAction:
	fmt.Println("value...")
		patched, err = applyValuePatch(root, action.(ValueAction))
	case PathAction:
	fmt.Println("path...")
		patched, err = applyPathPatch(root, action.(PathAction))
	case FromToAction:
	fmt.Println("from...")
		patched, err = applyFromToPatch(root, action.(FromToAction))

	default:
	fmt.Println(fmt.Sprintf("type: %+V", action))
	}

	return patched, err
}

func applyValuePatch(root map[string]interface{}, action ValueAction) (bool, error) {
	var patched bool

	switch action.Operation {
	case "add":
		fmt.Println("adding...")
		patched = addToMapByDotPath(action.Path, root, action.Value)
	case "replace":
		patched = replaceInMapByDotPath(action.Path, root, action.Value)
	}

	return patched, nil
}

func applyPathPatch(root map[string]interface{}, action PathAction) (bool, error) {
	var patched bool

	switch action.Operation {
	case "remove":
		patched = deleteFromMapByDotPath(action.Path, root)
	}
	return patched, nil
}

func applyFromToPatch(root map[string]interface{}, action FromToAction) (bool, error) {
	var patched bool

	switch action.Operation {
	case "move":
		patched = moveInMapByDotPath(action.From, action.To, root)
	case "copy":
		patched = copyInMapByDotPath(action.From, action.To, root)
	}

	return patched, nil
}
