package nanafshi

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/tokibi/nanafshi/config"
)

const version = "0.1.1"

var (
	opts Option
	conf *config.Config
)

type Nanafshi struct {
	Out, Err io.Writer
}

const (
	ExitCodeOK = iota
	ExitCodeError
)

func newRoot(conf *config.Config) *Root {
	now := time.Now()
	return &Root{
		Node: NewNode(),
		NodeInfo: NodeInfo{
			Mode:     os.ModeDir | 0777,
			Creation: now,
			LastMod:  now,
		},
		Services: conf.Services,
	}
}

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

	conf, err = config.LoadFile(opts.ConfPath)
	if err != nil {
		fmt.Fprintln(n.Out, err)
		return ExitCodeError
	}

	root := newRoot(conf)
	if err = root.MountAndServe(args[0], opts.Verbose); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ExitCodeError
	}

	return ExitCodeOK
}
