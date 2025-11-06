package slack

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ashwanthkumar/slack-go-webhook"
)

func PostMessage(msg string) {
	// Webhook URLが設定されていない場合はスキップ
	if webhook == "" {
		log.Printf("Slack webhook not configured, skipping message: %s", msg)
		return
	}

	payload := slack.Payload{
		Text:        msg,
		Username:    "robot",
		Channel:     "#server-log",
		IconEmoji:   ":thinking_face:",
		Attachments: nil,
	}
	err := slack.Send(webhook, "", payload)
	if len(err) > 0 {
		log.Printf("Error sending Slack message (ignoring): %s", err)
	}
}

func PostScrapePageContent(url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("PostScrapePageContent: HTTP GET failed for %s: %v", url, err)
		PostMessage(fmt.Sprintf("ERROR: Failed to fetch %s", url))
		return
	}
	defer resp.Body.Close()

	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("PostScrapePageContent: Failed to read response body: %v", err)
		PostMessage("ERROR: Failed to read response")
		return
	}
	PostMessage(string(byteArray))
}