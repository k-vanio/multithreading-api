package request

import "time"

type Data struct {
	From     string      `json:"from"`
	Response interface{} `json:"response"`
}

type Request interface {
	GetUrlWithFastestResponse(urls []string, limitTime time.Duration) (*Data, error)
}
