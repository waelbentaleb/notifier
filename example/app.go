package main

import (
	"bufio"
	"flag"
	"github.com/waelbentaleb/notifier"
	"log"
	"os"
	"time"
)

var (
	endpoint  string
	interval  int
	maxWorker int
	maxQueue  int
	messages  []string
)

const (
	DefaultEndpoint  = "https://waelbentaleb.free.beeceptor.com/my/api/path"
	DefaultInterval  = 5
	DefaultMaxWorker = 2000
	DefaultMaxQueue  = 2500
)

func main() {

	// Define application flags
	flag.StringVar(&endpoint, "url", DefaultEndpoint, "notification endpoint to send messages")
	flag.IntVar(&interval, "interval", DefaultInterval, "interval in seconds before sending new messages")
	flag.IntVar(&maxWorker, "maxWorker", DefaultMaxWorker, "represents the maximum number of concurrent worker")
	flag.IntVar(&maxQueue, "maxQueue", DefaultMaxQueue, "represents the maximum length of messages queue")
	flag.Parse()

	scanner := bufio.NewScanner(os.Stdin)
	messagesLock := make(chan int, 1)

	// Goroutine to send messages every interval (in seconds)
	go sendNewMessages(&messages, messagesLock)

	// Reading new stdin messages
	for scanner.Scan() {
		// Lock messages array
		messagesLock <- 1

		// Read new message and append it to messages array
		msg := scanner.Text()
		messages = append(messages, msg)

		// Unlock messages array
		<-messagesLock
	}

	// Check scanner error if exist
	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

	// I'm confusing by the task description between two different expected behaviors,
	// that's why I add this select to prevent exit
	select {}
}

func sendNewMessages(messages *[]string, messagesLock chan int) {
	for {
		time.Sleep(time.Second * time.Duration(interval))

		if len(*messages) != 0 {
			log.Printf("Start sending %d messages\n", len(*messages))
		}

		// Lock messages array
		messagesLock <- 1

		// Callback function called when all messages are sent
		notifierCallback := func(failedRequests []notifier.FailedRequest) {
			if failedRequests != nil {
				log.Printf("You have %d failed requests\n", len(failedRequests))
				return
			}
			log.Println("All messages are send successfully")
		}

		// Initiate notification client
		notificationClient := notifier.CreateNewNotifier(endpoint, maxWorker, maxQueue, notifierCallback)

		// Non-blocking call to send notifications
		notificationClient.Notify(*messages)

		// Clean messages array
		*messages = nil

		// Unlock messages array
		<-messagesLock
	}
}
