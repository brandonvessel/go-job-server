package client

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	libjob "github.com/brandonvessel/go-job-server/lib/job"
)

// Client is a client that can send work to the server
type Client struct {
	// server ip
	ServerIP string

	// server port
	ServerPort string

	// http client
	Client *http.Client

	// amount of time in milliseconds to wait between job status requests
	StatusRefreshFrequency int `default:"1000"`

	// job cache
	jobCache map[string]libjob.Job
}

// NewClient creates a new client and returns it
func NewClient(serverIP string, serverPort string) *Client {
	return &Client{
		ServerIP:   serverIP,
		ServerPort: serverPort,
		Client:     &http.Client{},
		jobCache:   make(map[string]libjob.Job),
	}
}

// NewJob sends a new job to the server
func (client *Client) NewJob(name string, arguments []string) (string, error) {
	// create job
	job := libjob.Job{
		Name: name,
		Data: arguments,
	}

	// submit job
	err := client.submitWork(&job)

	// if error
	if err != nil {
		return "", err
	}

	// if uuid empty
	if job.UUID == "" {
		return "", errors.New("empty uuid")
	}

	fmt.Println("JOB UUID", job.UUID)

	return job.UUID, nil
}

// IsFinished returns true if the job is finished and false if it is not finished. Errors if job has an error
func (client *Client) IsFinished(uuid string) (bool, error) {
	// get status
	status, err := client.getStatusByUUID(uuid)

	// if error
	if err != nil {
		return false, err
	}

	// return values based on status
	switch status {
	case "done":
		// if status is done
		return true, nil
	case "error":
		// if status is error
		return true, nil
	case "queued":
		// if status is queued
		return false, nil
	default:
		// if status is unknown
		return false, nil
	}
}

// GetResult returns the result of the job.
func (client *Client) GetResult(jobUUID string) (string, error) {
	// TODO: figure out if I need this
	//defer func(){
	// remove job from jobcache
	//	delete(client.jobCache, jobUUID)
	//}()

	// get job from jobcache
	job, exists := client.jobCache[jobUUID]

	// if job not in cache
	if !exists {
		// return error
		return "", errors.New("Invalid uuid")
	}

	// if is finished
	if job.Endtime != 0 {
		// return job result
		return job.Result, nil
	} else {
		// wait for result
		for {
			// check status
			isFinished, err := client.IsFinished(jobUUID)
			if err != nil {
				// TODO: Check for networking errors
			}

			// if finished
			if isFinished {
				break
			}

			// wait between requests
			time.Sleep(1 * time.Second)
		}

		// get result and return it
		result, err := job.GetResult()
		if err != nil {
			// TODO: Check for networking errors
		}

		// return based on job status
		switch job.Status {
		case "error":
			return result, errors.New(job.Error)
		case "done":
			return result, nil
		default:
			return result, nil
		}
	}
}
