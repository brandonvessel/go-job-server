package job

// job datatype
type Job struct {
	JobServerIP   string   `json:"job_server_ip"`       // job server ip
	JobServerPort string   `json:"job_server_port"`     // job server port
	UUID          string   `json:"uuid"`                // job data
	Name          string   `json:"name"`                // job name
	Data          []string `json:"data"`                // job data (parameters)
	Status        string   `json:"status"`              // job status
	Starttime     int      `json:"starttime"`           // job submission time
	Endtime       int      `json:"endtime default=0"` // job end time
	Result        string   `json:"result"`              // job result
	Error         string   `json:"error"`               // job error
}
