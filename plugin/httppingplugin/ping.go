package httppingplugin

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var httpGet = func(url string) (resp *http.Response, err error) {
	return http.Get(url)
}

func httpPing(u string) (int, error) {
	_, err := url.ParseRequestURI(u)
	if err != nil {
		return -1, fmt.Errorf("Not a valid url")
	}

	start := time.Now()
	resp, err := httpGet(u)
	elapsed := int(time.Since(start).Nanoseconds() / 1000 / 1000)

	if err != nil {
		return elapsed, fmt.Errorf("Error pinging the url: %s", err)
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	return elapsed, nil
}
