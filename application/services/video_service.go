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

func (video *VideoService) Encode() error {
	args := []string{}
	args = append(args, os.Getenv("localStoragePath")+"/"+video.Video.ID+".frag")
	args = append(args, "--use-segment-timeline")
	args = append(args, "-o")
	args = append(args, os.Getenv("localStoragePath")+"/"+video.Video.ID)
	args = append(args, "-f")
	args = append(args, "--exec-dir")
	args = append(args, "/opt/bento4/bin")

	comand := exec.Command("mp4dash", args...)

	output, erro := comand.CombinedOutput()

	if erro != nil {
		return erro
	}

	printOutput(output)
	return nil
}

func (video *VideoService) Finish() error {
	erro := os.Remove(os.Getenv("localStoragePath") + "/" + video.Video.ID + ".mp4")
	if erro != nil {
		log.Println("Error deleting video", video.Video.ID+".mp4")
		return erro
	}

	erro = os.Remove(os.Getenv("localStoragePath") + "/" + video.Video.ID + ".frag")
	if erro != nil {
		log.Println("Error deleting video", video.Video.ID+".frag")
		return erro
	}

	erro = os.RemoveAll(os.Getenv("localStoragePath") + "/" + video.Video.ID)
	if erro != nil {
		log.Println("Error deleting directory", video.Video.ID)
		return erro
	}

	log.Println("Files have been deleted", video.Video.ID)
	return nil
}

func printOutput(output []byte) {
	if len(output) > 0 {
		log.Printf("Output: %s\n", string(output))
	}
}
