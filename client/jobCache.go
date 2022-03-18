package client

import (
	libjob "github.com/brandonvessel/go-job-server/lib/job"
)

// purgeCache removes all jobs from the cache
func (client *Client) purgeCache() {
	for k := range client.jobCache {
		delete(client.jobCache, k)
	}
}

// CancelJob cancels a job
func (client *Client) CancelJob(job *libjob.Job) {
	if job.UUID == "" {
		return
	}

	client.deleteJob(job)
	delete(client.jobCache, job.UUID)
}

// CancelJobByUUID cancels a job associated with the given UUID
func (client *Client) CancelJobByUUID(uuid string) {
	job := client.jobCache[uuid]

	if job.UUID == "" {
		return
	}

	client.deleteJob(&job)
	delete(client.jobCache, job.UUID)
}

// CancelAllJobs cancels all jobs
func (client *Client) CancelAllJobs() {
	for _, job := range client.jobCache {
		go client.CancelJob(&job)
	}

	client.purgeCache()
}
