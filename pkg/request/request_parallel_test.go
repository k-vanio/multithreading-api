package request_test

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/k-vanio/multithreading-api/pkg/request"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	jsonStr = `{ 
		"USDBRL": { 
			"code": "USD", 
			"codein": 
			"BRL", 
			"name": "Dólar Americano/Real Brasileiro",
			"high": "4.75",
			"low": "4.6963",
			"varBid": "-0.0095",
			"pctChange": "-0.2",
			"bid": "4.7314",
			"ask": "4.7344",
			"timestamp": "1690577990",
			"create_date": "2023-07-28 17:59:50"
		}
	}`
)

type ClientStub struct {
	mock.Mock
}

func (m *ClientStub) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	return args.Get(0).(*http.Response), args.Error(1)
}

func TestGetUrlWithFastestResponse(t *testing.T) {
	urls := []string{
		"https://viacep.com.br/ws/01311100/json/",
		"https://brasilapi.com.br/api/cep/v2/01311100",
	}

	clientStub := new(ClientStub)

	response := &http.Response{Body: io.NopCloser(strings.NewReader(jsonStr))}
	clientStub.On("Do", mock.Anything).Return(response, nil)

	httpRequest := request.NewRequestParallel(clientStub)

	result, err := httpRequest.GetUrlWithFastestResponse(urls, time.Second*1)

	assert.NoError(t, err)
	assert.NotNil(t, result, "")
}

func TestGetUrlWithFastestResponseError(t *testing.T) {
	urls := []string{
		"https://viacep.com.br/ws/01311100/json/",
		"https://brasilapi.com.br/api/cep/v2/01311100",
	}

	clientStub := new(ClientStub)

	response := &http.Response{Body: io.NopCloser(strings.NewReader(jsonStr))}
	clientStub.On("Do", mock.Anything).Return(response, errors.New("bad request"))

	httpRequest := request.NewRequestParallel(clientStub)

	result, err := httpRequest.GetUrlWithFastestResponse(urls, time.Second*1)

	assert.Error(t, err)
	assert.Nil(t, result)
}
