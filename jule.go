package jules

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

type JuleSet []Jule

type Jule struct {
	Actions    []Action  `json:"actions"`
	Conditions Composite `json:"conditions"`
}

// Action is the base interface for an action.
//
// @jsonSchema(
// 	anyOf=["github.com/brainicorn/jules/AddAction"
//  ,"github.com/brainicorn/jules/RemoveAction"]
// )
type Action interface{}

type AddAction struct {
	// @jsonSchema(required=true)
	AddPath string `json:"add"`

	// @jsonSchema(required=true, type="any")
	Value interface{} `json:"value"`
}

type RemoveAction struct {
	// @jsonSchema(required=true)
	RemovePath string `json:"remove"`

	// @jsonSchema(required=true, type="any")
	Value interface{} `json:"value"`
}

// CompositeOrCondition is the base interface for a set of conditions.
//
// @jsonSchema(
// 	anyOf=["github.com/brainicorn/jules/Composite"
//  		,"github.com/brainicorn/jules/Condition"]
// )
type CompositeOrCondition interface{}

type Composite struct {

	// @jsonSchema(required=true, pattern="^all$|^any$|^none$")
	Match string `json:"match"`

	// @jsonSchema(required=true, minitems=1)
	Conditions []CompositeOrCondition `json:"conditions"`
}

type Condition struct {
	// @jsonSchema(required=true)
	Path string `json:"path"`

	// @jsonSchema(required=true, pattern="^eq$|^neq$|^gt$|^lt$|^gte$|^lte$|^exists$|^notexists$")
	Operation string `json:"op"`

	// @jsonSchema(required=true, type="any")
	Value interface{} `json:"value"`
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
	var schemaValidationResult *gojsonschema.Result

	schemaLoader := gojsonschema.NewStringLoader(GithubComBrainicornJulesJuleSet)
	docLoader := gojsonschema.NewBytesLoader(descriptorBytes)

	schemaValidationResult, err = gojsonschema.Validate(schemaLoader, docLoader)

	if err == nil && len(schemaValidationResult.Errors()) > 0 {
		var errBuf bytes.Buffer
		errBuf.WriteString("Error validating jules json:")
		for _, re := range schemaValidationResult.Errors() {
			errBuf.WriteString(fmt.Sprintf("  - %s", re))
		}

		err = errors.New(errBuf.String())
	}

	if err == nil {
		err = json.Unmarshal(descriptorBytes, &juleset)
	}

	b, _ := json.MarshalIndent(juleset, " ", "   ")
	fmt.Println(string(b))
	return juleset, err
}
