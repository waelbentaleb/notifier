package notifier

type Notifier struct {
	// endpoint represents the post url
	endpoint string

	// maxWorker represents the number of concurred worker
	maxWorker int

	// jobQueue A buffered channel that we can send work requests on.
	jobQueue chan Job

	// jobResponsesChannel A buffered channel to send job responses on.
	jobResponsesChannel chan JobResponse

	// callback is function to be called when all messages had responses
	callback func([]FailedRequest)
}

type FailedRequest struct {
	id  int
	err error
}

// CreateNewNotifier Create and initialise Notifier struct
func CreateNewNotifier(endpoint string, maxWorker int, maxQueue int, callback func([]FailedRequest)) *Notifier {
	return &Notifier{
		endpoint:            endpoint,
		maxWorker:           maxWorker,
		jobQueue:            make(chan Job, maxQueue),
		jobResponsesChannel: make(chan JobResponse, maxQueue),
		callback:            callback,
	}
}

func (n *Notifier) Notify(messages []string) {
	go func() {
		// Initiate the dispatcher
		dispatcher := newDispatcher(n.maxWorker, n.jobQueue, n.jobResponsesChannel)
		dispatcher.run()

		// Fill the JobQueue
		for i, message := range messages {
			job := Job{
				id:       i,
				message:  `{"data": "` + message + `"}`,
				endpoint: n.endpoint,
			}
			n.jobQueue <- job
		}

		// Total number of jobs based on messages size
		jobCount := len(messages)

		// Retrieve job responses and return failures to the callback
		failedRequests := handleJobResponses(n.jobResponsesChannel, jobCount)
		n.callback(failedRequests)
	}()
}

func handleJobResponses(jobResponsesChannel chan JobResponse, jobCount int) []FailedRequest {
	var failedRequests []FailedRequest
	responsesCount := 0

	for response := range jobResponsesChannel {
		// Update responses count
		responsesCount++

		// Save failed requests
		if response.err != nil {
			failedRequests = append(failedRequests, FailedRequest{response.id, response.err})
		}

		// Wait to receive all job responses before breaking
		if responsesCount == jobCount {
			break
		}
	}

	if len(failedRequests) == 0 {
		return nil
	}

	return failedRequests
}
