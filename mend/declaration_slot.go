package mend

import (
	"bufio"
)

type SlotDeclaration struct {
}

func (declaration SlotDeclaration) Match() (string, string) {
	return "@slot", ""
}

func (declaration SlotDeclaration) Parse(in *bufio.Reader) (Declaration, error) {
	return declaration, nil
}

func (declaration SlotDeclaration) Validate() error {
	return nil
}
