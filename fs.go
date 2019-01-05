package nanafshi

import (
	"errors"
	"fmt"
	"os"
	"time"

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
	Services []Service
}

func (r Root) Stat() *NodeInfo {
	return &r.NodeInfo
}

func (r Root) ListNodes() ([]Node, error) {
	now := time.Now()
	nodes := make([]Node, 0, len(r.Services))

	for _, s := range r.Services {
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
	Service
}

func (d ServiceDir) Stat() *NodeInfo {
	return &d.NodeInfo
}

func (d ServiceDir) ListNodes() ([]Node, error) {
	now := time.Now()
	nodes := make([]Node, 0, len(d.Service.Files))

	for _, f := range d.Service.Files {
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
	File
}

func (f CommandFile) Stat() *NodeInfo {
	return &f.NodeInfo
}

func (f CommandFile) ListNodes() ([]Node, error) {
	return nil, errProtocol
}

func (f CommandFile) ReadFile(ctx *fuse.Context) ([]byte, error) {
	cmd, err := f.ReadCommand.Command.Build(config.Shell, f.makeEnv(ctx))
	if err != nil {
		return nil, err
	}
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (f CommandFile) WriteFile(data []byte, ctx *fuse.Context) error {
	env := f.makeEnv(ctx)
	env = append(env, "FUSE_STDIN="+string(data))
	cmd, err := f.WriteCommand.Command.Build(config.Shell, env)
	if err != nil {
		return err
	}

	if f.WriteCommand.Async {
		err = cmd.Start()
	} else {
		err = cmd.Run()
	}

	if err != nil {
		return err
	}
	return nil
}

func (f CommandFile) makeEnv(ctx *fuse.Context) []string {
	return []string{
		"FUSE_FILENAME=" + f.File.Name,
		"FUSE_OPENPID=" + fmt.Sprint(ctx.Pid),
		"FUSE_OPENUID=" + fmt.Sprint(ctx.Uid),
		"FUSE_OPENGID=" + fmt.Sprint(ctx.Gid),
	}
}
