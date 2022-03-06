package notifier

// Worker represents the worker that executes the job
type Worker struct {
	workerPool    chan chan Job
	jobChannel    chan Job
	quit          chan bool
	responseQueue chan Response
}

func newWorker(workerPool chan chan Job, responseQueue chan Response) Worker {
	return Worker{
		workerPool:    workerPool,
		jobChannel:    make(chan Job),
		quit:          make(chan bool),
		responseQueue: responseQueue,
	}
}

// Start method starts the run loop for the worker, listening for a quit channel in
// case we need to stop it
func (w Worker) start() {
	go func() {
		for {
			// register the current worker into the worker queue.
			w.workerPool <- w.jobChannel

			select {
			case job := <-w.jobChannel:
				// we have received a work request.
				err := job.postNotification()
				w.responseQueue <- Response{job.id, err}
			}
		}
	}()
}
