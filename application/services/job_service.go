package services

import (
	"github.com/brandaoplaster/encoder/application/repositories"
	"github.com/brandaoplaster/encoder/domain"
)

type JobService struct {
	Job           *domain.Job
	JobRepository repositories.JobRepository
	VideoService  VideoService
}

func (job *JobService) Start() error {

}

func (job *JobService) ChangeJobStatus(status string) error {
	var erro error

	job.Job.Status = status

	job.Job, erro = job.JobRepository.Update(job.Job)
	if erro != nil {
		return job.failJob(erro)
	}

	return nil
}

func (job *JobService) failJob(error error) error {
	job.Job.Status = "FAILED"
	job.Job.Error = error.Error()

	_, err := job.JobRepository.Update(job.Job)
	if err != nil {
		return err
	}

	return error
}
