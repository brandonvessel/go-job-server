package server

import (
	"time"

	libjob "github.com/brandonvessel/go-job-server/lib/job"
)
// workMapCleaner cleans the result map every minute to remove old items from the map if they are too old or the map is too big
func (server *JobServer) workMapCleaner() {
	for {
		// sleep for a minute
		time.Sleep(time.Minute)

		// iterate through the work map and remove items that are too old
		for uuid, job := range server.workMap {
			// if the job is too old
			if (int(time.Now().Unix()) - job.Starttime) > server.MaxworkMapDuration*60 {
				// remove the job from the map
				delete(server.workMap, uuid)
			}
		}

		// while the map is too big, remove oldest item
		for len(server.workMap) > server.MaxMapSize {
			// iterate and calculate oldest item
			oldest := int(time.Now().Unix())
			oldestUUID := ""
			for uuid, job := range server.workMap {
				// if job is older than oldest
				if job.Starttime < oldest {
					// set oldest to job starttime
					oldest = job.Starttime
					// set oldestUUID to job UUID
					oldestUUID = uuid
				}
			}

			// delete oldest item
			server.removeJob(oldestUUID)
		}
	}
}

// workQueueCleaner cleans the work queue every minute to remove old items from the queue if they are too old or the queue is too big.
func (server *JobServer) workQueueCleaner() {
	for {
		// sleep for a minute
		time.Sleep(time.Minute)

		// if the queue is empty, continue
		if server.workQueue.Empty() {
			continue
		}

		// while the queue is too big, remove items and delete them from the work map
		for server.workQueue.Size() > server.MaxQueueSize {
			// get oldest item
			oldest, err := server.workQueue.Dequeue()

			// if there was an error
			if err != nil {
				break
			}

			// remove oldest item from work map (also typecase oldest to libjob.Job)
			server.removeJob(oldest.(libjob.Job).UUID)
		}
	}
}