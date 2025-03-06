package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/bbfh-dev/mend.html/mend"
	"github.com/bbfh-dev/parsex/parsex"
)

var Version string

var CLI = parsex.New("mend", Program, []parsex.Arg{
	{Name: "version", Match: "--AUTO,-v", Desc: "print version and exit"},
	{
		Name:   "build",
		Match:  "build",
		Desc:   "builds an html file. Outputs to STDOUT when no output is provided",
		Branch: BuildCLI,
	},
})

var BuildCLI = parsex.New("build", BuildProgram, []parsex.Arg{
	{
		Name:  "set",
		Match: "--AUTO,-s",
		Desc:  "specify parameter map in the following format: `{\"key\":\"value\",\"key\":[{\"arg\":3},\"mixed array\"]}`",
		Check: parsex.ValidString,
	},
})

func Program(in parsex.Input, args ...string) error {
	if in.Has("version") {
		fmt.Println("Mend", Version)
		return nil
	}

	return nil
}

func BuildProgram(in parsex.Input, args ...string) error {
	if len(args) == 0 {
		return errors.New("No input file is provided")
	}

	if len(args) > 1 {
		return errors.New("Too many arguments! Only provide a single input file")
	}

	inputFileName := args[0]
	if _, err := os.Stat(inputFileName); os.IsNotExist(err) {
		return err
	}

	params, err := mend.NewParameters(in.Default("set", "{}").(string))
	if err != nil {
		return fmt.Errorf("parsing params: %w", err)
	}

	tempFile, err := os.CreateTemp(os.TempDir(), "mend*.html")
	if err != nil {
		return fmt.Errorf("Creating temporary file: %w", err)
	}
	defer tempFile.Close()
	defer os.Remove(tempFile.Name())

	file, err := os.OpenFile(inputFileName, os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	prepassProc := mend.NewProcessor(file, params)
	if err = mend.Flatten(prepassProc, tempFile, 1); err != nil {
		return fmt.Errorf("flattening file: %w", err)
	}

	proc := mend.NewProcessor(tempFile, params)
	if err = mend.Build(proc, os.Stdout); err != nil {
		return fmt.Errorf("building file: %w", err)
	}

	return nil
}

func main() {
	if err := CLI.FromArgs().Run(); err != nil {
		os.Stderr.Write(append([]byte(err.Error()), '\n'))
		os.Exit(1)
	}
}
