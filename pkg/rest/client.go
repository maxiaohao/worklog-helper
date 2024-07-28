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
	GET    Method = "GET"
	POST   Method = "POST"
	PUT    Method = "PUT"
	PATCH  Method = "PATCH"
	DELETE Method = "DELETE"
	OPTION Method = "OPTION"
)

func SimpleExchange(method Method, apiUrl, authorization string, params map[string]string, reqBody any) ([]byte, error) {
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	parsedUrl, err := url.Parse(apiUrl)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	if len(authorization) > 0 {
		req.Header.Set("Authorization", authorization)
	}

	// req.Header.Set("Content-Type", "application/json")
	// auth := "username:password"
	// token := base64.StdEncoding.EncodeToString([]byte(auth))
	// req.Header.Set("Authorization", "Basic "+token)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Read the response body
	respBodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 400 {
		return respBodyBytes, fmt.Errorf("error: http response code: %v", response.StatusCode)
	}

	return respBodyBytes, nil
}
