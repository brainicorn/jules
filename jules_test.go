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
        "add": "meta.actions.delete",
        "value": "anyone"
      },
      {
        "add": "meta.actions.edit",
        "value": "anyone"
      }
    ],
    "conditions": {
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
        "add": "meta.decorations.labels[0]",
        "value": "APP"
      }
    ],
    "conditions": {
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
