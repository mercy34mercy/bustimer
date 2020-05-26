package infrastructure

import (
	"github.com/shun-shun123/bus-timer/src/config"
	"net/http"
)

func TimeTableRequest(c Context) error {
	_, from := c.GetApproachInfoUrls()
	timeTable, ok := TimeTable[from]
	if ok == false {
		if from == config.Unknown {
			return c.Response(http.StatusBadRequest, timeTable)
		} else {
			return c.Response(http.StatusNoContent, timeTable)
		}
	}
	return c.Response(http.StatusOK, timeTable)
}
