package mend

import (
	"bufio"
	"io"
	"strings"
)

type ExtendDeclaration struct {
	Filename string
	Params   *Parameters
}

func (declaration ExtendDeclaration) Match() (string, string) {
	return "#extend", "/extend"
}

func (declaration ExtendDeclaration) Parse(in *bufio.Reader) (Declaration, error) {
	var err error

	declaration.Filename, _ = in.ReadString(' ')
	declaration.Filename = strings.TrimSuffix(declaration.Filename, " ")

	data, _ := io.ReadAll(in)
	if len(data) != 0 {
		declaration.Params, err = NewParameters(string(data))
	}

	return declaration, err
}

func (declaration ExtendDeclaration) Validate() error {
	if declaration.Filename == "" {
		return ErrNoFilename
	}
	return nil
}
