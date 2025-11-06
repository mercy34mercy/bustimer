package infrastructure

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/shun-shun123/bus-timer/src/config"
	"github.com/shun-shun123/bus-timer/src/domain"
	"github.com/shun-shun123/bus-timer/src/slack"
)

func SystemInfoRequest(c Context) error {
	systemInfo := domain.SystemInfo{
		Status: http.StatusOK,
		Message: "",
		Time: "",
	}
	return c.Response("SystemInfoRequest", http.StatusOK, systemInfo)
}

func DebugErrorSystemInfoRequest(c Context) error {
	systemInfo := domain.SystemInfo{
		Status: http.StatusInternalServerError,
		Message: "Debug用出力。只今サーバでエラーが発生しており、復旧中です。",
		Time: "12:00ごろ復旧予定です。",
	}
	slack.PostMessage("SystemInfoRequestError(DEBUG)")
	return c.Response("SystemInfoRequest", http.StatusOK, systemInfo)
}

func DebugSuccessSystemInfoRequest(c Context) error {
	systemInfo := domain.SystemInfo{
		Status: http.StatusOK,
		Message: "Debug用出力",
		Time: "",
	}
	slack.PostMessage("DebugSuccessSystemInfoRequest(DEBUG)")
	return c.Response("SystemInfoRequest", http.StatusOK, systemInfo)
}

func DebugNavitimeHTML(c Context) error {
	// クエリパラメータからfrとtoを取得
	fr := c.Context.QueryParam("fr")
	to := c.Context.QueryParam("to")
	log.Printf("Debug endpoint called with fr=%s, to=%s", fr, to)

	approachInfoUrls := c.GetApproachInfoUrls()
	log.Printf("Generated %d URLs", len(approachInfoUrls))

	if len(approachInfoUrls) == 0 {
		log.Printf("No URLs generated from fr=%s, to=%s", fr, to)
		return c.Response("DebugNavitimeHTML", http.StatusBadRequest, map[string]string{
			"error": "No URLs generated. Please provide fr and to query parameters.",
		})
	}

	url := approachInfoUrls[0]
	log.Printf("Fetching HTML from: %s", url)

	// HTTPクライアントを作成してヘッダーを設定
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return c.Response("DebugNavitimeHTML", http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to create request: %v", err),
		})
	}

	// ブラウザを模倣するヘッダーを追加
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "ja,en-US;q=0.9,en;q=0.8")
	req.Header.Set("Referer", "https://transfer-cloud.navitime.biz/")

	// NAVITIME APIにリクエスト
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("HTTP GET failed: %v", err)
		return c.Response("DebugNavitimeHTML", http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("HTTP GET failed: %v", err),
		})
	}
	defer resp.Body.Close()

	log.Printf("HTTP GET response: status=%d, statusCode=%s", resp.StatusCode, resp.Status)

	// レスポンスボディを読み込み
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return c.Response("DebugNavitimeHTML", http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to read response body: %v", err),
		})
	}

	log.Printf("Response body length: %d bytes", len(body))

	// HTMLをそのまま返す（ステータスコードも一緒に返す）
	return c.Context.HTMLBlob(resp.StatusCode, body)
}