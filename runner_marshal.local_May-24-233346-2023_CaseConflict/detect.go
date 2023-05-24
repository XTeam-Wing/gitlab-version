package Runner

import (
	"encoding/json"
	"fmt"
)

type Version struct {
	Build    string
	Versions []string
}

type Data struct {
	Versions map[string]Version
}

type Resp struct {
	Hash string `json:"hash"`
}

func (r *Runner) Detect(content string) (build string, result []string, err error) {
	var resp Resp
	err = json.Unmarshal([]byte(content), &resp)
	if err != nil {
		return build, result, err
	}
	if resp.Hash == "" {
		return build, result, fmt.Errorf("hash is empty")
	}
	for k, data := range r.Data {
		if resp.Hash == k {
			for _, v := range data.Versions {
				result = append(result, v)
			}
			return data.Build, result, nil
		}
	}
	return "", nil, err
}
