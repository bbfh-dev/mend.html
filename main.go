package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

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
		Desc:  "specify parameter map in the following format: `{key=value;key=[{arg=3},'mixed array']}`",
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
	params := map[string]string{}
	if in.Has("params") {
		pairs := strings.Split(in["params"].(string), ";")
		for _, pair := range pairs {
			parts := strings.Split(pair, "=")
			if len(parts) != 2 {
				return fmt.Errorf("Invalid --param: %q must be `key=value`", pair)
			}
			key := parts[0]
			value := parts[1]

			params[key] = value
		}
	}

	if len(args) == 0 {
		return errors.New("No input file is provided")
	}

	if len(args) > 1 {
		return errors.New("Too many arguments! Only provide a single input file")
	}

	if _, err := os.Stat(args[0]); os.IsNotExist(err) {
		return err
	}

	outputFile, err := os.CreateTemp(os.TempDir(), "mend*.html")
	if err != nil {
		return fmt.Errorf("Creating temporary file: %w", err)
	}
	defer outputFile.Close()

	parser, err := mend.NewParser(args[0], os.Stdout)
	if err != nil {
		return err
	}
	defer parser.Close()

	parser.Params = params
	return parser.Flatten()
}

func main() {
	if err := CLI.FromArgs().Run(); err != nil {
		os.Stderr.Write([]byte(err.Error()))
		os.Exit(1)
	}
}
