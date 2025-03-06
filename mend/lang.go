package mend

import (
	"bufio"
	"bytes"
	"io"
	"path/filepath"
	"strings"
	"unicode"
)

func Indent(indent int) []byte {
	return bytes.Repeat([]byte{' '}, indent)
}

// readToken reads the next token from reader.
// A token is either a quoted string (if it begins with a `"` or `'`)
// or a word consisting of non-whitespace characters.
func readToken(reader *bufio.Reader) (string, error) {
	// Skip any leading whitespace.
	for {
		char, err := reader.ReadByte()
		if err != nil {
			return "", err
		}
		if !unicode.IsSpace(rune(char)) {
			// Put back the non-whitespace character.
			if err := reader.UnreadByte(); err != nil {
				return "", err
			}
			break
		}
	}

	// Check the first character.
	char, err := reader.ReadByte()
	if err != nil {
		return "", err
	}

	var builder strings.Builder
	if quote := char; quote == '"' || quote == '\'' {
		for {
			char, err = reader.ReadByte()
			if err != nil {
				return "", err
			}
			if char == quote {
				break
			}
			builder.WriteByte(char)
		}
	} else {
		builder.WriteByte(char)
		for {
			char, err = reader.ReadByte()
			if err != nil {
				if err == io.EOF {
					break
				}
				return "", err
			}
			if unicode.IsSpace(rune(char)) {
				break
			}
			builder.WriteByte(char)
		}
	}

	return builder.String(), nil
}

func getRelativePath(anchor string, path string) string {
	baseDir := filepath.Dir(anchor)
	return filepath.Clean(filepath.Join(baseDir, path))
}
