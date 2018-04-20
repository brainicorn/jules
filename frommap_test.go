package jules

import (
	"reflect"
	"testing"
)

var valueTests = []struct {
	name           string
	path           string
	data           map[string]interface{}
	expectedExists bool
	expectedValue  interface{}
}{
	{
		"Root Level Value",
		"foo",
		map[string]interface{}{"foo": "bar"},
		true,
		"bar",
	},
	{
		"Deeply Nested Value",
		"foo.doo.moo.loo",
		map[string]interface{}{"foo": map[string]interface{}{"doo": map[string]interface{}{"moo": map[string]interface{}{"loo": "bar"}}}},
		true,
		"bar",
	},
	{
		"Missing Root Path",
		"moo.foo.loo",
		map[string]interface{}{"foo": map[string]interface{}{"loo": "bar"}},
		false,
		nil,
	},
	{
		"Root Level Array Value",
		"foo[0]",
		map[string]interface{}{"foo": []interface{}{"bar"}},
		true,
		"bar",
	},
	{
		"Deeply Nested Array Value",
		"foo.doo.moo.loo[0]",
		map[string]interface{}{"foo": map[string]interface{}{"doo": map[string]interface{}{"moo": map[string]interface{}{"loo": []interface{}{"bar"}}}}},
		true,
		"bar",
	},
	{
		"Root Level Array Deep Index Value",
		"foo[2]",
		map[string]interface{}{"foo": []interface{}{"la", "fa", "bar", "ja"}},
		true,
		"bar",
	},
	{
		"Deeply Nested Array Deep Index Value",
		"foo.doo.moo.loo[2]",
		map[string]interface{}{"foo": map[string]interface{}{"doo": map[string]interface{}{"moo": map[string]interface{}{"loo": []interface{}{"la", "fa", "bar", "ja"}}}}},
		true,
		"bar",
	},
	{
		"Deeply Nested Array No Index",
		"foo.doo.moo.loo",
		map[string]interface{}{"foo": map[string]interface{}{"doo": map[string]interface{}{"moo": map[string]interface{}{"loo": []interface{}{"la", "fa", "bar", "ja"}}}}},
		true,
		[]interface{}{"la", "fa", "bar", "ja"},
	},
	{
		"Root Level Object Array Value",
		"foo[0].name",
		map[string]interface{}{"foo": []interface{}{map[string]interface{}{"name": "bar"}}},
		true,
		"bar",
	},
	{
		"Deeply Nested Object Array Value",
		"foo.doo.moo.loo[0].name",
		map[string]interface{}{"foo": map[string]interface{}{"doo": map[string]interface{}{"moo": map[string]interface{}{"loo": []interface{}{map[string]interface{}{"name": "bar"}}}}}},
		true,
		"bar",
	},
}

func TestValueFromMap(t *testing.T) {
	for _, vt := range valueTests {
		val, exists := getFromMapByDotPath(vt.path, vt.data)

		if exists != vt.expectedExists {
			t.Errorf("%s - mismatch exists for path '%s': expected '%t', got '%t'", vt.name, vt.path, vt.expectedExists, exists)
		}

		if !reflect.DeepEqual(val, vt.expectedValue) {
			t.Errorf("%s - mismatch value for path '%s': expected '%+v', got '%+v'", vt.name, vt.path, vt.expectedValue, val)
		}
	}
}
