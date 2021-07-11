package main

import (
	"encoding/json"
	"fmt"
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

var t T1 = T1{
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

func access(data map[string]interface{}, keys []string) map[string]interface{} {
	tmp, ok := data[keys[0]].(map[string]interface{})
	if !ok {
		return nil
	}
	if len(keys) > 1 {
		return access(tmp, keys[1:])
	}
	return tmp
}

func remove(data map[string]interface{}, path []string) {
	if len(path) > 1 {
		delete(access(data, path[:len(path)-1]), path[len(path)-1])
	} else {
		delete(data, path[0])
	}
}

func isDenied(role string, roles []string, action string) bool {

	found := false

	for _, r := range roles {
		if r==role {
			found = true;
		}
	}

	if action == "allow" {
		return !found
	} else {
		return found
	}
}


func filter(policyFile string, role string, jsonData []byte) ([]byte, error){

	policy, err := readConfig(policyFile)
	if err != nil {
		return []byte{}, err
	}

	var mp map[string]interface{}
	err = json.Unmarshal([]byte(jsonData), &mp)
	if err != nil {
		return []byte{}, err
	}

	for _, grant := range policy.Grants {
		if isDenied(role, grant.Roles, grant.Action) {
			remove(mp, grant.Path)
		}
	}

	j, err := json.Marshal(mp)
	if err != nil {
		return []byte{}, err
	}

	return j, nil
}

func main() {
	j, _ := json.Marshal(t)
	filteredJson, _ := filter("grants.json", "guest", j)
	fmt.Println(string(filteredJson))
}
