package nanafshi

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/tokibi/nanafshi/config"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
)

type NodeInfo struct {
	Name     string
	Size     int64
	Mode     os.FileMode
	Creation time.Time
	LastMod  time.Time
}

func (i NodeInfo) IsDir() bool {
	return i.Mode&os.ModeDir != 0
}

type Node interface {
	nodefs.Node
	Stat() *NodeInfo
	ListNodes() ([]Node, error)
}

func NewNode() nodefs.Node {
	return nodefs.NewDefaultNode()
}

var errProtocol = errors.New("not implemented")

type Root struct {
	nodefs.Node
	NodeInfo
	Services []config.Service
}

func (n Root) Stat() *NodeInfo {
	return &n.NodeInfo
}

func (n Root) ListNodes() ([]Node, error) {
	now := time.Now()
	nodes := make([]Node, 0, len(n.Services))

	for _, s := range n.Services {
		node := &ServiceDir{
			Node: NewNode(),
			NodeInfo: NodeInfo{
				Name:     s.Name,
				Mode:     os.ModeDir | 0777,
				Creation: now,
				LastMod:  now,
			},
			Service: s,
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

type ServiceDir struct {
	nodefs.Node
	NodeInfo
	Service config.Service
}

func (n ServiceDir) Stat() *NodeInfo {
	return &n.NodeInfo
}

func (n ServiceDir) ListNodes() ([]Node, error) {
	now := time.Now()
	nodes := make([]Node, 0, len(n.Service.Files))

	for _, f := range n.Service.Files {
		node := &CommandFile{
			Node: NewNode(),
			NodeInfo: NodeInfo{
				Name:     f.Name,
				Mode:     0777,
				Creation: now,
				LastMod:  now,
			},
			File: f,
		}
		nodes = append(nodes, node)
	}
	return nodes, nil
}

type CommandFile struct {
	nodefs.Node
	NodeInfo
	config.File
}

func (n CommandFile) Stat() *NodeInfo {
	return &n.NodeInfo
}

func (n CommandFile) ListNodes() ([]Node, error) {
	return nil, errProtocol
}

func (n CommandFile) ReadFile(ctx *fuse.Context) ([]byte, error) {
	cmd, err := n.ReadCommand.Command.Build(conf.Shell, n.makeEnv(ctx))
	if err != nil {
		return nil, err
	}
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (n CommandFile) WriteFile(data []byte, ctx *fuse.Context) error {
	env := n.makeEnv(ctx)
	env = append(env, "FUSE_STDIN="+string(data))
	cmd, err := n.WriteCommand.Command.Build(conf.Shell, env)
	if err != nil {
		return err
	}

	if n.WriteCommand.Async {
		err = cmd.Start()
	} else {
		err = cmd.Run()
	}

	if err != nil {
		return err
	}
	return nil
}

func (n CommandFile) makeEnv(ctx *fuse.Context) []string {
	return []string{
		"FUSE_FILENAME=" + n.Stat().Name,
		"FUSE_OPENPID=" + fmt.Sprint(ctx.Pid),
		"FUSE_OPENUID=" + fmt.Sprint(ctx.Uid),
		"FUSE_OPENGID=" + fmt.Sprint(ctx.Gid),
	}
}
