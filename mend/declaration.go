package mend

import (
	"bufio"
	"errors"
	"strings"
)

var AllDeclarations = []Declaration{
	ExtendDeclaration{},
	IncludeDeclaration{},
	IfDeclaration{},
	RangeDeclaration{},
	SlotDeclaration{},
}

var ErrNoFilename = errors.New("requires at least one non-empty argument: filename")

type Declaration interface {
	Match() (opening string, closing string)
	Parse(*bufio.Reader) (Declaration, error)
	Validate() error
}

func ParseDeclaration(line string) (Declaration, bool, error) {
	line = strings.TrimPrefix(line, "<!--")
	line = strings.TrimSuffix(line, "-->")
	line = strings.TrimSpace(line)

	reader := bufio.NewReader(strings.NewReader(line))
	match, _ := reader.ReadString(' ')
	if len(match) == 0 {
		return nil, true, nil
	}
	match = strings.TrimSuffix(match, " ")

	for _, declaration := range AllDeclarations {
		opening, closing := declaration.Match()

		if len(opening) != 0 && match == opening {
			var err error
			declaration, err = declaration.Parse(reader)
			if err != nil {
				return declaration, true, err
			}
			return declaration, true, declaration.Validate()
		}

		if len(closing) != 0 && match == closing {
			return declaration, false, nil
		}
	}

	return nil, true, nil
}
