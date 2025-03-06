package mend

import (
	"os"
	"strings"
)

const BRACKET_OPEN = "{{"
const BRACKET_CLOSE = "}}"

func Build(proc *Processor, out *os.File) error {
	for proc.Scan() {
		line := string(proc.Bytes())
		i := 0

		for {
			start := strings.Index(line[i:], BRACKET_OPEN)
			if start == -1 {
				out.WriteString(line[i:])
				break
			}

			out.WriteString(line[i : i+start])
			i += start + len(BRACKET_OPEN)

			end := strings.Index(line[i:], BRACKET_CLOSE)
			if end == -1 {
				// This isn't a block, BRACKET_OPEN just happened to be there
				out.WriteString(BRACKET_OPEN)
				out.WriteString(line[i:])
				break
			}

			out.WriteString(ComputeExpression(
				proc.Params,
				strings.TrimSpace(line[i:i+end]),
			))
			i += end + len(BRACKET_CLOSE)
		}

		out.Write([]byte{'\n'})
	}

	return nil
}
