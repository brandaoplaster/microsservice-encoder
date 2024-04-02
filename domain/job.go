package domain

import (
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

type Job struct {
	ID           string    `json:"job_id" valid:"uuid" gorm:"type:uuid;primary_key"`
	OutputBucket string    `json:"output_bucket" valid:"notnull"`
	Status       string    `json:"status" valid:"notnull"`
	Video        *Video    `json:"video" valid:"-"`
	VideoID      string    `json:"-" valid:"-" gorm:"column:video_id;type:uuid;notnull"`
	Error        string    `valid:"-"`
	CreatedAt    time.Time `json:"created_at" valid:"-"`
	UpdatedAt    time.Time `json:"updated_at" valid:"-"`
}

func init() {
	govalidator.SetFieldsRequiredByDefault(true)
}

func (job *Job) Validate() error {
	_, erro := govalidator.ValidateStruct(job)

	if erro != nil {
		return erro
	}

	return nil
}

func (job *Job) prepare() {
	job.ID = uuid.NewV4().String()
	job.CreatedAt = time.Now()
	job.UpdatedAt = time.Now()
}

func NewJob(output string, video *Video, status string) (*Job, error) {
	job := Job{
		OutputBucket: output,
		Video:        video,
		Status:       status,
	}

	job.prepare()

	error := job.Validate()

	if error != nil {
		return nil, error
	}

	return &job, nil
}
