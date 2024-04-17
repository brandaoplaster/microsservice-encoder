package services

import (
	"errors"
	"os"
	"strconv"

	"github.com/brandaoplaster/encoder/application/repositories"
	"github.com/brandaoplaster/encoder/domain"
)

type JobService struct {
	Job           *domain.Job
	JobRepository repositories.JobRepository
	VideoService  VideoService
}

func (job *JobService) Start() error {
	erro := job.changeJobStatus("DOWNLOADING")
	if erro != nil {
		return job.failJob(erro)
	}

	err := job.VideoService.Download(os.Getenv("inputBucketName"))
	if err != nil {
		return job.failJob(err)
	}

	erro = job.changeJobStatus("FRAGMENTING")
	if erro != nil {
		return job.failJob(erro)
	}

	err = job.VideoService.Fragment()
	if err != nil {
		return job.failJob(err)
	}

	erro = job.changeJobStatus("ENCODING")
	if erro != nil {
		return job.failJob(erro)
	}

	err = job.VideoService.Encode()
	if err != nil {
		return job.failJob(err)
	}

	err = job.performUpload()

	if err != nil {
		return job.failJob(err)
	}

	err = job.changeJobStatus("FINISHING")

	if err != nil {
		return job.failJob(err)
	}

	err = job.VideoService.Finish()

	if err != nil {
		return job.failJob(err)
	}

	err = job.changeJobStatus("COMPLETED")

	if err != nil {
		return job.failJob(err)
	}

	return nil
}

func (job *JobService) performUpload() error {

	err := job.changeJobStatus("UPLOADING")

	if err != nil {
		return job.failJob(err)
	}

	videoUpload := NewVideoUpload()
	videoUpload.OutputBucket = os.Getenv("outputBucketName")
	videoUpload.VideoPath = os.Getenv("localStoragePath") + "/" + job.VideoService.Video.ID
	concurrency, _ := strconv.Atoi(os.Getenv("CONCURRENCY_UPLOAD"))
	doneUpload := make(chan string)

	go videoUpload.ProcessUpload(concurrency, doneUpload)

	var uploadResult string
	uploadResult = <-doneUpload

	if uploadResult != "upload completed" {
		return job.failJob(errors.New(uploadResult))
	}

	return err
}

func (job *JobService) changeJobStatus(status string) error {
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
