// Binstubs is a nearly zero-configuration way to generate project binstubs
// that run the correct version of a go command tracked in go.mod.
//
//  1. Create a `tools.go` in the root of your project that looks like this,
//     pointing to the executable program path for each dependency:
//
//     package main
//
//     import (
//     _ "github.com/brasic/binstubs"
//     _ "github.com/path/to/dep/cmd/something"
//     _ "github.com/golang-migrate/migrate/v4/cmd/migrate"
//     )
//
//  2. Run `go run github.com/brasic/binstubs` to create corresponding shell
//     scripts in `bin/`. You can also add a `go:generate` comment to the top
//     of tools.go instead.
//
//  3. If you need to import a tool that you don't want a binstub for, add an
//     inline `binstub:ignore` comment after the import.
//
//  4. If you need extra flags passed to `go run`, add an inline
//     `binstub:args="ARGS"` comment after the import.
package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"text/template"
)

// regexp to split up a line
var lineRe = regexp.MustCompile(`\s+_\s"([^"]+)"\s*/?/?\s*(.*)`)

// regexp to parse the `binstub:` comment
var commentRe = regexp.MustCompile(`binstub:(\w+)=?"?([^"]*)"?\s*`)

var tpl = template.Must(template.New("binstub").Parse(`#!/bin/sh
# Code generated by binstubs. DO NOT EDIT.

exec {{ .GoRunCommand }} {{ .Module }} "$@"
`))

// TemplateArgs are passed to binstubTpl.
type TemplateArgs struct {
	GoRunCommand string
	Module       string
}

// CommentOption is a parsed set of generator options specified via a postfix
// inline comment.
type CommentOption struct {
	Ignore bool
	Args   string
}

func parseComment(comment string) *CommentOption {
	opt := CommentOption{}
	matches := commentRe.FindStringSubmatch(comment)
	if len(matches) > 1 {
		switch matches[1] {
		case "ignore":
			opt.Ignore = true
		case "args":
			if len(matches) != 3 {
				panic("bad syntax for comment " + comment)
			}
			opt.Args = matches[2]
		default:
			panic("bad syntax for comment " + comment)
		}

	}
	return &opt
}

func generateBinstub(module string, comment string) {
	options := parseComment(comment)
	if options.Ignore {
		return
	}
	binstubName := filepath.Base(module)
	binstubPath := filepath.Join("bin", binstubName)

	f, err := os.OpenFile(binstubPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	tplArgs := &TemplateArgs{Module: module}
	if options.Args != "" {
		tplArgs.GoRunCommand = "go run " + options.Args
	} else {
		tplArgs.GoRunCommand = "go run"
	}
	err = tpl.Execute(f, tplArgs)
	if err != nil {
		panic(err)
	}
	fmt.Fprintf(os.Stderr, "wrote %s\n", binstubPath)
}

func main() {
	f, err := os.Open("tools.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)

	count := 0

	os.MkdirAll("bin", os.ModePerm)
	for scanner.Scan() {
		matches := lineRe.FindStringSubmatch(scanner.Text())
		if len(matches) > 1 {
			count++
			generateBinstub(matches[1], matches[2])
		}
	}
	if count == 0 {
		fmt.Println("no imports found in tools.go")
		os.Exit(1)
	}
}
