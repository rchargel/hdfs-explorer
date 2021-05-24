package files

import (
	"os/user"

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
}

func (r *FileSystem) user() (string, error) {
	if user, err := user.Current(); err == nil {
		return user.Username, nil
	} else {
		return "", err
	}
}

// Connect opens a FileSystemClient.
func (r *FileSystem) Connect() (FileSystemClient, error) {
	if user, err := r.user(); err == nil {
		if conn, nerr := rpc.NewNamenodeConnectionWithOptions(
			rpc.NamenodeConnectionOptions{
				Addresses: r.Addresses,
				User:      user,
			},
		); nerr == nil {
			if client, cerr := hdfs.NewClient(
				hdfs.ClientOptions{
					Addresses: r.Addresses,
					Namenode:  conn,
					User:      user,
				},
			); cerr == nil {
				return CreateHDFSFileSystemClient(client), err
			} else {
				return nil, cerr
			}
		} else {
			return nil, nerr
		}
	} else {
		return nil, err
	}
}
