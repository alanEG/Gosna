package main

import (
	"bytes"
	"crypto/tls"
	"net/http"
	"time"
)

var (
	transport = &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		IdleConnTimeout:   time.Second,
		DisableKeepAlives: true,
	}

	Client = &http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
)

func send_request(url, method string, headers map[string]string, requestBody []byte) (*http.Response, error) {
	Client.Timeout = time.Duration(flagTimeout) * time.Second

	req, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))

	//handling headers
	for headerKey, headerValue := range headers {
		req.Header.Add(headerKey, headerValue)
	}

	resp, err := Client.Do(req)
	return resp, err

	//This handling the request without Content-Type header
	//If there is no headers in response add in resp object a Header
	if len(resp.Header["Content-Type"]) < 1 {
		resp.Header["Content-Type"] = append(resp.Header["Content-Type"], "None/Null")
	}

	return resp, err
}
