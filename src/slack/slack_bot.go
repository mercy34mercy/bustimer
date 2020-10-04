package slack

import (
	"fmt"
	"github.com/ashwanthkumar/slack-go-webhook"
	"io/ioutil"
	"net/http"
)

func PostMessage(msg string) {
	payload := slack.Payload{
		Text: msg,
		Username: "robot",
		Channel: "#server-log",
		IconEmoji: ":thinking_face:",
		Attachments: nil,
	}
	err := slack.Send(webhook, "", payload)
	if len(err) > 0 {
		fmt.Printf("Error: %s\n", err)
	}
}

func PostScrapePageContent(url string) {
	resp, err := http.Get(url)
	defer resp.Body.Close()
	if err != nil {
		PostMessage("ERROR")
		return
	}
	byteArray, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		PostMessage("ERROR")
		return
	}
	PostMessage(string(byteArray))
}