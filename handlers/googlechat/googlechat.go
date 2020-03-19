package googlechat

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/larderdev/kubewise/presenters"
	"helm.sh/helm/v3/pkg/release"
)

type GoogleChat struct {
	WebhookURL string
}

func (g *GoogleChat) Init() error {
	var webhookURL string
	if value, ok := os.LookupEnv("KW_GOOGLECHAT_WEBHOOK_URL"); ok {
		webhookURL = value
	} else {
		log.Fatalln("Missing environment variable KW_GOOGLECHAT_WEBHOOK_URL")
	}

	g.WebhookURL = webhookURL
	return nil
}

func (g *GoogleChat) ObjectCreated(currentRelease, previousRelease *release.Release) {
	if msg := presenters.PrepareObjectCreatedMsg(currentRelease, previousRelease); msg != "" {
		makeRequest(g, msg)
	}
}

func (g *GoogleChat) ObjectDeleted(currentRelease, previousRelease *release.Release) {
	if msg := presenters.PrepareObjectDeletedMsg(currentRelease, previousRelease); msg != "" {
		makeRequest(g, msg)
	}
}

func (g *GoogleChat) ObjectUpdated(currentRelease, previousRelease *release.Release) {
	if msg := presenters.PrepareObjectUpgradedMsg(currentRelease, previousRelease); msg != "" {
		makeRequest(g, msg)
	}
}

func makeRequest(g *GoogleChat, text string) []byte {
	responseBody := []byte{}
	values := map[string]string{"text": text}
	jsonValue, marshalError := json.Marshal(values)

	if marshalError != nil {
		log.Println("Error marshaling message into Json", marshalError)
		return responseBody
	}

	contentType := "application/json; charset=UTF-8"
	resp, requestErr := http.Post(g.WebhookURL, contentType, bytes.NewBuffer(jsonValue))

	if requestErr != nil {
		log.Println("Error making request to Google Hangouts Chat", requestErr)
		return responseBody
	}

	defer resp.Body.Close()
	responseBody, readBodyErr := ioutil.ReadAll(resp.Body)
	if readBodyErr != nil {
		log.Println("Malformed response received from Google Hangouts Chat", readBodyErr)
		return responseBody
	}

	return responseBody
}
