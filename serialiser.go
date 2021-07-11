package main

import (
	"encoding/json"
	"fmt"
)

type T3 struct {
	Visa string `json:"visa"`
	Notes string `json:"note"`
}

type T2 struct {
	Notice int `json:"notice"`
	Salary int `json:"salary"`
	Wp T3 `json:"wp"`
}

type T1 struct {
	Name string `json:"name"`
	Id string `json:"id"'`
	Expectations T2 `json:"expectations"`
}

var t T1 = T1 {
	Name: "mikhail",
	Id: "1123jj1hh123",
	Expectations: T2 {
		Notice: 3,
		Salary: 60000,
		Wp: T3 {
			Visa: "tier2",
			Notes: "foobar",
		},
	},
}

func access(data map[string]interface{}, keys []string) map[string]interface{} {
	tmp, ok := data[keys[0]].(map[string]interface{})
	if (!ok) {
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

func main() {

	j, err := json.Marshal(t)
	if err != nil {
		//return err
	}
	fmt.Println (string(j))

	var mp map[string]interface{}
	err = json.Unmarshal([]byte(j),&mp)
	if err != nil {
		//return err
	}

	//delete (access(mp, []string{"expectations","wp"}),"visa")

	remove (mp, []string{"name"})
	remove (mp, []string{"expectations","wp","visa"})

	j, err = json.Marshal(mp)
	fmt.Println (string(j))
}