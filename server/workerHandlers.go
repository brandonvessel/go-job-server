package server

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	libjob "github.com/brandonvessel/go-job-server/lib/job"
	"github.com/gin-gonic/gin"
)

// postWorkRequest processes a finished job and puts the results in the result map if the job exists
func (server *JobServer) postWorkRequest(c *gin.Context) {
	// get job
	var job libjob.Job

	// decode job
	err := json.NewDecoder(c.Request.Body).Decode(&job)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
	}

	// ensure job exists in map
	mjob, ok := server.workMap[job.UUID]

	if !ok {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Job not found"})
	}

	// if errored, set error
	if job.Error != "" {
		mjob.Error = job.Error
		mjob.Status = "error"
	} else {
		mjob.Result = job.Result
		mjob.Status = "done"
	}

	// set endtime
	mjob.Endtime = int(time.Now().Unix())

	// set job
	server.workMap[job.UUID] = mjob

	// return ok
	c.JSON(http.StatusOK, "accepted")
}

// getWorkRequest gets a job from the queue and returns it to the worker
func (server *JobServer) getWorkRequest(c *gin.Context) {
	// get job
	job, err := server.getJob()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
	}

	// add job to map
	server.workMap[job.UUID] = job

	// return job
	buf, err := json.Marshal(job)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, buf)
	}
}

// getJob gets a job from the queue
func (server *JobServer) getJob() (libjob.Job, error) {
	// get job from queue if queue is not empty
	// if empty, throw error
	if server.workQueue.Empty() {
		return libjob.Job{}, errors.New("queue empty")
	}

	// get job
	job, err := server.workQueue.Dequeue()

	if err != nil {
		return libjob.Job{}, err
	}

	// return job
	return job.(libjob.Job), nil
}