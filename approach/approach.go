package approach

const (
	url          = "https://ohmitetudo-bus.jorudan.biz/busstatedtl"
	frRits       = "立命館大学〔近江鉄道・湖国バス〕"
	frMinakusa   = "南草津駅〔近江鉄道・湖国バス〕"
	MAX_RESPONSE = 3
)

// frクエリに対応するdgmplクエリのスライス
var dgmplMap = map[string][]string{frMinakusa: []string{"南草津駅〔近江鉄道・湖国バス〕:1", "南草津駅〔近江鉄道・湖国バス〕:3", "南草津駅〔近江鉄道・湖国バス〕:4"},
	frRits: []string{"立命館大学〔近江鉄道・湖国バス〕:2"}}

// approachInfo 1件の接近情報を表す構造体
type approachInfo struct {
	// あと何分で到着か
	MoreMin string `json:"more_min"`
	// 実際の到着予定時刻(遅延を考慮する)
	RealArrivalTime string `json:"real_arrive_time"`
	// 系統
	Descent string `json:"descent"`
	// 行き先
	Direction string `json:"direction"`
}
