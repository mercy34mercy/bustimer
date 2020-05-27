package presenter

import (
	"github.com/shun-shun123/bus-timer/src/config"
	"github.com/shun-shun123/bus-timer/src/domain"
)

type IFetchTimetable interface {
	FetchTimetable(from config.From, to config.To) domain.TimeTable
}
