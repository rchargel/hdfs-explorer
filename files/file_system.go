package files

import (
	"github.com/colinmarc/hdfs"
	"github.com/colinmarc/hdfs/rpc"
)

// FileSystem defines a simple pointer to an HDFS file system.
type FileSystem struct {

	// Name is a unique identifier for the file system.
	Name string

	// Description should be a brief note.
	Description string

	// Addresses are a list of the name nodes for the file system.
	Addresses []string

	// User with access to the name nodes.
	User string
}

// Connect opens a FileSystemClient.
func (r *FileSystem) Connect() (FileSystemClient, error) {
	conn, nerr := rpc.NewNamenodeConnectionWithOptions(
		rpc.NamenodeConnectionOptions{
			Addresses: r.Addresses,
			User:      r.User,
		},
	)

	if nerr != nil {
		return nil, nerr
	}

	client, err := hdfs.NewClient(
		hdfs.ClientOptions{
			Addresses: r.Addresses,
			Namenode:  conn,
			User:      r.User,
		},
	)
	if client == nil {
		return nil, err
	}
	return CreateHDFSFileSystemClient(client), err
}
