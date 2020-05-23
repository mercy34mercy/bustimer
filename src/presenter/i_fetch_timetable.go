package presenter

import "github.com/shun-shun123/bus-timer/src/domain"

type IFetchTimetable interface {
	FetchTimetable(from string) domain.TimeTable
}
