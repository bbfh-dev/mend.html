package mend

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/bbfh-dev/mend.html/mend/assert"
)

type Parser struct {
	Filename string
	Input    *bufio.Scanner
	Output   *os.File
	Params   map[string]string

	file   *os.File
	index  int
	indent int
}

func NewParser(filename string, output *os.File) (*Parser, error) {
	assert.NotNil("NewParser(output)", output)

	inputFile, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	return &Parser{
		Filename: filename,
		Input:    bufio.NewScanner(inputFile),
		Output:   output,
		Params:   map[string]string{},
		file:     inputFile,
		index:    0,
		indent:   0,
	}, nil
}

func (parser *Parser) Close() error {
	assert.NotNil("Parser.file", parser.file)
	return parser.file.Close()
}

func (parser *Parser) SetParam(key, value string) *Parser {
	parser.Params[key] = value
	return parser
}

func (parser *Parser) Flatten() error {
	assert.NotNil("Parser(Input).Flatten", parser.Input)
	assert.NotNil("Parser(Output).Flatten", parser.Output)

	var discardUntilIf bool
	var subParser *Parser

	for parser.scan() {
		line := strings.TrimSpace(parser.Input.Text())
		indent := countLeadingSpaces(parser.Input.Text())
		if isComment(line) {
			comment := extractComment(line)
			if isMendStatement(comment) {
				reader := bufio.NewReader(strings.NewReader(comment))
				command, err := reader.ReadString(' ')
				if err != nil && err != io.EOF {
					return parser.err("reading comment: %s", err.Error())
				}
				command = strings.TrimSpace(command)

				if discardUntilIf {
					if command == "/if" {
						discardUntilIf = false
					}
					continue
				}

				switch command {
				case "#extend":
					token, err := readToken(reader)
					if err != nil {
						if err == io.EOF {
							return parser.err("#extend requires at least one argument (filepath)")
						}
						return parser.err("reading token: %s", err.Error())
					}

					if len(token) == 0 {
						return parser.err("#extend: first argument (filepath) cannot be empty")
					}

					params := map[string]string{}

					var key string
					var value string

					for {
						key, err = reader.ReadString('=')
						if err != nil {
							break
						}
						key = strings.TrimSuffix(key, "=")

						value, err = readToken(reader)
						if err != nil {
							break
						}

						if len(key) != 0 && len(value) != 0 {
							params[key] = value
						}
					}

					subParser, err = NewParser(
						getRelativePath(parser.Filename, token),
						parser.Output,
					)
					if err != nil {
						return parser.err("#extend: %s", err.Error())
					}
					defer subParser.Close()
					subParser.Params = params

					err = subParser.Flatten()
					if err != nil {
						return err
					}
					parser.indent += 4

				case "/extend":
					err = subParser.Flatten()
					if err != nil {
						return err
					}

				case "#if":
					name, err := reader.ReadString(' ')
					if err != nil {
						if err == io.EOF {
							return parser.err("#if requires at least one argument (param name)")
						}
						return parser.err("reading token: %s", err.Error())
					}
					name = strings.TrimSuffix(name, " ")

					op, err := reader.ReadString(' ')
					if err != nil {
						if err == io.EOF {
							return parser.err("#if requires an operator")
						}
						return parser.err("reading token: %s", err.Error())
					}
					op = strings.TrimSuffix(op, " ")

					value, err := readToken(reader)
					if err != nil {
						if err == io.EOF {
							return parser.err("#if requires the comparison value")
						}
						return parser.err("reading token: %s", err.Error())
					}

					param, ok := parser.Params[name]
					if !ok {
						return parser.err("undefined parameter %q", name)
					}
					switch op {
					case "==":
						if param != value {
							discardUntilIf = true
						}
					case "!=":
						if param == value {
							discardUntilIf = true
						}
					default:
						return parser.err(
							"#if statements only support '==' and '!=' operators at the moment",
						)
					}

				case "/if":
					// return parser.err("trying to close an <if> but it was never opened")

				case "@slot":
					return nil

				case "@put":
					token, err := readToken(reader)
					if err != nil {
						if err == io.EOF {
							return parser.err("@put requires at least one argument (variable name)")
						}
						return parser.err("reading token: %s", err.Error())
					}
					parser.writeIndent(indent)
					parser.write([]byte(parser.Params[token]))
					parser.write([]byte{'\n'})

				case "@include":
					token, err := readToken(reader)
					if err != nil {
						if err == io.EOF {
							return parser.err("@include requires at least one argument (filepath)")
						}
						return parser.err("reading token: %s", err.Error())
					}

					if len(token) == 0 {
						return parser.err("@include: first argument (filepath) cannot be empty")
					}

					data, err := os.ReadFile(getRelativePath(parser.Filename, token))
					if err != nil {
						return parser.err("@include: %s", err.Error())
					}

					scanner := bufio.NewScanner(bytes.NewReader(data))
					for scanner.Scan() {
						parser.writeIndent(parser.indent)
						parser.write(scanner.Bytes())
						parser.write([]byte{'\n'})
					}

				case "@alter":
					variable, err := reader.ReadString(' ')
					if err != nil {
						if err == io.EOF {
							return parser.err(
								"@alter requires at least one argument (variable name)",
							)
						}
						return parser.err("reading token: %s", err.Error())
					}
					variable = strings.TrimSuffix(variable, " ")

					op, err := reader.ReadString(' ')
					if err != nil {
						if err == io.EOF {
							return parser.err("@alter requires an operator")
						}
						return parser.err("reading token: %s", err.Error())
					}
					op = strings.TrimSuffix(op, " ")

					value, err := readToken(reader)
					if err != nil {
						if err == io.EOF {
							return parser.err("@alter requires a valid value")
						}
						return parser.err("reading token: %s", err.Error())
					}
					fmt.Println("---", variable, op, value)

				default:
					fmt.Println("===", comment)
				}
				continue
			}
		}

		if discardUntilIf {
			continue
		}

		parser.writeIndent(parser.indent)
		parser.write(parser.Input.Bytes())
		parser.write([]byte{'\n'})
	}

	return nil
}

// Increments index and performs bufio.Scanner.scan()
func (parser *Parser) scan() bool {
	parser.index++
	return parser.Input.Scan()
}

func (parser *Parser) write(bytes []byte) {
	assert.SafeWrite(parser.Output.Write(bytes))
}

func (parser *Parser) err(str string, format ...any) error {
	return fmt.Errorf("%s:%d -> %s", parser.Filename, parser.index, fmt.Sprintf(str, format...))
}

func (parser *Parser) writeIndent(indent int) {
	parser.write(bytes.Repeat([]byte{' '}, indent))
}

// readToken reads the next token from r.
// A token is either a quoted string (if it begins with a '"')
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

func countLeadingSpaces(str string) (count int) {
	for _, rune := range str {
		switch rune {
		case ' ':
			count++
		case '\t':
			count += 4
		default:
			return
		}
	}

	return
}

func isComment(line string) bool {
	return strings.HasPrefix(line, "<!--") && strings.HasSuffix(line, "-->")
}

func isMendStatement(comment string) bool {
	switch comment[0] {
	case '#', '@', '/':
		return true
	default:
		return false
	}
}

func extractComment(line string) string {
	return strings.TrimSpace(line[len("<!--") : len(line)-len("-->")])
}
