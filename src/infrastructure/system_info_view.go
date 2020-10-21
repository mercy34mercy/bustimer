package infrastructure

import (
	"github.com/shun-shun123/bus-timer/src/domain"
	"github.com/shun-shun123/bus-timer/src/slack"
	"net/http"
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