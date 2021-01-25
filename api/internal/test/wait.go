package test

import (
	internalHttp "api/internal/http"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

func WaitAppHealth(ticker *time.Ticker, healthCheckUrl string) *sync.WaitGroup {
	return waitApp(ticker, healthCheckUrl, func(res *http.Response, err error) bool {
		return err == nil && res.StatusCode == 200
	})
}

func WaitAppReady(ticker *time.Ticker, healthCheckUrl string) *sync.WaitGroup {
	return waitApp(ticker, healthCheckUrl, func(res *http.Response, err error) bool {
		if err != nil || res.StatusCode != 200 {
			return false
		}

		defer res.Body.Close()
		healthResponse := internalHttp.Health{}
		err = json.NewDecoder(res.Body).Decode(&healthResponse)
		if err != nil {
			return false
		}

		return healthResponse.Status == internalHttp.StatusUP
	})
}

func waitApp(ticker *time.Ticker, healthCheckUrl string, checkOk func(resp *http.Response, err error) bool) *sync.WaitGroup {
	var wg sync.WaitGroup
	wg.Add(1)

	client := http.Client{
		Timeout: 1 * time.Second,
	}

	go func() {
		for {
			select {
			case <-ticker.C:
				res, err := client.Get(healthCheckUrl)
				if checkOk(res, err) {
					wg.Done()
					return
				}

				println("wait...")
			}
		}
	}()

	return &wg
}
