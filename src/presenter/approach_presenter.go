package presenter

import "github.com/shun-shun123/bus-timer/src/domain"

// convert html data passed from infrastructure to domain.ApproachInfo
func HtmlToApproach(html string) (domain.ApproachInfos, error) {
	approach := domain.ApproachInfos{
		ApproachInfo: []domain.ApproachInfo{
			domain.ApproachInfo{
				MoreMin: "20分",
				RealArrivalTime: "14:05",
				Direction: "立命館大学行き",
				Via: "パナソニック東口経由",
				ScheduledTime: "14:00",
				Delay: "5分遅延しています",
				BusStop: "4番のりば",
			},
			domain.ApproachInfo{
				MoreMin: "25分",
				RealArrivalTime: "14:10",
				Direction: "立命館大学行き",
				Via: "パナソニック西経由",
				ScheduledTime: "14:10",
				Delay: "",
				BusStop: "4番のりば",
			},
			domain.ApproachInfo{
				MoreMin: "30分",
				RealArrivalTime: "14:15",
				Direction: "立命館大学行き",
				Via: "かがやき通り経由",
				ScheduledTime: "14:15",
				Delay: "",
				BusStop: "4番のりば",
			},
		},
	}
	return approach, nil
}


