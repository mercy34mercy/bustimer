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
	// HTTPクライアントを作成してヘッダーを設定
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("PostScrapePageContent: Failed to create request for %s: %v", url, err)
		PostMessage(fmt.Sprintf("ERROR: Failed to create request for %s", url))
		return
	}

	// ブラウザを模倣するヘッダーを追加
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "ja,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Referer", "https://transfer-cloud.navitime.biz/")

	resp, err := client.Do(req)
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