package main

import (
	"encoding/json"
	"io/ioutil"
)

type Grant struct {
	Path   []string `json:"path"`
	Roles  []string `json:"roles"`
	Action string   `json:"action"`
}

type Grants struct {
	Grants []Grant `json:"grants"`
}

func readConfig(path string) (*Grants, error) {
	//TODO: some sanitisation of the path
	file, err := ioutil.ReadFile(path)
	if err != nil {

	}

	data := Grants{}

	err = json.Unmarshal([]byte(file), &data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
