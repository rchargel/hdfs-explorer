package files

import (
	"io/fs"
	"os"
	"time"

	"github.com/colinmarc/hdfs/protocol/hadoop_hdfs"
)

// HdfsFileInfo wrapper around os.FileInfo with some additional functions
// related to HDFS files.
type HdfsFileInfo interface {
	IsDir() bool
	ModTime() time.Time
	AccessTime() time.Time
	Mode() fs.FileMode
	Name() string
	Size() int64
	Owner() string
	OwnerGroup() string
	Replication() uint32
	BlockSize() uint64
}

type hdfsFileInfo struct {
	fileInfo os.FileInfo
}

func (h *hdfsFileInfo) IsDir() bool {
	return h.fileInfo.IsDir()
}

func (h *hdfsFileInfo) ModTime() time.Time {
	return h.fileInfo.ModTime()
}

func (h *hdfsFileInfo) AccessTime() time.Time {
	return time.Unix(int64(h.getSys().GetAccessTime())/1000, 0)
}

func (h *hdfsFileInfo) Mode() fs.FileMode {
	return h.fileInfo.Mode()
}

func (h *hdfsFileInfo) Name() string {
	return h.fileInfo.Name()
}

func (h *hdfsFileInfo) Size() int64 {
	return h.fileInfo.Size()
}

func (h *hdfsFileInfo) Owner() string {
	return h.getSys().GetOwner()
}

func (h *hdfsFileInfo) OwnerGroup() string {
	return h.getSys().GetGroup()
}

func (h *hdfsFileInfo) Replication() uint32 {
	return h.getSys().GetBlockReplication()
}

func (h *hdfsFileInfo) BlockSize() uint64 {
	return h.getSys().GetBlocksize()
}

func (h *hdfsFileInfo) getSys() *hadoop_hdfs.HdfsFileStatusProto {
	return h.fileInfo.Sys().(*hadoop_hdfs.HdfsFileStatusProto)
}
