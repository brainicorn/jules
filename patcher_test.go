package jules_test

import (
	"testing"

	"github.com/brainicorn/jules"
)

const (
	patchChangeBarPayload = `{"foo":"bar"}`
	//	nestedPayload            = `{"stuff":{"foo":"bar"}}`
	//	deepNestedPayload        = `{"stuff":{"morestuff":{"foo":"bar"}}}`
	//	systemBobNamePayload     = `{"website":{"user":{"system":true,"name":"Bob"}}}`
	//	systemBobUsernamePayload = `{"website":{"user":{"system":true,"username":"Bob"}}}`
	//	systemNoNamePayload      = `{"website":{"user":{"system":true}}}`
	//	humanBobPayload          = `{"website":{"user":{"system":false,"name":"Bob"}}}`
	//	badBobPayload     = `{"website":{"user":{"system":true,"name":"Bob", "bad":true}}}`

	patchChangeBarRules = `[{"actions":[{"op":"replace","path":"foo","value":"fighters"}],"condition":{"path":"foo","op":"eq","value":"bar"}}]`
	//	userNameRules = `[{"condition":{"match":"all", "conditions": [{"path":"user.system","op":"eq","value":true},{"path":"user.bad","op":"notexists"},{"match":"any","conditions":[{"path":"user.name","op":"exists"},{"path":"user.username","op":"exists"}]}]}}]`

	patchExpectedChangeBarResult = `{"foo":"fighters"}`
)

var patchTests = []struct {
	name           string
	rootPath       string
	rules          string
	payload        string
	expectedPatch  bool
	expectedResult string
	expectedError  error
}{
	//	{
	//		"No Rules",
	//		"",
	//		"[{}]",
	//		simplePayload,
	//		false,
	//		fmt.Errorf("No rules found"),
	//	},
	//	{
	//		"Bad Payload",
	//		"",
	//		fooEQBarRules,
	//		"yup",
	//		false,
	//		fmt.Errorf("payload must be a json object or a json array of objects"),
	//	},
	{
		"Root Level",
		"",
		patchChangeBarRules,
		patchChangeBarPayload,
		true,
		patchExpectedChangeBarResult,
		nil,
	},
	//	{
	//		"Array Root Level",
	//		"",
	//		fooEQBarRules,
	//		"[" + simplePayload + "]",
	//		true,
	//		nil,
	//	},
	//	{
	//		"Nested Root",
	//		"stuff",
	//		fooEQBarRules,
	//		nestedPayload,
	//		true,
	//		nil,
	//	},
	//	{
	//		"Deep Nested Root",
	//		"stuff.morestuff",
	//		fooEQBarRules,
	//		deepNestedPayload,
	//		true,
	//		nil,
	//	},
	//	{
	//		"Missing Root",
	//		"notfound",
	//		fooEQBarRules,
	//		nestedPayload,
	//		false,
	//		fmt.Errorf("root object not found at path 'notfound'"),
	//	},
	//	{
	//		"System Bob Name",
	//		"website",
	//		userNameRules,
	//		systemBobNamePayload,
	//		true,
	//		nil,
	//	},
	//	{
	//		"System Bob Username",
	//		"website",
	//		userNameRules,
	//		systemBobUsernamePayload,
	//		true,
	//		nil,
	//	},
	//	{
	//		"System No Name",
	//		"website",
	//		userNameRules,
	//		systemNoNamePayload,
	//		false,
	//		nil,
	//	},
	//	{
	//		"Human Bob Name",
	//		"website",
	//		userNameRules,
	//		humanBobPayload,
	//		false,
	//		nil,
	//	},
	//	{
	//		"Bad Bob User",
	//		"website",
	//		userNameRules,
	//		badBobPayload,
	//		false,
	//		nil,
	//	},
}

func TestPatchAt(t *testing.T) {
	var err error
	var patcher *jules.JulesPatcher
	var result []byte

	for _, pt := range patchTests {
		var patched bool
		patcher, err = jules.NewPatcher([]byte(pt.rules))

		if err == nil {
			result, patched, err = patcher.PatchAt([]byte(pt.payload), pt.rootPath)
		}

		if err != nil && (pt.expectedError == nil || (err.Error() != pt.expectedError.Error())) {
			t.Errorf("%s - %s", pt.name, err)
		}

		if pt.expectedError != nil && err == nil {
			t.Errorf("%s - expected error '%s' but was nil", pt.name, pt.expectedError)
		}

		if patched != pt.expectedPatch {
			t.Errorf("%s - mismatch patched for path: expected '%t', got '%t'", pt.name, pt.expectedPatch, patched)
		}

		if string(result) != pt.expectedResult {
			t.Errorf("%s - mismatch result for path: expected '%s', got '%s'", pt.name, pt.expectedResult, string(result))
		}
	}
}

//func TestPatch(t *testing.T) {
//	pt := struct {
//		name          string
//		rootPath      string
//		rules         string
//		payload       string
//		expectedPatch bool
//		expectedError error
//	}{
//		"Match Test",
//		"",
//		fooEQBarRules,
//		simplePayload,
//		true,
//		nil,
//	}

//	var err error
//	var patcher *jules.JulesPatcher
//	var patched bool
//	patcher, err = jules.NewPatcher([]byte(pt.rules))

//	if err != nil {
//		t.Error(err)
//		t.FailNow()
//	}

//	patched, err = patcher.Patch([]byte(pt.payload))

//	if err != nil && (pt.expectedError == nil || (err.Error() != pt.expectedError.Error())) {
//		t.Errorf("%s - %s", pt.name, err)
//	}

//	if patched != pt.expectedPatch {
//		t.Errorf("%s - mismatch patched for path: expected '%t', got '%t'", pt.name, pt.expectedPatch, patched)
//	}
//}
