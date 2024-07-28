package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type Method string

const (
	METHOD_GET    Method = "GET"
	METHOD_POST   Method = "POST"
	METHOD_PUT    Method = "PUT"
	METHOD_PATCH  Method = "PATCH"
	METHOD_DELETE Method = "DELETE"
	METHOD_OPTION Method = "OPTION"
)

func SimpleExchange(method Method, apiUrl, authorization string, params map[string]string, reqBody any) (jsonBody string, er error) {
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	parsedUrl, err := url.Parse(apiUrl)
	if err != nil {
		return "", err
	}

	if len(params) > 0 {
		p := url.Values{}
		for k, v := range params {
			p.Add(k, v)
		}
		parsedUrl.RawQuery = p.Encode()
	}

	req, err := http.NewRequest(string(method), parsedUrl.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	if len(authorization) > 0 {
		req.Header.Set("Authorization", authorization)
	}

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	// Read the response body
	respBodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	if response.StatusCode >= 400 {
		return string(respBodyBytes), fmt.Errorf("error: http response code: %v", response.StatusCode)
	}

	return string(respBodyBytes), nil
}
