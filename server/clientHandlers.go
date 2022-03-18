package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	libjob "github.com/brandonvessel/go-job-server/lib/job"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// addJobRequest adds jobs to the queue
func (server *JobServer) addJobRequest(c *gin.Context) {
	// get job data
	data := strings.Split(c.PostForm("data"), "|")

	// create a new job
	job := libjob.Job{
		UUID:      uuid.New().String(),
		Name:      c.PostForm("name"),
		Data:      data,
		Status:    "queued",
		Starttime: int(time.Now().Unix()),
		Endtime:   0,
		Result:    "",
		Error:     "",
	}

	fmt.Println(job)

	// return job info
	c.JSON(http.StatusOK, job)
}

// deleteJobRequest deletes a job from the result map if it exists
func (server *JobServer) deleteJobRequest(c *gin.Context) {
	// get job uuid
	uuid := c.Param("uuid")

	// get job
	_, ok := server.workMap[uuid]

	// if job exists
	if ok {
		// delete job
		delete(server.workMap, uuid)

		// return success
		c.Status(http.StatusOK)
	} else {
		// return error
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Job not found",
		})
	}
}

// getResultRequest gets the result of a job and returns it to the client
func (server *JobServer) getResultRequest(c *gin.Context) {
	// get job uuid
	uuid := c.Param("uuid")

	// get job
	job, ok := server.workMap[uuid]

	// if job exists
	if ok {
		// if job is done
		if job.Status == "done" {
			// return job result
			c.JSON(http.StatusOK, gin.H{
				"result":  job.Result,
				"time":    job.Starttime,
				"endtime": job.Endtime,
				"error":   job.Error,
			})
		} else {
			// tell client job is not done
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"error": "Job not done",
			})
		}
	} else {
		// return error
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Job not found",
		})
	}
}

// getStatusRequest gets the status of a job and returns it to the client
func (server *JobServer) getStatusRequest(c *gin.Context) {
	// get job uuid
	uuid := c.Param("uuid")

	// get job
	job, ok := server.workMap[uuid]

	// if job exists
	if ok {
		// send json of status, result, and error
		c.JSON(http.StatusOK, gin.H{
			"status": job.Status,
			"result": job.Result,
			"error":  job.Error,
		})
	} else {
		// return error
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Job not found",
		})
	}
}

// getStatusesRequest gets the status of all jobs and returns it to the client
func (server *JobServer) getStatusesRequest(c *gin.Context) {
	// encode job map to json
	buf, err := json.Marshal(server.workMap)

	if err != nil {
		// error
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	} else {
		// send json
		c.JSON(http.StatusOK, buf)
	}
}
