package jsonfilter

import (
	"encoding/json"
	"fmt"
)

type filterDescriptor struct {
	filter bool
	path []string
}

func key(path []string) string {
	out := ""
	for _,p := range path {
		out += " " + p
	}
	return out
}

type JsonFilter struct {
	buffer   map[string]interface{}
	jsonData []byte
	policy   Grants
}

func (jf *JsonFilter) access(data map[string]interface{}, keys []string) map[string]interface{} {
	tmp, ok := data[keys[0]].(map[string]interface{})
	if !ok {
		return nil
	}
	if len(keys) > 1 {
		return jf.access(tmp, keys[1:])
	}
	return tmp
}

func (jf *JsonFilter) remove(path []string) {
	if len(path) > 1 {
		delete(jf.access(jf.buffer, path[:len(path)-1]), path[len(path)-1])
	} else {
		delete(jf.buffer, path[0])
	}
}

func (jf *JsonFilter) isDenied(role string, roles []string, action string) bool {

	found := false

	for _, r := range roles {
		if r == role {
			found = true
		}
	}

	if action == "allow" {
		return !found
	} else {
		return found
	}
}

func (jf *JsonFilter) Filter(roles []string) ([]byte, error) {

	filterMap := make(map[string]filterDescriptor)

	err := json.Unmarshal(jf.jsonData, &jf.buffer)
	if err != nil {
		return []byte{}, err
	}

	for _, grant := range jf.policy.Grants {
		if _, ok := filterMap[key(grant.Path)]; !ok {
			filterMap[key(grant.Path)] = filterDescriptor{false, grant.Path}
		}

		for _, userRole  := range roles {
			if !jf.isDenied(userRole, grant.Roles, grant.Action) {
				filterMap[key(grant.Path)] = filterDescriptor{true, grant.Path}
			}
		}
	}

	fmt.Println(filterMap)
	for _, filters := range filterMap {
		if filters.filter == false {
			jf.remove(filters.path)
		}
	}


	j, err := json.Marshal(jf.buffer)
	if err != nil {
		return []byte{}, err
	}

	return j, nil
}

func New(jsonData []byte, policy Grants) JsonFilter {
	buffer := make(map[string]interface{})

	return JsonFilter{
		buffer:   buffer,
		jsonData: jsonData,
		policy:   policy,
	}
}

func NewFromFile(jsonData []byte, policyFile string) (JsonFilter, error) {
	buffer := make(map[string]interface{})
	policy, err := readConfig(policyFile)
	if err != nil {
		return JsonFilter{}, err
	}
	return JsonFilter{
			buffer:   buffer,
			jsonData: jsonData,
			policy:   policy,
		},
		nil
}
