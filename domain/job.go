package domain

import (
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

type Job struct {
	ID           string    `valid:"uuid"`
	OutputBucket string    `valid:"notnull"`
	Status       string    `valid:"notnull"`
	Video        *Video    `valid:"-"`
	Error        string    `valid:"-"`
	CreatedAt    time.Time `valid:"-"`
	UpdatedAt    time.Time `valid:"-"`
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
