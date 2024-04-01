package domain_test

import (
	"testing"
	"time"

	"github.com/brandaoplaster/encoder/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestNewJob(test *testing.T) {
	video := domain.NewVideo()
	video.ID = uuid.NewV4().String()
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	job, erro := domain.NewJob("output", video, "pending")

	require.NotNil(test, job)
	require.Nil(test, erro)
}
