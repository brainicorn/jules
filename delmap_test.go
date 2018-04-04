package jules

import "testing"

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
		"Non-Object Array",
		"foo[2]",
		map[string]interface{}{"foo": []string{"la","fa","bar","ja"}},
		false,
	},
	{
		"Deeply Nested Object Array Value",
		"foo.doo.moo.loo[0].name",
		map[string]interface{}{"foo": map[string]interface{}{"doo": map[string]interface{}{"moo": map[string]interface{}{"loo": []map[string]interface{}{{"name": "bar"}}}}}},
		true,
	},
}

func TestDeleteFromMap(t *testing.T) {
	for _, dt := range delTests {
		deleted := deleteFromMapByDotPath(dt.path, dt.data)

		if deleted != dt.expectedDelete {
			t.Errorf("%s - mismatch delete for path '%s': expected '%t', got '%t'", dt.name, dt.path, dt.expectedDelete, deleted)
		}

		_, exists := valueFromMapByDotPath(dt.path, dt.data)

		if deleted && exists {
			t.Errorf("%s - value should have been deleted but exists for path '%s'", dt.name, dt.path)
		}
	}
}
