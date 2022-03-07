// TODO change package name
package main

import (
	"bufio"
	"fmt"
	"github.com/waelbentaleb/notifier"
	"log"
	"os"
	"time"
)

var (
	endpoint  = "https://waelbentaleb.free.beeceptor.com/my/api/path"
	maxWorker = 2000
	maxQueue  = 2500
)

func main() {
	var messages []string
	scanner := bufio.NewScanner(os.Stdin)
	messagesLock := make(chan int, 1)
	interval := 5

	go sendNewMessages(&messages, messagesLock, interval)

	for scanner.Scan() {
		messagesLock <- 1
		msg := scanner.Text()
		messages = append(messages, msg)
		<-messagesLock
	}

	// Check scanner error if exist
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

	select {}
}

func sendNewMessages(messages *[]string, chanLock chan int, interval int) {
	for {
		time.Sleep(time.Second * time.Duration(interval))
		fmt.Printf("Start sending %d messages\n", len(*messages))

		chanLock <- 1

		// Callback that receive a success boolean and an array with failed requests if they exist
		notifierCallback := func(failedRequests []notifier.FailedRequest) {
			if failedRequests != nil {
				log.Printf("You have %d failed requests\n", len(failedRequests))
			}
		}

		notificationClient := notifier.CreateNewNotifier(endpoint, maxWorker, maxQueue, notifierCallback)
		notificationClient.Notify(*messages)

		*messages = nil
		fmt.Println("Waiting for new messages")

		<-chanLock
	}
}
