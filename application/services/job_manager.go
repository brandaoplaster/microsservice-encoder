package services

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/brandaoplaster/encoder/application/repositories"
	"github.com/brandaoplaster/encoder/domain"
	"github.com/brandaoplaster/encoder/framework/queue"
	"github.com/jinzhu/gorm"
	"github.com/streadway/amqp"
)

type JobManager struct {
	Db               *gorm.DB
	Domain           domain.Job
	MessageChannel   chan amqp.Delivery
	JobReturnChannel chan JobWorkerResult
	RabbitMQ         *queue.RabbitMQ
}

type JobNotificationError struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

func NewJobManager(db *gorm.DB, rabbitMQ *queue.RabbitMQ, jobReturnChannel chan JobWorkerResult, messageChannel chan amqp.Delivery) *JobManager {
	return &JobManager{
		Db:               db,
		Domain:           domain.Job{},
		MessageChannel:   messageChannel,
		JobReturnChannel: jobReturnChannel,
		RabbitMQ:         rabbitMQ,
	}
}

func (job *JobManager) Start(ch *amqp.Channel) {

	videoService := NewVideoService()
	videoService.VideoRepository = repositories.VideoRepositoryDb{Db: job.Db}

	jobService := JobService{
		JobRepository: repositories.JobRepositoryDb{Db: job.Db},
		VideoService:  videoService,
	}

	concurrency, err := strconv.Atoi(os.Getenv("CONCURRENCY_WORKERS"))

	if err != nil {
		log.Fatalf("error loading var: CONCURRENCY_WORKERS.")
	}

	for qtdProcesses := 0; qtdProcesses < concurrency; qtdProcesses++ {
		go JobWorker(job.MessageChannel, job.JobReturnChannel, jobService, job.Domain, qtdProcesses)
	}

	for jobResult := range job.JobReturnChannel {
		if jobResult.Error != nil {
			err = job.checkParseErrors(jobResult)
		} else {
			err = job.notifySuccess(jobResult, ch)
		}

		if err != nil {
			jobResult.Message.Reject(false)
		}
	}
}

func (job *JobManager) checkParseErrors(jobResult JobWorkerResult) error {
	if jobResult.Job.ID != "" {
		log.Printf("MessageID: %v. Error during the job: %v with video: %v. Error: %v",
			jobResult.Message.DeliveryTag, jobResult.Job.ID, jobResult.Job.Video.ID, jobResult.Error.Error())
	} else {
		log.Printf("MessageID: %v. Error parsing message: %v", jobResult.Message.DeliveryTag, jobResult.Error)
	}

	errorMsg := JobNotificationError{
		Message: string(jobResult.Message.Body),
		Error:   jobResult.Error.Error(),
	}

	jobJson, err := json.Marshal(errorMsg)

	if err != nil {
		return err
	}

	err = job.notify(jobJson)

	if err != nil {
		return err
	}

	err = jobResult.Message.Reject(false)

	if err != nil {
		return err
	}

	return nil
}

func (job *JobManager) notifySuccess(jobResult JobWorkerResult, _ *amqp.Channel) error {

	Mutex.Lock()
	jobJson, err := json.Marshal(jobResult.Job)
	Mutex.Unlock()

	if err != nil {
		return err
	}

	err = job.notify(jobJson)

	if err != nil {
		return err
	}

	err = jobResult.Message.Ack(false)

	if err != nil {
		return err
	}

	return nil
}

func (job *JobManager) notify(jobJson []byte) error {

	err := job.RabbitMQ.Notify(
		string(jobJson),
		"application/json",
		os.Getenv("RABBITMQ_NOTIFICATION_EX"),
		os.Getenv("RABBITMQ_NOTIFICATION_ROUTING_KEY"),
	)

	if err != nil {
		return err
	}

	return nil
}
