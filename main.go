package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"

	parser "github.com/axetroy/changelog/1_parser"
	extractor "github.com/axetroy/changelog/2_extractor"
	transform "github.com/axetroy/changelog/3_transform"
	generator "github.com/axetroy/changelog/4_generator"
	"github.com/axetroy/changelog/internal/client"
	"github.com/pkg/errors"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func printHelp() {
	fmt.Println(`changelog - a cli to generate changelog from git project

USAGE:
  changelog [OPTIONS] [version]

ARGUMENTS:
  [version]     Optional version or version range.
                1.null.
                  If you do not specify the version, then it will automatically
                  generate a change log from "HEAD~<latest version>" or
                  "HEAD~<earliest commit>" or "<latest version>-<last version>"
                2.single version. eg. v1.2.0
                  Generate a specific version of the changelog.
                3.version range. eg v1.3.0~v1.2.0
                  Generate changelog within the specified range.
                  For more details, please check the following examples.

OPTIONS:
  --help        Print help information.
  --version     Print version information.
  --dir         Specify the directory to be generated.
                The directory should contain a .git folder. defaults to $PWD.
  --tpl         Specify the directory to be generated.

EXAMPLES:
  # generate changelog from HEAD to <latest version>
  $ changelog

  # generate changelog of the specified version
  $ changelog v1.2.0

  # generate changelog within the specified range
  $ changelog v1.3.0~v1.2.0

  # generate changelog from HEAD to specified version
  $ changelog HEAD~v1.3.0

  # generate all changelog
	$ changelog HEAD~

	# generate all changelog
  $ changelog HEAD~

  # generate changelog from two commit hashes
  $ changelog 770ed02~585445d

  # Generate changelog for the specified project
  $ changelog --dir=/path/to/project v1.0.0

SOURCE CODE:
  https://github.com/axetroy/changelog`)
}

func run() error {
	var (
		showHelp    bool
		showVersion bool
		projectDir  string
	)

	cwd, err := os.Getwd()

	if err != nil {
		return errors.WithStack(err)
	}
	flag.StringVar(&projectDir, "dir", cwd, "project dir")
	flag.StringVar(&projectDir, "tpl", cwd, "TODO: generate changelog with template")
	flag.BoolVar(&showHelp, "help", false, "print help information")
	flag.BoolVar(&showVersion, "version", false, "print version information")

	flag.Parse()

	if showHelp {
		printHelp()
		os.Exit(0)
	}

	if showVersion {
		fmt.Printf("%s %s %s\n", version, commit, date)
		os.Exit(0)
	}

	version := flag.Arg(0)

	client, err := client.NewGitClient(projectDir)

	if err != nil {
		return errors.WithStack(err)
	}

	scope, err := parser.Parse(client, version)

	if err != nil {
		return errors.WithStack(err)
	}

	splices, err := extractor.Extract(client, scope)

	if err != nil {
		return errors.WithStack(err)
	}

	ctxs, err := transform.Transform(client, splices)

	if err != nil {
		return errors.WithStack(err)
	}

	output, err := generator.Generate(client, ctxs)

	if err != nil {
		return errors.WithStack(err)
	}

	_, err = io.Copy(os.Stdout, bytes.NewBuffer(output))

	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func main() {
	var (
		err error
	)

	defer func() {
		// if r := recover(); r != nil {
		// 	fmt.Printf("%+v\n", r)
		// 	os.Exit(255)
		// }

		if err != nil {
			fmt.Printf("%+v\n", err)
			os.Exit(255)
		}
	}()

	err = run()
}
