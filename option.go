package nanafshi

import flags "github.com/jessevdk/go-flags"

type Option struct {
	ConfPath string `short:"c" long:"config" description:"Config filepath" default:"/etc/nanafshi/config.yml"`
	Verbose  bool   `short:"v" long:"verbose" description:"Show verbose debug information"`
	Version  bool   `short:"V" long:"version" description:"Show version"`
}

func newOptionParser(opts *Option) *flags.Parser {
	parser := flags.NewParser(opts, flags.Default)
	parser.Name = "nanafshi"
	parser.Usage = "[OPTIONS] MOUNTPOINT"
	return parser
}
