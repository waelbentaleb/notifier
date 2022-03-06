package notifier

type Dispatcher struct {
	// A pool of workers channels that are registered with the dispatcher
	workerPool    chan chan Job
	maxWorkers    int
	jobQueue      *chan Job
	responseQueue *chan Response
}

func newDispatcher(maxWorkers int, jobQueue *chan Job, responseQueue *chan Response) *Dispatcher {
	return &Dispatcher{
		workerPool:    make(chan chan Job, maxWorkers),
		maxWorkers:    maxWorkers,
		jobQueue:      jobQueue,
		responseQueue: responseQueue,
	}
}

func (d *Dispatcher) run() {
	// starting n number of workers
	for i := 0; i < d.maxWorkers; i++ {
		worker := newWorker(d.workerPool, d.responseQueue)
		worker.start()
	}

	go d.dispatch()
}

func (d *Dispatcher) dispatch() {
	for {
		select {
		case job := <-*d.jobQueue:
			// a job request has been received
			go func(job Job) {
				// try to obtain a worker job channel that is available.
				// this will block until a worker is idle
				jobChannel := <-d.workerPool

				// dispatch the job to the worker job channel
				jobChannel <- job
			}(job)
		}
	}
}
