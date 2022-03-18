package main

import (
	gs "github.com/brandonvessel/go-job-server/server"
)

func main() {
	// create server instance
	server := gs.NewJobServer()
	
	// initialize and set port
	server.Initialize("8888")
	
	// run server
	server.Run()

	// the server runs in the background, so the main thread can do other things
	select {}
}
