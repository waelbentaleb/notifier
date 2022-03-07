package notifier

// Worker represents the worker that executes the job
type Worker struct {
	workerPool          chan chan Job
	jobChannel          chan Job
	quit                chan bool
	jobResponsesChannel chan JobResponse
}

// JobResponse represents a single query response with query id and an error if exist
type JobResponse struct {
	id         int
	statusCode int
	err        error
}

func newWorker(workerPool chan chan Job, jobResponsesChannel chan JobResponse) Worker {
	return Worker{
		workerPool:          workerPool,
		jobChannel:          make(chan Job),
		quit:                make(chan bool),
		jobResponsesChannel: jobResponsesChannel,
	}
}

// Start method starts the run loop for the worker
func (w Worker) start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.workerPool <- w.jobChannel

			select {
			case job := <-w.jobChannel:
				// we have received a work request.
				statusCode, err := job.postNotification()
				w.jobResponsesChannel <- JobResponse{job.id, statusCode, err}
			}
		}
	}()
}
