package mend

import (
	"bufio"
	"io"
)

type IncludeDeclaration struct {
	Filename string
	Params   *Parameters
}

func (declaration IncludeDeclaration) Match() (string, string) {
	return "@include", ""
}

func (declaration IncludeDeclaration) Parse(in *bufio.Reader) (Declaration, error) {
	var err error
	declaration.Filename, err = readToken(in)

	data, _ := io.ReadAll(in)
	if len(data) != 0 {
		declaration.Params, err = NewParameters(string(data))
	}

	return declaration, err
}

func (declaration IncludeDeclaration) Validate() error {
	if declaration.Filename == "" {
		return ErrNoFilename
	}
	return nil
}
