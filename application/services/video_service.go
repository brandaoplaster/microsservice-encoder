package services

import (
	"context"
	"io/ioutil"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"github.com/brandaoplaster/encoder/application/repositories"
	"github.com/brandaoplaster/encoder/domain"
)

type VideoService struct {
	Video           *domain.Video
	VideoRepository repositories.VideoRepository
}

func NewVideoService() VideoService {
	return VideoService{}
}

func (video *VideoService) Download(bucketName string) error {
	context := context.Background()
	client, erro := storage.NewClient(context)

	if erro != nil {
		return erro
	}

	bucket := client.Bucket(bucketName)
	object := bucket.Object(video.Video.FilePath)

	read, erro := object.NewReader(context)
	if erro != nil {
		return erro
	}

	defer read.Close()

	body, erro := ioutil.ReadAll(read)
	if erro != nil {
		return erro
	}

	file, erro := os.Create(os.Getenv("localStoragePath") + "/" + video.Video.ID + ".mp4")
	if erro != nil {
		return erro
	}

	_, erro = file.Write(body)
	if erro != nil {
		return erro
	}

	defer file.Close()

	log.Printf("Video %v has been downloaded", video.Video.ID)

	return nil
}
