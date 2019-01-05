package nanafshi

import (
	"os"
	"time"

	"github.com/hanwen/go-fuse/fuse"
	"github.com/hanwen/go-fuse/fuse/nodefs"
)

func (i NodeInfo) FillAttr(out *fuse.Attr) {
	perm := uint32(i.Mode & os.ModePerm)
	if i.IsDir() {
		out.Mode = fuse.S_IFDIR | perm
	} else {
		out.Mode = fuse.S_IFREG | perm
	}
	out.Size = uint64(i.Size)
	out.Atime = uint64(i.LastMod.Unix())
	out.Mtime = uint64(i.LastMod.Unix())
}

func (i NodeInfo) FillDirEntry(out *fuse.DirEntry) {
	out.Name = i.Name
	out.Mode = uint32(i.Mode & os.ModePerm)
	if i.IsDir() {
		out.Mode |= fuse.S_IFDIR
	}
}

func (r *Root) MountAndServe(path string, debug bool) error {
	opts := &nodefs.Options{
		AttrTimeout:  time.Second,
		EntryTimeout: time.Second,
		Debug:        debug,
	}
	s, _, err := nodefs.MountRoot(path, r, opts)
	if err != nil {
		return err
	}
	s.Serve()
	return nil
}

func (r *Root) Lookup(out *fuse.Attr, name string, ctx *fuse.Context) (*nodefs.Inode, fuse.Status) {
	return lookupName(r, name, out, ctx)
}

func (r Root) GetAttr(out *fuse.Attr, file nodefs.File, ctx *fuse.Context) fuse.Status {
	r.NodeInfo.FillAttr(out)
	return fuse.OK
}

func (r *Root) OpenDir(ctx *fuse.Context) ([]fuse.DirEntry, fuse.Status) {
	return listNodes(r)
}

func (d *ServiceDir) Lookup(out *fuse.Attr, name string, ctx *fuse.Context) (*nodefs.Inode, fuse.Status) {
	return lookupName(d, name, out, ctx)
}

func (d ServiceDir) GetAttr(out *fuse.Attr, file nodefs.File, ctx *fuse.Context) fuse.Status {
	d.NodeInfo.FillAttr(out)
	return fuse.OK
}

func (d *ServiceDir) OpenDir(ctx *fuse.Context) ([]fuse.DirEntry, fuse.Status) {
	return listNodes(d)
}

func (f CommandFile) GetAttr(out *fuse.Attr, file nodefs.File, ctx *fuse.Context) fuse.Status {
	f.NodeInfo.FillAttr(out)
	return fuse.OK
}

func (f CommandFile) OpenDir(ctx *fuse.Context) ([]fuse.DirEntry, fuse.Status) {
	return nil, fuse.EINVAL
}

func (f CommandFile) Open(flags uint32, ctx *fuse.Context) (nodefs.File, fuse.Status) {
	if flags&fuse.O_ANYWRITE != 0 {
		return nodefs.NewDevNullFile(), fuse.OK
	}

	p, err := f.ReadFile(ctx)
	if err != nil {
		return nil, fuse.EIO
	}
	return &nodefs.WithFlags{
		File:      nodefs.NewDataFile(p),
		FuseFlags: fuse.FOPEN_DIRECT_IO,
	}, fuse.OK
}

func (f CommandFile) Truncate(file nodefs.File, size uint64, ctx *fuse.Context) fuse.Status {
	return fuse.OK
}

func (f CommandFile) Write(file nodefs.File, data []byte, off int64, ctx *fuse.Context) (uint32, fuse.Status) {
	err := f.WriteFile(data, ctx)
	if err != nil {
		return 0, fuse.EINVAL
	}
	return uint32(len(data)), fuse.OK
}

func lookupName(node Node, name string, out *fuse.Attr, ctx *fuse.Context) (*nodefs.Inode, fuse.Status) {
	_, status := listNodes(node)
	if status != fuse.OK {
		return nil, status
	}
	c := node.Inode().GetChild(name)
	if c == nil {
		return nil, fuse.ENOENT
	}
	status = c.Node().GetAttr(out, nil, ctx)
	if status != fuse.OK {
		return nil, status
	}
	return c, fuse.OK
}

func listNodes(node Node) ([]fuse.DirEntry, fuse.Status) {
	p := node.Inode()
	nodes, err := node.ListNodes()
	if err != nil {
		return nil, fuse.EIO
	}
	a := make([]fuse.DirEntry, len(nodes))
	for i, node := range nodes {
		info := node.Stat()
		if p.GetChild(info.Name) == nil {
			p.NewChild(info.Name, info.IsDir(), node)
		}
		info.FillDirEntry(&a[i])
	}
	return a, fuse.OK
}
