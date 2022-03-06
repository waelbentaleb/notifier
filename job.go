package notifier

import (
	"bytes"
	"net/http"
	"time"
)

// Job represents the job to be run
type Job struct {
	id       int
	message  string
	endpoint string
}

func (job Job) postNotification() error {
	jsonData := []byte(job.message)

	request, err := http.NewRequest("POST", job.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	request.Header.Set("Content-Type", "application/json; charset=UTF-8")
	client := &http.Client{Timeout: 10 * time.Second}

	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	//log.Println(p.message)
	return nil
}
