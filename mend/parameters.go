package mend

import (
	"fmt"

	"github.com/tidwall/gjson"
)

type Parameters struct {
	json string
}

func NewParameters(json string) (*Parameters, error) {
	if !gjson.Valid(json) {
		return nil, fmt.Errorf("Invalid JSON format of %q", json)
	}

	return &Parameters{json: json}, nil
}

func (params *Parameters) Get(path string) gjson.Result {
	return gjson.Get(params.json, path)
}

func (params *Parameters) GetMany(paths ...string) []gjson.Result {
	return gjson.GetMany(params.json, paths...)
}
