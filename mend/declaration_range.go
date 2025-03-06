package mend

import (
	"bufio"
	"errors"
)

type RangeDeclaration struct {
	Variable string
}

func (declaration RangeDeclaration) Match() (string, string) {
	return "#range", "/range"
}

func (declaration RangeDeclaration) Parse(in *bufio.Reader) (Declaration, error) {
	declaration.Variable, _ = in.ReadString(' ')
	return declaration, nil
}

func (declaration RangeDeclaration) Validate() error {
	if declaration.Variable == "" {
		return errors.New("requires at least one non-empty argument: variable")
	}
	return nil
}
