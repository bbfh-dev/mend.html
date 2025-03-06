package mend

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type Processor struct {
	Params *Parameters
	Indent int

	file      *os.File
	scanner   *bufio.Scanner
	scanIndex int
}

func NewProcessor(file *os.File, params *Parameters) *Processor {
	file.Seek(0, io.SeekStart)
	return &Processor{
		Params:    params,
		Indent:    0,
		file:      file,
		scanner:   bufio.NewScanner(file),
		scanIndex: 0,
	}
}

func (proc *Processor) Name() string {
	if proc.file == nil {
		return ""
	}
	return proc.file.Name()
}

func (proc *Processor) Scan() bool {
	proc.scanIndex++
	return proc.scanner.Scan()
}

func (proc *Processor) Bytes() []byte {
	return proc.scanner.Bytes()
}

func (proc *Processor) Line() (string, int) {
	text := proc.scanner.Text()
	return strings.TrimSpace(text), countLeadingSpaces(text)
}

func (proc *Processor) WriteTo(out *os.File) {
	out.Write(Indent(proc.Indent))
	out.Write(append(proc.Bytes(), '\n'))
}

func (proc *Processor) PrefixErr(err error) error {
	return fmt.Errorf("%s:%d -> %w", proc.file.Name(), proc.scanIndex, err)
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
