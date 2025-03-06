package mend

import (
	"bufio"
	"errors"
	"strings"
)

var (
	OP_EQUAL     = "=="
	OP_NOT_EQUAL = "!="
	OP_HAS       = "has"
	OP_LACKS     = "lacks"
)

var AllOperators = []string{
	OP_EQUAL,
	OP_NOT_EQUAL,
	OP_HAS,
	OP_LACKS,
}

type IfDeclaration struct {
	LeftExpression  string
	Operator        string
	RightExpression string
}

func (declaration IfDeclaration) Match() (string, string) {
	return "#if", "/if"
}

func (declaration IfDeclaration) Parse(in *bufio.Reader) (Declaration, error) {
	declaration.LeftExpression, _ = in.ReadString(' ')
	declaration.LeftExpression = strings.TrimSuffix(declaration.LeftExpression, " ")

	op, _ := in.ReadString(' ')
	declaration.Operator = strings.TrimSuffix(op, " ")

	declaration.RightExpression, _ = readToken(in)

	return declaration, nil
}

func (declaration IfDeclaration) Validate() error {
	if declaration.LeftExpression == "" {
		return errors.New("no variable name on the left side of comparison")
	}

	if declaration.Operator == "" {
		return errors.New("no operation is provided")
	}
	for _, op := range AllOperators {
		if declaration.Operator == op {
			goto hasOperations
		}
	}
	return errors.New("unknown operation")

hasOperations:
	return nil
}
