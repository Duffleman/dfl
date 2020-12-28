package auth

import (
	"testing"
)

var suite = []struct {
	Name   string
	Action string
	Scopes string
	Result bool
}{
	{
		Name:   "match up",
		Action: "dflauth:login",
		Scopes: "dflauth:*",
		Result: true,
	},
	{
		Name:   "match exact",
		Action: "dflauth:*",
		Scopes: "dflauth:*",
		Result: true,
	},
	{
		Name:   "do no match down",
		Action: "dflauth:*",
		Scopes: "dflauth:login",
		Result: false,
	},
	{
		Name:   "do not match unrelated",
		Action: "dflauth:list",
		Scopes: "dflauth:login",
		Result: false,
	},
	{
		Name:   "match in set",
		Action: "dflauth:login",
		Scopes: "dflauth:login dflauth:list dflimg:upload",
		Result: true,
	},
}

func TestCan(t *testing.T) {
	for _, test := range suite {
		res := Can(test.Action, test.Scopes)
		if res != test.Result {
			t.Errorf("unexpected \"Can\" result in test \"%s\", expected %t, got %t", test.Name, test.Result, res)
		}
	}
}
