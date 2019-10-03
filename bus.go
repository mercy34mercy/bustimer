package main

// busTime busの時間情報をもつ構造体
type busTime struct {
	Leave       string `json:"leave"`
	Arrive      string `json:"arrive"`
	Via         string `json:"via"`
	IsDoubleBus bool   `json:"is_double_bus"`
	Fee         int    `json:"fee"`
	BusStop     string `json:"bus_stop"`
}

func newBusTime(leave, arrive, via string, isDoubleBus bool, fee int, busStop string) busTime {
	return busTime{
		Leave:       leave,
		Arrive:      arrive,
		Via:         via,
		IsDoubleBus: isDoubleBus,
		Fee:         fee,
		BusStop:     busStop,
	}
}
