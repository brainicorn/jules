package jules

import (
	"encoding/json"
)

type JuleSet []Jule

type Jule struct {
	Actions []Action `json:"actions"`
	// @jsonSchema(required=true)
	Condition CompositeOrCondition `json:"condition"`
}

// Action is the base interface for an action.
//
// @jsonSchema(
// 	anyOf=["github.com/brainicorn/jules/ValueAction"
//  ,"github.com/brainicorn/jules/PathAction"
//  ,"github.com/brainicorn/jules/FromToAction"]
// )
type Action interface{}

type ValueAction struct {
	// @jsonSchema(required=true, pattern="^add$|^replace$")
	Operation string `json:"op"`

	// @jsonSchema(required=true)
	Path string `json:"path"`

	// @jsonSchema(required=true, type="any")
	Value interface{} `json:"value"`
}

type PathAction struct {
	// @jsonSchema(required=true, pattern="^remove$")
	Operation string `json:"op"`

	// @jsonSchema(required=true)
	Path string `json:"path"`
}

type FromToAction struct {
	// @jsonSchema(required=true, pattern="^move$|^copy$")
	Operation string `json:"op"`

	// @jsonSchema(required=true)
	From string `json:"from"`

	// @jsonSchema(required=true)
	To string `json:"to"`
}

// CompositeOrCondition is the base interface for a set of conditions.
//
// @jsonSchema(
// 	anyOf=["github.com/brainicorn/jules/Composite"
//  		,"github.com/brainicorn/jules/Condition"]
// )
type CompositeOrCondition interface{}

type ConditionSet []CompositeOrCondition
type Composite struct {

	// @jsonSchema(required=true, pattern="^all$|^any$|^none$")
	Match string `json:"match"`

	// @jsonSchema(required=true, minitems=1)
	Conditions ConditionSet `json:"conditions"`
}

type Condition struct {
	// @jsonSchema(required=true)
	Path string `json:"path"`

	// @jsonSchema(required=true, pattern="^eq$|^neq$|^gt$|^lt$|^gte$|^lte$|^exists$|^notexists$")
	Operation string `json:"op"`

	// @jsonSchema(type="any")
	Value interface{} `json:"value"`
}

func (jul *Jule) UnmarshalJSON(data []byte) error {
	var err error
	var stuff map[string]interface{}
	err = json.Unmarshal(data, &stuff)

	if err == nil {
		for k, v := range stuff {
			switch k {
			case "condition":
				vv := v.(map[string]interface{})
				var jsbytes []byte
				jsbytes, err = json.Marshal(v)

				if err == nil {
					if _, hasMatch := vv["match"]; hasMatch {
						var composite Composite
						err = json.Unmarshal(jsbytes, &composite)

						if err == nil {
							jul.Condition = composite
						}
					} else {
						var conditon Condition
						err = json.Unmarshal(jsbytes, &conditon)

						if err == nil {
							jul.Condition = conditon
						}
					}
				}

			case "actions":
				actionSlice := []Action{}
				actions := v.([]interface{})
				for _, a := range actions {
					var jsbytes []byte
					jsbytes, err = json.Marshal(a)

					if err == nil {
						if _, hasValue := a.(map[string]interface{})["value"]; hasValue {
							var typedAction ValueAction
							err = json.Unmarshal(jsbytes, &typedAction)

							if err == nil {
								actionSlice = append(actionSlice, typedAction)
							}
						} else if _, hasPath := a.(map[string]interface{})["path"]; hasPath {
							var typedAction PathAction
							err = json.Unmarshal(jsbytes, &typedAction)

							if err == nil {
								actionSlice = append(actionSlice, typedAction)
							}
						} else {
							var typedAction FromToAction
							err = json.Unmarshal(jsbytes, &typedAction)

							if err == nil {
								actionSlice = append(actionSlice, typedAction)
							}
						}
					}

				}
				jul.Actions = actionSlice
			}
		}
	}

	return err
}

// UnmarshalJSON cretaes a template object from a JSON structure
func (cmp *Composite) UnmarshalJSON(data []byte) error {
	var err error
	var stuff map[string]interface{}
	err = json.Unmarshal(data, &stuff)

	if err == nil {
		for k, v := range stuff {
			switch k {
			case "match":
				cmp.Match = v.(string)
			case "conditions":
				cndSlice := []CompositeOrCondition{}
				corcs := v.([]interface{})
				for _, cc := range corcs {
					var jsbytes []byte
					jsbytes, err = json.Marshal(cc)

					if err == nil {
						if _, isCondition := cc.(map[string]interface{})["path"]; isCondition {
							var cond Condition
							err = json.Unmarshal(jsbytes, &cond)

							if err == nil {
								cndSlice = append(cndSlice, cond)
							}
						} else {
							var cmp Composite
							err = json.Unmarshal(jsbytes, &cmp)

							if err == nil {
								cndSlice = append(cndSlice, cmp)
							}
						}
					}

				}
				cmp.Conditions = cndSlice
			}
		}
	}

	return err
}

func validateJuleSet(descriptorBytes []byte) (JuleSet, error) {
	var err error
	var juleset JuleSet

	// BUG https://github.com/xeipuuv/gojsonschema/issues/198
	//	var schemaValidationResult *gojsonschema.Result

	//	schemaLoader := gojsonschema.NewStringLoader(GithubComBrainicornJulesJuleSet)
	//	docLoader := gojsonschema.NewBytesLoader(descriptorBytes)

	//	schemaValidationResult, err = gojsonschema.Validate(schemaLoader, docLoader)

	//	if err == nil && len(schemaValidationResult.Errors()) > 0 {
	//		var errBuf bytes.Buffer
	//		errBuf.WriteString("Error validating jules json:")
	//		for _, re := range schemaValidationResult.Errors() {
	//			errBuf.WriteString(fmt.Sprintf("  - %s", re))
	//		}

	//		err = errors.New(errBuf.String())
	//	}

	if err == nil {
		err = json.Unmarshal(descriptorBytes, &juleset)
	}

	return juleset, err
}
