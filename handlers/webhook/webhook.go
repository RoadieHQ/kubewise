package webhook

import (
	"bytes"
	"log"
	"net/http"
	"os"

	"github.com/larderdev/kubewise/kwrelease"
	"github.com/larderdev/kubewise/presenters"
)

type Webhook struct {
	URL    string
	Method string
}

func (w *Webhook) Init() error {
	var url string

	method := "POST"
	if value, ok := os.LookupEnv("KW_WEBHOOK_METHOD"); ok {
		method = value
	}

	if value, ok := os.LookupEnv("KW_WEBHOOK_URL"); ok {
		url = value
	} else {
		log.Fatalln("Missing environment variable KW_WEBHOOK_URL")
	}

	w.Method = method
	w.URL = url
	return nil
}

func (w *Webhook) HandleEvent(releaseEvent *kwrelease.Event) {
	jsonStr, jsonErr := presenters.ReleaseEventToJSON(releaseEvent)

	if jsonErr != nil {
		log.Println("Error encoding JSON in webhook event", jsonErr)
	}

	client := &http.Client{}
	req, reqErr := http.NewRequest(w.Method, w.URL, bytes.NewBuffer(jsonStr))

	if reqErr != nil {
		log.Println("Error forming request in webhook event", reqErr)
	}

	req.Header.Add("Content-Type", "application/json")

	if value, ok := os.LookupEnv("KW_WEBHOOK_AUTH_TOKEN"); ok {
		req.Header.Add("Authorization", "Bearer "+value)
	}

	resp, respErr := client.Do(req)

	if respErr != nil {
		log.Println("Error handling response to webhook event", respErr)
	}

	log.Println("Successful response received from", w.Method, w.URL, ":", resp.StatusCode)
}
