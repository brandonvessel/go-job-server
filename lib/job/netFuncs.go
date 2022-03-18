package job

import (
	"encoding/json"
	"errors"
	"net/http"
)

// GetResult queries the job server for the result of the job and returns it if it is available.
func (job *Job) GetResult() (string, error) {
	// if job is not done, fetch result from server
	if job.Status != "done" && job.Status != "error" {
		// ask server for job result
		// create request
		req, err := http.NewRequest("GET", "http://"+job.JobServerIP+":"+job.JobServerPort+"/job/"+job.UUID + "/result", nil)

		if err != nil {
			return "", err
		}

		// create http client
		client := &http.Client{}

		// send request
		resp, err := client.Do(req)

		if err != nil {
			return "", err
		}

		defer resp.Body.Close()

		// decode json response
		var mjob Job
		err = json.NewDecoder(resp.Body).Decode(&mjob)

		// save remote job stats to local job
		job.Endtime = mjob.Endtime
		job.Status = mjob.Status
		job.Result = mjob.Result
		job.Error = mjob.Error

		if err != nil {
			return "", err
		}

		// return result if job is done or error
		if job.Status == "queued" || job.Status == "error"{
			return job.Result, nil
		} else {
			return "", errors.New("job incomplete")
		}
	}

	// return job result
	return job.Result, nil
}
