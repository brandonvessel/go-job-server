package client

import (
	"bytes"
	"encoding/json"
	"net/http"

	libjob "github.com/brandonvessel/go-job-server/lib/job"
)

// submitWork submits work to the server
func (client *Client) submitWork(job *libjob.Job) error {
	// encode job
	buf, err := json.Marshal(job)

	if err != nil {
		return err
	}

	// create request
	req, err := http.NewRequest("POST", "http://"+client.ServerIP+":"+client.ServerPort+"/job", bytes.NewBuffer(buf))

	if err != nil {
		return err
	}

	// send request
	resp, err := client.Client.Do(req)

	if err != nil {
		return err
	}

	// decode job
	var djob libjob.Job

	err = json.NewDecoder(resp.Body).Decode(&djob)

	if err != nil {
		return err
	}

	// add job to cache
	client.jobCache[job.UUID] = *job

	return nil
}

// GetStatus gets the status of a job associated with the given UUID
func (client *Client) getStatusByUUID(uuid string) (string, error) {
	// create request
	req, err := http.NewRequest("GET", "http://"+client.ServerIP+":"+client.ServerPort+"/job/"+uuid, nil)

	if err != nil {
		return "", err
	}

	// send request
	resp, err := client.Client.Do(req)

	if err != nil {
		return "", err
	}

	// decode json response
	var status string

	err = json.NewDecoder(resp.Body).Decode(&status)

	if err != nil {
		return "", err
	}

	return status, nil
}

// getStatus gets the status of a job
func (client *Client) getStatus(job *libjob.Job) (string, error) {
	return client.getStatusByUUID(job.UUID)
}

// deleteJob deletes a job from the server and removes it from the jobCache
func (client *Client) deleteJob(job *libjob.Job) {
	// create request
	req, err := http.NewRequest("DELETE", "http://"+client.ServerIP+":"+client.ServerPort+"/job/"+job.UUID, nil)

	if err != nil {
		return
	}

	// send request
	resp, err := client.Client.Do(req)

	if err != nil {
		return
	}

	// close response body
	resp.Body.Close()

	// delete job from cache
	delete(client.jobCache, job.UUID)
}

// getResult gets the result of a job
func (client *Client) getResult(job libjob.Job) {
	// create request
	req, err := http.NewRequest("GET", "http://"+client.ServerIP+":"+client.ServerPort+"/result/"+job.UUID, nil)

	if err != nil {
		return
	}

	// send request
	resp, err := client.Client.Do(req)

	if err != nil {
		return
	}
	defer resp.Body.Close()

	// get result string from response
	var mjob libjob.Job

	err = json.NewDecoder(resp.Body).Decode(&mjob)

	if err != nil {
		return
	}

	// set result in job cache
	cacheJob := client.jobCache[job.UUID]
	cacheJob.Result = mjob.Result
	cacheJob.Status = mjob.Status
	cacheJob.Error = mjob.Error

	client.jobCache[job.UUID] = cacheJob
}
