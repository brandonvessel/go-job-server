package main

import (
	"fmt"

	gc "github.com/brandonvessel/go-job-server/client"
)

func main() {
	// create client for sending jobs
	client := gc.NewClient("127.0.0.1", "8888")

	// send Add job
	jobUUID, err := client.NewJob(
		"Add",
		[]string{"1", "2"},
	)
	if err != nil {
		fmt.Println("Could not send job")
	}

	result, err := client.GetResult(jobUUID)
	if err != nil {
		fmt.Println("Job error:", err)
	}

	fmt.Println(result)
}
