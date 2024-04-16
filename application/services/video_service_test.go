package services_test

import (
	"log"
	"testing"
	"time"

	"github.com/brandaoplaster/encoder/application/repositories"
	"github.com/brandaoplaster/encoder/application/services"
	"github.com/brandaoplaster/encoder/domain"
	"github.com/brandaoplaster/encoder/framework/database"
	"github.com/joho/godotenv"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func prepare() (*domain.Video, repositories.VideoRepositoryDb) {
	db := database.NewDatabaseTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "convite.mp4"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: db}

	return video, repo
}

func TestVideoServiceDownload(test *testing.T) {

	video, repo := prepare()

	videoService := services.NewVideoService()
	videoService.Video = video
	videoService.VideoRepository = repo

	err := videoService.Download("bucketTest")
	require.Nil(test, err)

	err = videoService.Fragment()
	require.Nil(test, err)

	err = videoService.Encode()
	require.Nil(test, err)

	err = videoService.Finish()
	require.Nil(test, err)
}
