package server

// the objective of the job server is to store jobs in a queue and process them to clients that request work

import (
	libjob "github.com/brandonvessel/go-job-server/lib/job"
	libqueue "github.com/brandonvessel/go-job-server/lib/queue"
	"github.com/gin-gonic/gin"
)

// represents a job queue
type JobServer struct {
	// listen and serve port
	Port string

	// result map
	workMap map[string]libjob.Job
	// work queue
	workQueue libqueue.Queue

	// gin router
	router *gin.Engine

	// config

	// PrintJobInfo enables the job server printing job information as it comes in (for debug, but not extensive info)
	PrintJobInfo bool `default:"false"`

	// Debug determines whether the server should print debug information
	Debug bool `default:"false"`

	// MaxQueueSize is the max queue size.
	MaxQueueSize int `default:"100"`

	// MaxworkMapDuration is the max result map item duration in minutes.
	MaxworkMapDuration int `default:"60"`

	// MaxworkQueueDuration is the maximum duration of a job in the work queue in minutes
	MaxworkQueueDuration int `default:"60"`

	// MaxMapSize is the maximum size of the map before deleting old jobs
	MaxMapSize int `default:"100"`
}

// NewJobServer returns a new job server object
func NewJobServer() *JobServer {
	// make new job server instance and return it
	js := &JobServer{}
	return js
}

// Initialize initializes the job server.
// The port variable is the port the server will listen on.
func (server *JobServer) Initialize(port string) {
	// set port
	server.Port = port

	// create gin router
	server.router = gin.Default()

	// set debug printing
	// TODO: Setup gin to hush when the PrintJobInfo variable is not set

	// Set a lower memory limit for multipart forms (default is 32 MiB)
	server.router.MaxMultipartMemory = 8 << 20 // 8 MiB

	// initialize the result map
	server.workMap = make(map[string]libjob.Job)

	// client routers
	// TODO: Make these paths make more sense. probably base them on client type (client, worker, etc)
	// job status router
	server.router.GET("/job/:uuid", server.getStatusRequest)

	// job statuses router
	server.router.GET("/job/results", server.getStatusesRequest)

	// job result router
	server.router.GET("/job/:uuid/result", server.getResultRequest)

	// job submission router
	server.router.POST("/job", server.addJobRequest)

	// job delete router
	server.router.DELETE("/job/:uuid", server.deleteJobRequest)

	// worker routers
	// get work router
	server.router.GET("/work", server.getWorkRequest)

	// finish work router
	server.router.POST("/work", server.postWorkRequest)
}

// Run starts the job server and listens for requests. Starts daemons in goroutines.
func (server *JobServer) Run() {
	// run the server
	// start result map cleaner
	go server.workMapCleaner()

	go server.workQueueCleaner()

	// start the job server
	go server.router.Run(":" + server.Port)
}
