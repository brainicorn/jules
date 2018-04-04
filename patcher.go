package jules

import (
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

	if err == nil && len(juleset) < 1 {
		err = fmt.Errorf("No rules found")
	}

	return &JulesPatcher{Rules: juleset}, err
}
