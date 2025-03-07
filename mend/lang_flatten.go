package mend

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

const MAX_RECURSION_DEPTH = 32
const INDENT = 4

var (
	ErrProcIsNil = errors.New("Processor is nil")
	ErrOutIsNil  = errors.New("Out is nil")
	ErrRecursive = errors.New("Max recusion (include/extend) limit is reached!")
	FoundSlot    = errors.New("(Expected behavior) encountered a slot block")
)

func Flatten(proc *Processor, out *os.File, depth int) error {
	if depth > MAX_RECURSION_DEPTH {
		return ErrRecursive
	}

	var discardContent bool
	nestedProcs := NewStack[*Processor]()

	for proc.Scan() {
		line, lineIndent := proc.Line()
		declaration, opening, err := ParseDeclaration(line)
		if err != nil {
			return proc.PrefixErr(fmt.Errorf("parsing statement: %w", err))
		}

		if declaration == nil && !discardContent {
			proc.WriteTo(out)
		}

		switch declaration := declaration.(type) {
		case ExtendDeclaration:
			if discardContent {
				continue
			}
			if opening {
				filepath := getRelativePath(proc.Name(), declaration.Filename)
				file, err := os.OpenFile(filepath, os.O_RDONLY, os.ModePerm)
				if err != nil {
					return proc.PrefixErr(err)
				}
				defer file.Close()

				nestedProc := NewProcessor(file, proc.Params)
				nestedProc.Indent = proc.Indent + lineIndent
				nestedProcs.Add(nestedProc)

				err = Flatten(nestedProc, out, depth+1)
				if err == nil {
					nestedProcs.Pop()
				}
				if err == FoundSlot {
					proc.Indent += INDENT
				} else {
					return proc.PrefixErr(err)
				}
			} else {
				if nestedProcs.Length() == 0 {
					return proc.PrefixErr(errors.New("Cannot close /extend block because it was never opened!"))
				}
				// Pop the last nested processor and continue its flattening
				sub := nestedProcs.Pop()
				if err := Flatten(sub, out, depth+1); err != nil {
					return proc.PrefixErr(err)
				}
			}

		case SlotDeclaration:
			if discardContent {
				continue
			}
			return FoundSlot

		case IncludeDeclaration:
			if discardContent {
				continue
			}
			if opening {
				filepath := getRelativePath(proc.Name(), declaration.Filename)
				file, err := os.OpenFile(filepath, os.O_RDONLY, os.ModePerm)
				if err != nil {
					return proc.PrefixErr(err)
				}
				defer file.Close()

				if declaration.Params == nil {
					declaration.Params, _ = NewParameters("{}")
				}
				nestedProc := NewProcessor(file, proc.Params)
				nestedProc.Indent = proc.Indent + lineIndent
				if err = Flatten(nestedProc, out, depth+1); err != nil {
					return proc.PrefixErr(err)
				}
			}

		case IfDeclaration:
			if opening {
				if discardContent {
					continue
				}
				discardContent = !declaration.Run(proc.Params)
			} else {
				discardContent = false
			}

		case RangeDeclaration:
			if opening {
				var builder strings.Builder

			scanning:
				for proc.Scan() {
					line, lineIndent := proc.Line()
					declaration, opening, err := ParseDeclaration(line)
					if err != nil {
						return proc.PrefixErr(fmt.Errorf("parsing statement: %w", err))
					}

					switch declaration.(type) {
					case RangeDeclaration:
						if !opening {
							break scanning
						}
					}

					builder.Write(Indent(lineIndent))
					builder.Write(proc.Bytes())
					builder.WriteRune('\n')
				}

				result := proc.Params.Get(declaration.Variable)
				if !result.Exists() {
					return proc.PrefixErr(fmt.Errorf("trying to range over a non-existant value %q", declaration.Variable))
				}
				if !result.IsArray() {
					return proc.PrefixErr(fmt.Errorf("trying to range over a value %q which isn't an array", declaration.Variable))
				}
				for i := range result.Array() {
					out.WriteString(strings.ReplaceAll(
						builder.String(),
						"#.",
						fmt.Sprintf(".%s.%d.", declaration.Variable, i),
					))
				}
			}

		}
	}

	return nil
}
