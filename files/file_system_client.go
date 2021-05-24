package files

import (
	"io"
	"os"

	"github.com/colinmarc/hdfs"
)

// FileSystemClient interface is an abstraction around requests to an hdfs file system.
type FileSystemClient interface {
	Close() error
	Upload(source, destination string) error
	Download(source, destination string) error
	Move(oldpath, newpath string) error
	MakeDir(path string, permission os.FileMode) error
	List(path string) ([]HdfsFileInfo, error)
	Delete(path string) error
	Read(path string) (*io.Reader, error)
	Write(path string, append bool) (*io.Writer, error)
	Status(path string) (HdfsFileInfo, error)
}

// CreateHDFSFileSystemClient wraps an hdfs.Client in order to produce the FileSystemClient interface.
func CreateHDFSFileSystemClient(client *hdfs.Client) FileSystemClient {
	return &hdfsFileSystemClient{client}
}

type hdfsFileSystemClient struct {
	client *hdfs.Client
}

func (f *hdfsFileSystemClient) Close() error {
	return f.client.Close()
}

func (f *hdfsFileSystemClient) Upload(source, destination string) error {
	return f.client.CopyToRemote(source, destination)
}

func (f *hdfsFileSystemClient) Download(source, destination string) error {
	return f.client.CopyToLocal(source, destination)
}

func (f *hdfsFileSystemClient) Move(oldpath, newpath string) error {
	return f.client.Rename(oldpath, newpath)
}

func (f *hdfsFileSystemClient) MakeDir(path string, permission os.FileMode) error {
	return f.client.MkdirAll(path, permission)
}

func (f *hdfsFileSystemClient) List(path string) ([]HdfsFileInfo, error) {
	files, err := f.client.ReadDir(path)
	if err == nil {
		hdfsFiles := make([]HdfsFileInfo, len(files))

		for idx, value := range files {
			hdfsFiles[idx] = &hdfsFileInfo{value}
		}
		return hdfsFiles, nil
	}
	return nil, err
}

func (f *hdfsFileSystemClient) Delete(path string) error {
	return f.client.Remove(path)
}

func (f *hdfsFileSystemClient) Read(path string) (*io.Reader, error) {
	return nil, nil
}

func (f *hdfsFileSystemClient) Write(path string, append bool) (*io.Writer, error) {
	if !append {
		f.Delete(path)
	}
	return nil, nil
}

func (f *hdfsFileSystemClient) Status(path string) (HdfsFileInfo, error) {
	file, err := f.client.Stat(path)
	return &hdfsFileInfo{file}, err
}
