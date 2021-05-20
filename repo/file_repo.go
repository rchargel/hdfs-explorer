package repo

import (
	"io"
	"os"

	"github.com/colinmarc/hdfs"
)

type FileRepo interface {
	Close() error
	Upload(source, destination string) error
	Download(source, destination string) error
	Move(oldpath, newpath string) error
	MakeDir(path string, permission os.FileMode) error
	List(path string) ([]os.FileInfo, error)
	Delete(path string) error
	Read(path string) (*io.Reader, error)
	Write(path string, append bool) (*io.Writer, error)
}

func CreateHDFSFileRepo(client *hdfs.Client) FileRepo {
	return HdfsFileRepo{client}
}

type HdfsFileRepo struct {
	client *hdfs.Client
}

func (f HdfsFileRepo) Close() error {
	return f.client.Close()
}

func (f HdfsFileRepo) Upload(source, destination string) error {
	return f.client.CopyToRemote(source, destination)
}

func (f HdfsFileRepo) Download(source, destination string) error {
	return f.client.CopyToLocal(source, destination)
}

func (f HdfsFileRepo) Move(oldpath, newpath string) error {
	return f.client.Rename(oldpath, newpath)
}

func (f HdfsFileRepo) MakeDir(path string, permission os.FileMode) error {
	return f.client.MkdirAll(path, permission)
}

func (f HdfsFileRepo) List(path string) ([]os.FileInfo, error) {
	return f.client.ReadDir(path)
}

func (f HdfsFileRepo) Delete(path string) error {
	return f.client.Remove(path)
}

func (f HdfsFileRepo) Read(path string) (*io.Reader, error) {
	return nil, nil
}

func (f HdfsFileRepo) Write(path string, append bool) (*io.Writer, error) {
	return nil, nil
}
