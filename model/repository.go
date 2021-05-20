package model

import (
	"github.com/colinmarc/hdfs"
	"github.com/rchargel/hdfs-explorer/repo"
)

type Repository struct {
	Id          uint16
	Name        string
	Description string
	location    string
}

func (r *Repository) Connect() (repo.FileRepo, error) {
	client, err := hdfs.New(r.location)
	if client == nil {
		return nil, err
	}
	return repo.CreateHDFSFileRepo(client), err
}
