package jules_test

import (
	"fmt"
	"testing"

	"github.com/brainicorn/jules"
)

const (
	patchFooBarPayload               = `{"foo":"bar"}`
	patchFooBarNestedPayload         = `{"stuff":{"foo":"bar"}}`
	patchFooBarDeepNestedPayload     = `{"stuff":{"morestuff":{"foo":"bar"}}}`
	patchReplaceBarRules             = `[{"actions":[{"op":"replace","path":"foo","value":"fighters"}],"condition":{"path":"foo","op":"eq","value":"bar"}}]`
	patchAddDogRules                 = `[{"actions":[{"op":"add","path":"my.pet.dog","value":true}],"condition":{"path":"foo","op":"eq","value":"bar"}}]`
	patchAddPetsRules                = `[{"actions":[{"op":"add","path":"my.pets[-]","value":"dog"}],"condition":{"path":"foo","op":"eq","value":"bar"}}]`
	patchFooFightersResult           = `{"foo":"fighters"}`
	patchFooFightersNestedResult     = `{"stuff":{"foo":"fighters"}}`
	patchFooFightersDeepNestedResult = `{"stuff":{"morestuff":{"foo":"fighters"}}}`
	patchDogResult                   = `{"foo":"bar","my":{"pet":{"dog":true}}}`
	patchPetsResult                  = `{"foo":"bar","my":{"pets":["dog"]}}`
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
		"Replace Root Level",
		"",
		patchReplaceBarRules,
		patchFooBarPayload,
		true,
		patchFooFightersResult,
		nil,
	},
	{
		"Replace Array Root Level",
		"",
		patchReplaceBarRules,
		"[" + patchFooBarPayload + "]",
		true,
		"[" + patchFooFightersResult + "]",
		nil,
	},
	{
		"Replace Nested Root",
		"stuff",
		patchReplaceBarRules,
		patchFooBarNestedPayload,
		true,
		patchFooFightersNestedResult,
		nil,
	},
	{
		"Replace Deep Nested Root",
		"stuff.morestuff",
		patchReplaceBarRules,
		patchFooBarDeepNestedPayload,
		true,
		patchFooFightersDeepNestedResult,
		nil,
	},
	{
		"Replace Missing Root",
		"notfound",
		patchReplaceBarRules,
		patchFooBarPayload,
		false,
		"",
		fmt.Errorf("root object not found at path 'notfound'"),
	},
	{
		"Add Object Root Level",
		"",
		patchAddDogRules,
		patchFooBarPayload,
		true,
		patchDogResult,
		nil,
	},
	{
		"Add To Array Root Level",
		"",
		patchAddPetsRules,
		patchFooBarPayload,
		true,
		patchPetsResult,
		nil,
	},
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
