package worker

// this package is used to create a worker than can process work and send the result back to the server

import (
	"encoding/json"
	"log"
	"net/http"

	libjob "github.com/brandonvessel/go-job-server/lib/job"
)

type Worker struct {
	// server ip
	ServerIP string

	// server port
	ServerPort string

	// current job
	job libjob.Job

	// task ids the worker is capable of processing
	taskIds []string
}

type Task struct {
	UUID       string
	Parameters []string

	// task result
	Result string
	// task error
	Error string
	// task status
	Status string
}

// TaskFunc defines the function for a task
type TaskFunc func(task Task) (Task, error)

// ProcessWork processes work and returns the result
func (worker *Worker) ProcessWork(taskFunc TaskFunc) (Task, error) {
	// get work
	job, err := worker.GetWork()

	// if error
	if err != nil {
		// log error
		log.Println(err)
		return Task{}, nil
	} else {
		// create task
		task := Task{
			UUID:       job.UUID,
			Parameters: job.Data,
		}

		// process task
		task, err = taskFunc(task)

		// if error
		if err != nil {
			// log error
			log.Println(err)
			return Task{}, nil
		} else {
			// return task
			return task, nil
		}
	}
}

// parseWork parses the work and returns the job information
func (worker *Worker) parseWork(resp *http.Response) (libjob.Job, error) {
	// get job
	var job libjob.Job

	// decode job
	err := json.NewDecoder(resp.Body).Decode(&job)

	// if error
	if err != nil {
		// log error
		log.Println(err)
		return libjob.Job{}, nil
	} else {
		// return job
		return job, nil
	}
}

// get work from server and return as map
func (worker *Worker) GetWork() (libjob.Job, error) {
	// get work from server with get request
	resp, err := http.Get("http://" + worker.ServerIP + ":" + worker.ServerPort + "/work")

	// if error
	if err != nil {
		// log error
		log.Println(err)
		return libjob.Job{}, nil
	} else {
		// return job information
		return worker.parseWork(resp)
	}

}
