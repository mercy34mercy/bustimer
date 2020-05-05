package infrastructure

type Context interface {
	GetApproachInfoUrl() []string
	Response(int, interface{}) error
}
