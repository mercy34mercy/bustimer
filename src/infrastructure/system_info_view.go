package infrastructure

import (
	"github.com/shun-shun123/bus-timer/src/domain"
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