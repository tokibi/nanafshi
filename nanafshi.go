package nanafshi

import (
	"fmt"
	"io"
	"os"
)

const version = "0.1"

var (
	opts   Option
	config Config
)

type Nanafshi struct {
	Out, Err io.Writer
}

const (
	ExitCodeOK = iota
	ExitCodeError
)

func (n Nanafshi) Run(args []string) int {
	parser := newOptionParser(&opts)
	args, err := parser.ParseArgs(args)
	if err != nil {
		parser.WriteHelp(os.Stderr)
		return ExitCodeError
	}
	if opts.Version {
		fmt.Printf("nanafshi version %s\n", version)
		return ExitCodeOK
	}
	if len(args) < 1 {
		parser.WriteHelp(os.Stderr)
		return ExitCodeError
	}

	config = Config{}
	if err := LoadConfig(opts.ConfPath, &config); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ExitCodeError
	}
	root := config.NewRoot()
	if err = root.MountAndServe(args[0], opts.Verbose); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ExitCodeError
	}

	return ExitCodeOK
}
