package timetable

// newTimeTable timeTableを初期化して返す
func newTimeTable() timeTable {
	return timeTable{
		Version:  0,
		WeekDays: map[int][]oneTimeTable{},
		Saturday: map[int][]oneTimeTable{},
		Holidays: map[int][]oneTimeTable{},
	}
}
