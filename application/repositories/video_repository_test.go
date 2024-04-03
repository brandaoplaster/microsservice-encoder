package repositories_test

import (
	"testing"
	"time"

	"github.com/brandaoplaster/encoder/application/repositories"
	"github.com/brandaoplaster/encoder/domain"
	"github.com/brandaoplaster/encoder/framework/database"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestVideoRepositoryDBInsert(test *testing.T) {
	db := database.NewDatabaseTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: db}
	repo.Insert(video)

	findVideo, err := repo.Find(video.ID)

	require.NotEmpty(test, findVideo.ID)
	require.Nil(test, err)
	require.Equal(test, video.ID, findVideo.ID)
}
