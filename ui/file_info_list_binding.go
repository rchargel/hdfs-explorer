package ui

import (
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"strings"
	"time"

	"fyne.io/fyne/v2/data/binding"
	"github.com/rchargel/hdfs-explorer/files"
)

type HdfsFileInfoList interface {
	binding.DataList

	Get() ([]files.HdfsFileInfo, error)
	GetValue(int) (files.HdfsFileInfo, error)
	UpdatePath(string) error
	Reload() error
	Close()
}

func NewHdfsFileInfoList(client files.FileSystemClient) HdfsFileInfoList {
	return &boundFileInfoList{binding.BindStringList(&[]string{}), client}
}

type boundFileInfoList struct {
	list   binding.ExternalStringList
	client files.FileSystemClient
}

func (f *boundFileInfoList) AddListener(l binding.DataListener)      { f.list.AddListener(l) }
func (f *boundFileInfoList) RemoveListener(l binding.DataListener)   { f.list.RemoveListener(l) }
func (f *boundFileInfoList) Length() int                             { return f.list.Length() }
func (f *boundFileInfoList) GetItem(i int) (binding.DataItem, error) { return f.list.GetItem(i) }
func (f *boundFileInfoList) Reload() error                           { return f.list.Reload() }
func (f *boundFileInfoList) Close()                                  { f.client.Close() }
func (f *boundFileInfoList) Get() ([]files.HdfsFileInfo, error) {
	a := make([]files.HdfsFileInfo, f.Length())
	for idx := 0; idx < f.Length(); idx++ {
		i, err := f.GetValue(idx)
		if err != nil {
			return a, err
		}
		a[idx] = i
	}

	return a, nil
}
func (f *boundFileInfoList) GetValue(i int) (files.HdfsFileInfo, error) {
	s, err := f.list.GetValue(i)
	if err == nil {
		return fromString(s)
	}
	return nil, err
}
func (f *boundFileInfoList) UpdatePath(path string) error {
	a := make([]string, 0)
	if path != "/" {
		p := files.Parent(path)
		pInfo, err := f.client.Status(p)
		if err != nil {
			return err
		}
		w := wrap(pInfo)
		w.name = ".."
		a = append(a, toString(w))
	}
	list, err := f.client.List(path)
	if err != nil {
		return err
	}
	for _, info := range list {
		a = append(a, toString(info))
	}
	return f.list.Set(a)
}
func (f *boundFileInfoList) set(infos []files.HdfsFileInfo) error {
	a := make([]string, len(infos))
	for i, info := range infos {
		j := toString(info)
		a[i] = j
		println(j)
	}

	return f.list.Set(a)
}

type fileInfoHolder struct {
	isDir       bool
	modTime     time.Time
	accessTime  time.Time
	mode        fs.FileMode
	name        string
	size        int64
	owner       string
	ownerGroup  string
	replication uint32
	blockSize   uint64
}

func (f *fileInfoHolder) IsDir() bool           { return f.isDir }
func (f *fileInfoHolder) ModTime() time.Time    { return f.modTime }
func (f *fileInfoHolder) AccessTime() time.Time { return f.accessTime }
func (f *fileInfoHolder) Mode() fs.FileMode     { return f.mode }
func (f *fileInfoHolder) Name() string          { return f.name }
func (f *fileInfoHolder) Size() int64           { return f.size }
func (f *fileInfoHolder) Owner() string         { return f.owner }
func (f *fileInfoHolder) OwnerGroup() string    { return f.ownerGroup }
func (f *fileInfoHolder) Replication() uint32   { return f.replication }
func (f *fileInfoHolder) BlockSize() uint64     { return f.blockSize }
func toString(f files.HdfsFileInfo) string {
	return fmt.Sprintf(
		"%v|%d|%d|0%v|%v|%d|%v|%v|%d|%d",
		f.IsDir(),
		f.ModTime().Unix(),
		f.AccessTime().Unix(),
		strconv.FormatUint(uint64(f.Mode()), 8),
		f.Name(),
		f.Size(),
		f.Owner(),
		f.OwnerGroup(),
		f.Replication(),
		f.BlockSize(),
	)
}
func fromString(j string) (files.HdfsFileInfo, error) {
	a := strings.Split(j, "|")
	if len(a) != 10 {
		return nil, fmt.Errorf("Invalid value length %d", len(a))
	}
	isDir, _ := strconv.ParseBool(a[0])
	mtu, _ := strconv.ParseInt(a[1], 10, 64)
	mt := time.Unix(mtu, 0)
	atu, _ := strconv.ParseInt(a[2], 10, 64)
	at := time.Unix(atu, 0)
	fm, _ := strconv.ParseUint(a[3], 0, 32)
	sz, _ := strconv.ParseInt(a[5], 10, 64)
	rp, _ := strconv.ParseUint(a[8], 10, 32)
	bz, _ := strconv.ParseUint(a[9], 10, 64)

	return &fileInfoHolder{
		isDir:       isDir,
		modTime:     mt,
		accessTime:  at,
		mode:        os.FileMode(fm),
		name:        a[4],
		size:        sz,
		owner:       a[6],
		ownerGroup:  a[7],
		replication: uint32(rp),
		blockSize:   bz,
	}, nil
}
func wrap(f files.HdfsFileInfo) *fileInfoHolder {
	return &fileInfoHolder{
		f.IsDir(),
		f.ModTime(),
		f.AccessTime(),
		f.Mode(),
		f.Name(),
		f.Size(),
		f.Owner(),
		f.OwnerGroup(),
		f.Replication(),
		f.BlockSize(),
	}
}
