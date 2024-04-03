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

func TestJobRepositoryDbInsert(test *testing.T) {
	db := database.NewDatabaseTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: db}
	repo.Insert(video)

	job, err := domain.NewJob("output_path", video, "Pending")
	require.Nil(test, err)

	repoJob := repositories.JobRepositoryDb{Db: db}
	repoJob.Insert(job)

	findJob, err := repoJob.Find(job.ID)
	require.NotEmpty(test, findJob.ID)
	require.Nil(test, err)
	require.Equal(test, findJob.ID, job.ID)
	require.Equal(test, findJob.VideoID, video.ID)
}

func TestJobRepositoryDbUpdate(test *testing.T) {
	db := database.NewDatabaseTest()
	defer db.Close()

	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	repo := repositories.VideoRepositoryDb{Db: db}
	repo.Insert(video)

	job, err := domain.NewJob("output_path", video, "Pending")
	require.Nil(test, err)

	repoJob := repositories.JobRepositoryDb{Db: db}
	repoJob.Insert(job)

	job.Status = "Complete"

	repoJob.Update(job)

	findJob, err := repoJob.Find(job.ID)
	require.NotEmpty(test, findJob.ID)
	require.Nil(test, err)
	require.Equal(test, findJob.Status, job.Status)
}
