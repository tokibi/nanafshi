package main

import (
	"os"

	"github.com/tokibi/nanafshi"
)

func main() {
	fs := nanafshi.Nanafshi{Out: os.Stdout, Err: os.Stderr}
	exitCode := fs.Run(os.Args[1:])
	os.Exit(exitCode)
}
