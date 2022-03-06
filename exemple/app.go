// TODO change package name
package main

import (
	"fmt"
	"github.com/waelbentaleb/notifier"
	"log"
)

var (
	endpoint  = "https://waelbentaleb.free.beeceptor.com/my/api/path"
	maxWorker = 1500
	maxQueue  = 2000
)

func main() {
	var messages []string

	// TODO Optimizing reading and sending
	// Generate random messages
	for i := 0; i < 100000; i++ {
		messages = append(messages, fmt.Sprintf("message %d", i))
	}

	// Callback that receive a success boolean and an array with failed requests if they exist
	notifierCallback := func(failedRequests []notifier.FailedRequest) {
		if failedRequests != nil {
			log.Printf("You have %d failed requests\n", len(failedRequests))
		}
	}

	notificationService := notifier.CreateNewNotifier(endpoint, maxWorker, maxQueue, messages, notifierCallback)
	notificationService.Notify()

	log.Println("We are good BOSS !!")

	// Simulate server listening with an infinite wait
	select {}
}
