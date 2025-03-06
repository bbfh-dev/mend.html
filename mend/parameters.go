package mend

import (
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

type Parameters struct {
	root string
	json string
}

func NewParameters(json string) (*Parameters, error) {
	if !gjson.Valid(json) {
		return nil, fmt.Errorf("Invalid JSON format of %q", json)
	}

	return &Parameters{json: json}, nil
}

func (params *Parameters) Get(path string) gjson.Result {
	return gjson.Get(strings.Join([]string{params.root, params.json}, "."), path)
}

func (params *Parameters) GetMany(paths ...string) []gjson.Result {
	return gjson.GetMany(strings.Join([]string{params.root, params.json}, "."), paths...)
}

func (params *Parameters) SetRoot(root string) *Parameters {
	params.root = root
	return params
}
