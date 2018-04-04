package jules_test

import (
	"fmt"
	"testing"

	"github.com/brainicorn/jules"
)

const (
	simplePayload     = `{"foo":"bar"}`
	nestedPayload     = `{"stuff":{"foo":"bar"}}`
	deepNestedPayload = `{"stuff":{"morestuff":{"foo":"bar"}}}`
	simpleRules       = `[{"conditions":{"match":"any", "conditions": [{"path": "foo","op": "eq", "value": "bar"}]}}]`
)

var matchTests = []struct {
	name          string
	rootPath      string
	rules         string
	payload       string
	expectedMatch bool
	expectedError error
}{
	{
		"Root Level",
		"",
		simpleRules,
		simplePayload,
		true,
		nil,
	},
	{
		"Nested Root",
		"stuff",
		simpleRules,
		nestedPayload,
		true,
		nil,
	},
	{
		"Deep Nested Root",
		"stuff.morestuff",
		simpleRules,
		deepNestedPayload,
		true,
		nil,
	},
	{
		"Missing Root",
		"notfound",
		simpleRules,
		nestedPayload,
		false,
		fmt.Errorf("root object not found at path 'notfound'"),
	},
}

func TestMatcher(t *testing.T) {
	var err error
	var matcher *jules.JulesMatcher

	for _, mt := range matchTests {
		var matched bool
		matcher, err = jules.NewMatcher([]byte(mt.rules))

		if err != nil {
			t.Error(err)
			t.FailNow()
		}

		matched, err = matcher.MatchAt([]byte(mt.payload), mt.rootPath)

		if err != nil && err.Error() != mt.expectedError.Error() {
			t.Error(err)
		}

		if matched != mt.expectedMatch {
			t.Errorf("%s - mismatch matched for path: expected '%t', got '%t'", mt.name, mt.expectedMatch, matched)
		}
	}
}
