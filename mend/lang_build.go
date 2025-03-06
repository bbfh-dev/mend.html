package mend

import (
	"os"
)

func Build(proc *Processor, out *os.File) error {
	for proc.Scan() {
		// line, indent := proc.Line()

		proc.WriteTo(out)
	}

	return nil
}
