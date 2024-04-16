package services

import (
	"context"
	"io"
	"os"
	"strings"

	"cloud.google.com/go/storage"
)

type VideoUpload struct {
	Paths        []string
	VideoPath    string
	OutputBucket string
	Errors       []string
}

func NewVideoUpload() *VideoUpload {
	return &VideoUpload{}
}

func (upload *VideoUpload) UploadObject(objectPath string, client *storage.Client, context context.Context) error {
	path := strings.Split(objectPath, os.Getenv("localStoragePath")+"/")

	file, erro := os.Open(objectPath)

	if erro != nil {
		return erro
	}

	defer file.Close()

	writeBucket := client.Bucket(upload.OutputBucket).Object(path[1]).NewWriter(context)

	writeBucket.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	if _, erro = io.Copy(writeBucket, file); erro != nil {
		return erro
	}

	if erro = writeBucket.Close(); erro != nil {
		return erro
	}

	return nil
}
