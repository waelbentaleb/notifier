package notifier

import (
	"fmt"
	"log"
	"time"
)

type Notifier struct {
	// endpoint represents the post url
	endpoint string

	// maxWorker represents the number of concurred worker
	maxWorker int

	// jobQueue A buffered channel that we can send work requests on.
	jobQueue chan Job

	// jobResponsesChannel A buffered channel to send job responses on.
	jobResponsesChannel chan JobResponse

	// messages represents an array of messages to be sent
	messages []string

	// callback is function to be called when all messages had responses
	callback func([]FailedRequest)
}

type FailedRequest struct {
	id  int
	err error
}

// CreateNewNotifier Create and initialise Notifier struct
func CreateNewNotifier(endpoint string, maxWorker int, maxQueue int, messages []string, callback func([]FailedRequest)) *Notifier {
	return &Notifier{
		endpoint:            endpoint,
		maxWorker:           maxWorker,
		jobQueue:            make(chan Job, maxQueue),
		jobResponsesChannel: make(chan JobResponse, maxQueue),
		messages:            messages,
		callback:            callback,
	}
}

func (n *Notifier) Notify() {
	go func() {
		start := time.Now()
		defer func() { fmt.Println(time.Since(start)) }()

		// Initiate the dispatcher
		dispatcher := newDispatcher(n.maxWorker, n.jobQueue, n.jobResponsesChannel)
		dispatcher.run()

		// Fill the JobQueue
		for i, message := range n.messages {
			job := Job{
				id:       i,
				message:  `{"data": "` + message + `"}`,
				endpoint: n.endpoint,
			}
			n.jobQueue <- job
		}

		// Total number of jobs based on messages size
		jobCount := len(n.messages)

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

		log.Println(jobCount, responsesCount)

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
