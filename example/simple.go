package main

import (
	"encoding/json"
	"fmt"
	jsonfilter "json_filter"
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

func main() {
	j, _ := json.Marshal(t)
	policyFile := "simple.json"

	myFilter, _ := jsonfilter.NewFromFile(j, policyFile)
	newJson, _ := myFilter.Filter("guest")
	fmt.Println(string(newJson))
}
