package googlechat

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/RoadieHQ/kubewise/kwrelease"
	"github.com/RoadieHQ/kubewise/presenters"
	"helm.sh/helm/v3/pkg/release"
)

// GoogleChat represents the ability to send notifications to Google Hangouts Chat.
// https://chat.google.com
type GoogleChat struct {
	WebhookURL string
}

// Init retrieves configuration properties from environment variables and stores them in the
// GoogleChat instance.
func (g *GoogleChat) Init() {
	var webhookURL string
	if value, ok := os.LookupEnv("KW_GOOGLECHAT_WEBHOOK_URL"); ok {
		webhookURL = value
	} else {
		log.Fatalln("Missing environment variable KW_GOOGLECHAT_WEBHOOK_URL")
	}

	g.WebhookURL = webhookURL
}

// HandleEvent sends notifications when release events occur.
func (g *GoogleChat) HandleEvent(releaseEvent *kwrelease.Event) {
	if msg := presenters.PrepareMsg(releaseEvent); msg != "" {
		makeRequest(g, msg)
	}
}

// HandleServerStartup sends notifications when KubeWise starts up.
func (g *GoogleChat) HandleServerStartup(releases []*release.Release) {
	if msg := presenters.PrepareServerStartupMsg(releases); msg != "" {
		makeRequest(g, msg)
	}
}

func makeRequest(g *GoogleChat, text string) []byte {
	responseBody := []byte{}
	values := map[string]string{"text": text}
	jsonValue, marshalError := json.Marshal(values)

	if marshalError != nil {
		// msg should never contain sensitive information because it's being sent to a third-party
		// application so logging this error is secure.
		log.Println("Error marshaling message into Json", marshalError)
		return responseBody
	}

	contentType := "application/json; charset=UTF-8"
	resp, requestErr := http.Post(g.WebhookURL, contentType, bytes.NewBuffer(jsonValue))

	if requestErr != nil {
		// Do NOT log the err. It contains the URL which contains sensitive authentication data.
		// If this is to be logged in future, strip the sensitive data from the URL before logging.
		log.Println("Error making request to Google Hangouts Chat")
		return responseBody
	}

	defer resp.Body.Close()
	responseBody, readBodyErr := ioutil.ReadAll(resp.Body)
	if readBodyErr != nil {
		// Do NOT log the err. It could contain the URL which contains sensitive authentication data.
		log.Println("Malformed response received from Google Hangouts Chat")
		return responseBody
	}

	return responseBody
}
