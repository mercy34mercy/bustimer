package domain

type ApproachInfo struct {
	MoreMin 		string `json:"more_min"`
	RealArrivalTime string `json:"real_arrival_time"`
	Direction 		string `json:"direction"`
	Via 			string `json:"via"`
	ScheduledTime 	string `json:"scheduled_time"`
	Delay 			string `json:"delay"`
	BusStop 		string `json:"bus_stop"`
}

type ApproachInfos struct {
	ApproachInfo []ApproachInfo `json:"approach_infos"`
}