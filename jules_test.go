package jules_test

import (
	"fmt"
	"testing"

	"github.com/brainicorn/jules"
)

const (
	rules1 = `[
  {
    "actions": [
      {
        "op": "add",
		"path": "meta.actions.delete",
        "value": "anyone"
      },
      {
        "op": "add",
		"path": "meta.actions.edit",
        "value": "anyone"
      }
    ],
    "condition": {
      "match": "all",
      "conditions": [
        {
          "path": "sender.system",
          "op": "eq",
          "value": "true"
        },
		{
          "path": "sender.funny",
          "op": "eq",
          "value": "false"
        }
      ]
    }
  },
  {
    "actions": [
      {
        "op": "add",
		"path": "meta.decorations.labels[0]",
        "value": "APP"
      }
    ],
    "condition": {
      "match": "all",
      "conditions": [
        {
          "match": "all",
          "conditions": [
            {
              "path": "sender.system",
              "op": "eq",
              "value": "true"
            },
            {
              "path": "sender.notmentionable",
              "op": "eq",
              "value": "false"
            }
          ]
        },
        {
          "path": "sender.notmentionable",
          "op": "eq",
          "value": "false"
        }
      ]
    }
  }
]`
)

func TestSimpleJule(t *testing.T) {
	eng, err := jules.NewPatcher([]byte(rules1))

	if err != nil {
		t.Errorf("error should have been nil but was: %s", err)
	}

	fmt.Println(eng)
}
