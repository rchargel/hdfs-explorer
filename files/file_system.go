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
	user, err := user.Current()
	if err == nil {
		return user.Username, nil
	}
	return "", err
}

// Connect opens a FileSystemClient.
func (r *FileSystem) Connect() (FileSystemClient, error) {
	user, err := r.user()
	if err == nil {
		if conn, nerr := rpc.NewNamenodeConnectionWithOptions(
			rpc.NamenodeConnectionOptions{
				Addresses: r.Addresses,
				User:      user,
			},
		); nerr == nil {
			client, cerr := hdfs.NewClient(
				hdfs.ClientOptions{
					Addresses: r.Addresses,
					Namenode:  conn,
					User:      user,
				},
			)
			if cerr == nil {
				return CreateHDFSFileSystemClient(r.Name, client), err
			}
			return nil, cerr
		} else {
			return nil, nerr
		}
	}
	return nil, err
}
