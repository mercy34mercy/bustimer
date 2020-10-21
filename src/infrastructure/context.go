package infrastructure

import "github.com/shun-shun123/bus-timer/src/config"

type Context interface {
	GetApproachInfoUrls() []string
	Response(string, int, interface{}) error
	GetFromToQuery() (config.From, config.To)
}
