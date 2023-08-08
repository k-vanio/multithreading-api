package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/k-vanio/multithreading-api/pkg/request"
)

func main() {
	urls := []string{
		"https://viacep.com.br/ws/01311100/json/",
		"https://brasilapi.com.br/api/cep/v2/01311100",
	}

	client := request.NewRequestParallel(&http.Client{})

	response, err := client.GetUrlWithFastestResponse(urls, time.Second*1)
	if err != nil {
		fmt.Println(err)
		return
	}

	json.NewEncoder(os.Stdout).Encode(response)
}
