package infrastructure

import "net/http"

func TimeTableRequest(c Context) error {
	_, fr := c.GetApproachInfoUrls()
	timeTable, ok := TimeTable[fr]
	if ok == false {
		return c.Response(http.StatusBadRequest, timeTable)
	}
	return c.Response(http.StatusOK, timeTable)
}
