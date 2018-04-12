package jules

import (
	"testing"
)

var delTests = []struct {
	name           string
	path           string
	data           map[string]interface{}
	expectedDelete bool
}{
	{
		"Root Level Value",
		"foo",
		map[string]interface{}{"foo": "bar"},
		true,
	},
	{
		"Deeply Nested Value",
		"foo.doo.moo.loo",
		map[string]interface{}{"foo": map[string]interface{}{"doo": map[string]interface{}{"moo": map[string]interface{}{"loo": "bar"}}}},
		true,
	},
	{
		"Missing Root Path",
		"moo.foo.loo",
		map[string]interface{}{"foo": map[string]interface{}{"loo": "bar"}},
		false,
	},
	{
		"Deeply Nested Object Array Value",
		"foo.doo.moo.loo[0].name",
		map[string]interface{}{"foo": map[string]interface{}{"doo": map[string]interface{}{"moo": map[string]interface{}{"loo": []interface{}{map[string]interface{}{"name": "bar"}}}}}},
		true,
	},
{
		"Non-Object Array",
		"foo[3]",
		map[string]interface{}{"foo": []interface{}{"la", "fa", "bar", "ja"}},
		true,
	},
}

func TestDeleteFromMap(t *testing.T) {
	for _, dt := range delTests {
		deleted := deleteFromMapByDotPath(dt.path, dt.data)

		if deleted != dt.expectedDelete {
			t.Errorf("%s - mismatch delete for path '%s': expected '%t', got '%t'", dt.name, dt.path, dt.expectedDelete, deleted)
		}

		val, exists := getFromMapByDotPath(dt.path, dt.data)

		if deleted && exists {
			t.Errorf("%s - value should have been deleted but exists for path '%+v'", dt.name, dt.path, val)
		}
	}
}
