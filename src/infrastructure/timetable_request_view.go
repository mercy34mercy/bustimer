package infrastructure

import (
	"fmt"
	"github.com/shun-shun123/bus-timer/src/config"
	"github.com/shun-shun123/bus-timer/src/slack"
	"net/http"
)

func TimeTableRequest(c Context) error {
	from, _ := c.GetFromToQuery()
	timeTable, ok := TimeTable[from]
	if ok == false {
		if from == config.Unknown {
			slack.PostMessage(fmt.Sprint("TimeTableRequest.BadQuery"))
			return c.Response("TimeTableRequest", http.StatusBadRequest, timeTable)
		} else {
			slack.PostMessage(fmt.Sprint("TimeTableRequest.StatusNoContent"))
			return c.Response("TimeTableRequest", http.StatusNoContent, timeTable)
		}
	}
	return c.Response("TimeTableRequest", http.StatusOK, timeTable)
}
