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

	// responseQueue A buffered channel that we can send request results on.
	responseQueue chan Response
}

// Response represents a single query response with query id and an error if exist
type Response struct {
	id  int
	err error
}

// CreateNewNotifier Create and initialise Notifier struct
func CreateNewNotifier(endpoint string, maxWorker int, maxQueue int) *Notifier {
	return &Notifier{
		endpoint:      endpoint,
		maxWorker:     maxWorker,
		jobQueue:      make(chan Job, maxQueue),
		responseQueue: make(chan Response, maxQueue),
	}
}

func (n Notifier) Notify(messages []string) (bool, []Response) {
	start := time.Now()
	defer func() { fmt.Println(time.Since(start)) }()

	// Initiate the dispatcher
	dispatcher := newDispatcher(n.maxWorker, &n.jobQueue, &n.responseQueue)
	dispatcher.run()

	// Fill the JobQueue
	for i, message := range messages {
		job := Job{
			id:       i,
			message:  message,
			endpoint: n.endpoint,
		}

		n.jobQueue <- job
	}

	messagesLength := len(messages)
	responsesLength := 0

	var failedResponses []Response

	// Retrieve workers result
	for response := range n.responseQueue {

		// Save failed requests
		if response.err != nil {
			failedResponses = append(failedResponses, response)
		}

		responsesLength++
		log.Println(messagesLength, responsesLength)

		// Wait to receive all results before killing main goroutine
		if responsesLength == messagesLength {
			break
		}
	}

	return len(failedResponses) == 0, failedResponses
}
