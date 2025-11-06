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
	approachInfoUrls := c.GetApproachInfoUrls()

	if len(approachInfoUrls) == 0 {
		return c.Response("DebugNavitimeHTML", http.StatusBadRequest, map[string]string{
			"error": "No URLs generated. Please provide fr and to query parameters.",
		})
	}

	url := approachInfoUrls[0]
	log.Printf("Fetching HTML from: %s", url)

	// NAVITIME APIにリクエスト
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("HTTP GET failed: %v", err)
		return c.Response("DebugNavitimeHTML", http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("HTTP GET failed: %v", err),
		})
	}
	defer resp.Body.Close()

	// レスポンスボディを読み込み
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return c.Response("DebugNavitimeHTML", http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to read response body: %v", err),
		})
	}

	// HTMLをそのまま返す
	return c.Context.HTML(http.StatusOK, string(body))
}