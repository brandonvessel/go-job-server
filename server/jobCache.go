package server

// removeJob removes a job from the work queue map by its UUID
func (server *JobServer) removeJob(uuid string) {
	// remove job from job map
	delete(server.workMap, uuid)
}