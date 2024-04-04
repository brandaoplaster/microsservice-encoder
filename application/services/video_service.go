package services

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"os/exec"

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

func (video *VideoService) Fragment() error {
	localStoragePath := os.Getenv("localStoragePath")
	erro := os.Mkdir(localStoragePath+"/"+video.Video.ID, os.ModePerm)

	if erro != nil {
		return erro
	}

	source := localStoragePath + "/" + video.Video.ID + ".mp4"
	target := localStoragePath + "/" + video.Video.ID + ".frag"

	comand := exec.Command("mp4fragment", source, target)
	output, erro := comand.CombinedOutput()

	if erro != nil {
		return erro
	}

	printOutput(output)

	return nil
}

func printOutput(output []byte) {
	if len(output) > 0 {
		log.Printf("Output: %s\n", string(output))
	}
}
