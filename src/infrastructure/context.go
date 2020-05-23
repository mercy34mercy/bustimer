package infrastructure

type Context interface {
	GetApproachInfoUrls() ([]string, string)
	Response(int, interface{}) error
}
