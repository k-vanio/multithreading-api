package request

import (
	"context"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

type Response struct {
	Url      string
	Response *http.Response
}

type RequestParallel struct {
	client    Client
	waitGroup *sync.WaitGroup
}

func (r *RequestParallel) byRequest(url string, ctx context.Context, responseChan chan<- Response, errorChan chan<- error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		errorChan <- err
		r.waitGroup.Done()
		return
	}

	res, err := r.client.Do(req)
	if err != nil {
		errorChan <- err
		r.waitGroup.Done()
		return
	}

	responseChan <- struct {
		Url      string
		Response *http.Response
	}{Url: url, Response: res}
	r.waitGroup.Done()
}

func (r *RequestParallel) GetUrlWithFastestResponse(urls []string, limitTime time.Duration) (*Data, error) {
	data := new(Data)

	responseChan := make(chan Response)
	errorChan := make(chan error)

	ctx, cancel := context.WithTimeout(context.Background(), limitTime)
	defer cancel()

	for _, value := range urls {
		r.waitGroup.Add(1)
		go r.byRequest(value, ctx, responseChan, errorChan)
	}

	go func(w *sync.WaitGroup) {
		w.Wait()
		close(responseChan)
		close(errorChan)
	}(r.waitGroup)

	select {
	case res := <-responseChan:
		json.NewDecoder(res.Response.Body).Decode(&data.Response)
		data.From = res.Url

		return data, nil
	case err := <-errorChan:
		return nil, err
	}
}

func NewRequestParallel(client Client) *RequestParallel {
	return &RequestParallel{client, &sync.WaitGroup{}}
}
