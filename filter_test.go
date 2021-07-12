package jsonfilter

import (
	"encoding/json"
	"github.com/kinbiko/jsonassert"
	"testing"
)

type T3 struct {
	Visa  string `json:"visa"`
	Notes string `json:"note"`
}

type T2 struct {
	Notice int `json:"notice"`
	Salary int `json:"salary"`
	Wp     T3  `json:"wp"`
}

type T1 struct {
	Name         string `json:"name"`
	Id           string `json:"id"'`
	Expectations T2     `json:"expectations"`
}

var mockData T1 = T1{
	Name: "mikhail",
	Id:   "1123jj1hh123",
	Expectations: T2{
		Notice: 3,
		Salary: 60000,
		Wp: T3{
			Visa:  "tier2",
			Notes: "foobar",
		},
	},
}

var denyNameToGuestAndUser Grant = Grant{
	Path:   []string{"name"},
	Roles:  []string{"guest", "user"},
	Action: "deny",
}

var denyAccessToAWholeObject Grant = Grant{
	Path:   []string{"expectations", "wp"},
	Roles:  []string{"guest", "user"},
	Action: "deny",
}

var allowNameToUser Grant = Grant{
	Path:   []string{"name"},
	Roles:  []string{"user"},
	Action: "allow",
}

var policyFilterFirstrLevel Grants = Grants{
	[]Grant{denyNameToGuestAndUser},
}

var policyFilterBlock Grants = Grants{
	[]Grant{denyAccessToAWholeObject},
}

var policyAllow Grants = Grants{
	[]Grant{allowNameToUser},
}

func TestJsonFilter_DoNotFiter(t *testing.T) {
	j, _ := json.Marshal(mockData)

	myFilter := New(j, policyFilterFirstrLevel)
	json, _ := myFilter.Filter("superadmin")

	ja := jsonassert.New(t)
	ja.Assertf(string(j), string(json))

}

func TestJsonFilter_FilterFirstLevelWithTwoRoles(t *testing.T) {
	//prepare a mock json
	j, _ := json.Marshal(mockData)

	myFilter := New(j, policyFilterFirstrLevel)
	filteredJson, _ := myFilter.Filter("user")
	ja := jsonassert.New(t)

	ja.Assertf(
		`{"expectations":{"notice":3,"salary":60000,"wp":{"note":"foobar","visa":"tier2"}},"id":"1123jj1hh123"}`,
		string(filteredJson))

	filteredJson, _ = myFilter.Filter("guest")
	ja.Assertf(
		`{"expectations":{"notice":3,"salary":60000,"wp":{"note":"foobar","visa":"tier2"}},"id":"1123jj1hh123"}`,
		string(filteredJson))
}

func TestJsonFilter_FilterBlock(t *testing.T) {
	//prepare a mock json
	j, _ := json.Marshal(mockData)

	myFilter := New(j, policyFilterBlock)
	filteredJson, _ := myFilter.Filter("user")
	ja := jsonassert.New(t)

	ja.Assertf(
		`{"expectations":{"notice":3,"salary":60000},"id":"1123jj1hh123","name":"mikhail"}
`,
		string(filteredJson))

	filteredJson, _ = myFilter.Filter("guest")
	ja.Assertf(
		`{"expectations":{"notice":3,"salary":60000},"id":"1123jj1hh123","name":"mikhail"}`,
		string(filteredJson))
}

func TestJsonFilter_AllowField(t *testing.T) {
	//prepare a mock json
	j, _ := json.Marshal(mockData)

	myFilter := New(j, policyAllow)
	filteredJson, _ := myFilter.Filter("user")
	ja := jsonassert.New(t)

	ja.Assertf(
		`{"expectations":{"notice":3,"salary":60000,"wp":{"note":"foobar","visa":"tier2"}},"id":"1123jj1hh123","name":"mikhail"}`,
		string(filteredJson))

	filteredJson, _ = myFilter.Filter("guest")
	ja.Assertf(
		`{"expectations":{"notice":3,"salary":60000,"wp":{"note":"foobar","visa":"tier2"}},"id":"1123jj1hh123"}`,
		string(filteredJson))
}
