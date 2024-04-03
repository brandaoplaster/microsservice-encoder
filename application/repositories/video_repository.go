package repositories

import (
	"fmt"

	"github.com/brandaoplaster/encoder/domain"
	"github.com/jinzhu/gorm"
	uuid "github.com/satori/go.uuid"
)

type VideoRepository interface {
	insert(video *domain.Video) (*domain.Video, error)
	find(id string) (*domain.Video, error)
}

type VideoRepositoryDb struct {
	Db *gorm.DB
}

func NewVideoRepository(db *gorm.DB) *VideoRepositoryDb {
	return &VideoRepositoryDb{Db: db}
}

func (repo VideoRepositoryDb) Insert(video *domain.Video) (*domain.Video, error) {
	if video.ID == "" {
		video.ID = uuid.NewV4().String()
	}

	erro := repo.Db.Create(video).Error

	if erro != nil {
		return nil, erro
	}

	return video, nil
}

func (repo VideoRepositoryDb) Find(id string) (*domain.Video, error) {
	var video domain.Video

	repo.Db.First(&video, "id = ?", id)

	if video.ID == "" {
		return nil, fmt.Errorf("video does not exist")
	}

	return &video, nil
}
