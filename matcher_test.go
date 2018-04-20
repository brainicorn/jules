package jules_test

import (
	"fmt"
	"testing"

	"github.com/brainicorn/jules"
)

const (
	simplePayload            = `{"foo":"bar"}`
	nestedPayload            = `{"stuff":{"foo":"bar"}}`
	deepNestedPayload        = `{"stuff":{"morestuff":{"foo":"bar"}}}`
	arrayRootMatchPayload        = `{"stuff":{"morestuff":[{"foo":"bar"},{"foo":"bar"}]}}`
	arrayRootPayload        = `{"stuff":{"morestuff":[{"foo":"bar"},{"foo":"jar"}]}}`
	systemBobNamePayload     = `{"website":{"user":{"system":true,"name":"Bob"}}}`
	systemBobUsernamePayload = `{"website":{"user":{"system":true,"username":"Bob"}}}`
	systemNoNamePayload      = `{"website":{"user":{"system":true}}}`
	humanBobPayload          = `{"website":{"user":{"system":false,"name":"Bob"}}}`
	badBobPayload            = `{"website":{"user":{"system":true,"name":"Bob", "bad":true}}}`

	fooEQBarRules  = `[{"condition":{"match":"all", "conditions": [{"path": "foo","op": "eq", "value": "bar"}]}}]`
	fooEQBarSingle = `[{"condition":{"path": "foo","op": "eq", "value": "bar"}}]`
	emptyRules     = `[{"condition":{"match":"all", "conditions": []}}]`
	userNameRules  = `[{"condition":{"match":"all", "conditions": [{"path":"user.system","op":"eq","value":true},{"path":"user.bad","op":"notexists"},{"match":"any","conditions":[{"path":"user.name","op":"exists"},{"path":"user.username","op":"exists"}]}]}}]`
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
		"No Rules",
		"",
		"[{}]",
		simplePayload,
		false,
		fmt.Errorf("No rules found"),
	},
		{
			"Bad Payload",
			"",
			fooEQBarRules,
			"yup",
			false,
			fmt.Errorf("payload must be a json object or a json array of objects"),
		},
		{
			"Root Level",
			"",
			fooEQBarRules,
			simplePayload,
			true,
			nil,
		},
		{
			"Root Level Single Condition",
			"",
			fooEQBarSingle,
			simplePayload,
			true,
			nil,
		},
		{
			"Array Root Level",
			"",
			fooEQBarRules,
			"[" + simplePayload + "]",
			true,
			nil,
		},
		{
			"Nested Root",
			"stuff",
			fooEQBarRules,
			nestedPayload,
			true,
			nil,
		},
		{
			"Deep Nested Root",
			"stuff.morestuff",
			fooEQBarRules,
			deepNestedPayload,
			true,
			nil,
		},
		{
			"Array Root No Match",
			"stuff.morestuff",
			fooEQBarRules,
			arrayRootPayload,
			false,
			nil,
		},
		{
			"Array Root Match",
			"stuff.morestuff",
			fooEQBarRules,
			arrayRootMatchPayload,
			true,
			nil,
		},
		{
			"Missing Root",
			"notfound",
			fooEQBarRules,
			nestedPayload,
			false,
			fmt.Errorf("root object not found at path 'notfound'"),
		},
		{
			"System Bob Name",
			"website",
			userNameRules,
			systemBobNamePayload,
			true,
			nil,
		},
		{
			"System Bob Username",
			"website",
			userNameRules,
			systemBobUsernamePayload,
			true,
			nil,
		},
		{
			"System No Name",
			"website",
			userNameRules,
			systemNoNamePayload,
			false,
			nil,
		},
		{
			"Human Bob Name",
			"website",
			userNameRules,
			humanBobPayload,
			false,
			nil,
		},
		{
			"Bad Bob User",
			"website",
			userNameRules,
			badBobPayload,
			false,
			nil,
		},
}

func TestMatchAt(t *testing.T) {
	var err error
	var matcher *jules.JulesMatcher

	for _, mt := range matchTests {
		var matched bool
		matcher, err = jules.NewMatcher([]byte(mt.rules))

		if err == nil {
			matched, err = matcher.MatchAt([]byte(mt.payload), mt.rootPath)
		}

		if err != nil && (mt.expectedError == nil || (err.Error() != mt.expectedError.Error())) {
			t.Errorf("%s - %s", mt.name, err)
		}

		if mt.expectedError != nil && err == nil {
			t.Errorf("%s - expected error '%s' but was nil", mt.name, mt.expectedError)
		}

		if matched != mt.expectedMatch {
			t.Errorf("%s - mismatch matched for path: expected '%t', got '%t'", mt.name, mt.expectedMatch, matched)
		}
	}
}

func TestMatch(t *testing.T) {
	mt := struct {
		name          string
		rootPath      string
		rules         string
		payload       string
		expectedMatch bool
		expectedError error
	}{
		"Match Test",
		"",
		fooEQBarRules,
		simplePayload,
		true,
		nil,
	}

	var err error
	var matcher *jules.JulesMatcher
	var matched bool
	matcher, err = jules.NewMatcher([]byte(mt.rules))

	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	matched, err = matcher.Match([]byte(mt.payload))

	if err != nil && (mt.expectedError == nil || (err.Error() != mt.expectedError.Error())) {
		t.Errorf("%s - %s", mt.name, err)
	}

	if matched != mt.expectedMatch {
		t.Errorf("%s - mismatch matched for path: expected '%t', got '%t'", mt.name, mt.expectedMatch, matched)
	}
}
