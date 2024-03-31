package domain_test

import (
	"testing"
	"time"

	"github.com/brandaoplaster/encoder/domain"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
)

func TestValidateIfVideoIsEmpty(test *testing.T) {
	video := domain.NewVideo()
	error := video.Validate()

	require.Error(test, error)
}

func TestVideoIdIsNotAUUID(test *testing.T) {
	video := domain.NewVideo()

	video.ID = "abc"
	video.ResourceID = "a"
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	error := video.Validate()

	require.Error(test, error)
}

func TestVideoValidation(test *testing.T) {
	video := domain.NewVideo()

	video.ID = uuid.NewV4().String()
	video.ResourceID = "a"
	video.FilePath = "path"
	video.CreatedAt = time.Now()

	err := video.Validate()
	require.Nil(test, err)
}
