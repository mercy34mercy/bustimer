package infrastructure

import (
	"github.com/shun-shun123/bus-timer/src/domain"
	"net/http"
)

func SystemInfoRequest(c Context) error {
	systemInfo := domain.SystemInfo{
		Status: http.StatusInternalServerError,
		Message: "只今サーバでエラーが発生しており、原因を究明中です。ご不便をお掛けしてしまい申し訳ありません。",
		Time: "06:00ごろ復旧予定です。",
	}
	return c.Response("SystemInfoRequest", http.StatusOK, systemInfo)
}