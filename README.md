# binstubs

_Automatically generate binstubs using tools.go_

Binstubs is a nearly zero-configuration way to generate binstubs to execute
the version of a go command tracked in `go.mod`.

So instead of running `go run tool@version` which is verbose and causes 
difficulty when upgrading versions, or `go run tool`, which might use the
wrong version of a tool, you just run `bin/tool`.

1. Create a `tools.go` in the root of your project that looks like this,
   pointing to the executable program path for each dependency:

       package main

       import (
       	_ "github.com/brasic/binstubs"
       	_ "github.com/path/to/dep/cmd/something"
       	_ "github.com/golang-migrate/migrate/v4/cmd/migrate"
       )

2. Run `go run github.com/brasic/binstubs` to create corresponding shell
   scripts in `bin/`. You can also add a `go:generate` comment to the top
   of tools.go instead.

3. If you need to import a tool that you don't want a binstub for, add an
   inline `binstub:ignore` comment after the import.

4. If you need extra flags passed to `go run`, add an inline
   `binstub:args="ARGS"` comment after the import.
