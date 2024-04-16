package services

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
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

func (upload *VideoUpload) LoadPaths() error {
	erro := filepath.Walk(upload.VideoPath, func(path string, info os.FileInfo, erro error) error {
		if !info.IsDir() {
			upload.Paths = append(upload.Paths, path)
		}
		return nil
	})

	if erro != nil {
		return erro
	}
	return nil
}

func getClientUpload() (*storage.Client, context.Context, error) {
	context := context.Background()

	client, erro := storage.NewClient(context)

	if erro != nil {
		return nil, nil, erro
	}

	return client, context, nil
}

func (upload *VideoUpload) ProcessUpload(concurrency int, doneUpload chan string) error {
	in := make(chan int, runtime.NumCPU())
	returnChannel := make(chan string)

	erro := upload.LoadPaths()
	if erro != nil {
		return erro
	}

	uploadClient, context, erro := getClientUpload()
	if erro != nil {
		return erro
	}

	for process := 0; process < concurrency; process++ {
		go upload.uploadWorker(in, returnChannel, uploadClient, context)
	}

	go func() {
		for x := 0; x < len(upload.Paths); x++ {
			in <- x
		}
		close(in)
	}()

	for r := range returnChannel {
		if r != "" {
			doneUpload <- r
			break
		}
	}

	return nil
}

func (upload *VideoUpload) uploadWorker(int chan int, returnChan chan string, uploadClient *storage.Client, context context.Context) {
	for x := range int {
		erro := upload.UploadObject(upload.Paths[x], uploadClient, context)
		if erro != nil {
			upload.Errors = append(upload.Errors, upload.Paths[x])
			log.Printf("Error during the upload: %v. Error: %v", upload.Paths[x], erro)
			returnChan <- erro.Error()
		}
		returnChan <- ""
	}
	returnChan <- "Upload completed"
}
