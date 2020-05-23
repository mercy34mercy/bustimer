package infrastructure

import "net/http"

func TimeTableRequest(c Context) error {
	fr := c.GetTimeTableFrom()
	timeTable, ok := TimeTable[fr]
	if ok == false {
		return c.Response(http.StatusBadRequest, timeTable)
	}
	return c.Response(http.StatusOK, timeTable)
}
