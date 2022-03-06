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
	// generate random messages
	for i := 0; i < 10000; i++ {
		messages = append(messages, fmt.Sprintf(`{"data": "message %v"}`, i))
	}

	notificationService := notifier.CreateNewNotifier(endpoint, maxWorker, maxQueue)
	success, failedResponses := notificationService.Notify(messages)

	if !success {
		log.Printf("%d failed responses\n", len(failedResponses))
	}
}
