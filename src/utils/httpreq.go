package utils

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
)

func DoRequest(req *http.Request) (res *http.Response, err error) {
	res, err = http.DefaultClient.Do(req)
	if err == nil && res != nil && res.StatusCode != http.StatusOK {
		err = errors.New("Got bad response " + res.Status)
	}
	return res, err
}

func HttpRequest(ctx context.Context, method, url, endpoint, token string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url+endpoint, body)
	if err != nil {
		return nil, err
	}

	if token != "" {
		req.Header.Add("Authorization", "Bearer "+token)
	}
	req.Header.Add("Content-Type", "application/json")

	if response, err := DoRequest(req); err != nil {
		return nil, err
	} else {
		return io.ReadAll(response.Body)
	}
}

func HttpGetBody(ctx context.Context, url, endpoint, token string) ([]byte, error) {
	return HttpRequest(ctx, "GET", url, endpoint, token, nil)
}

func HttpPost(ctx context.Context, url, endpoint, token string, body []byte) error {
	_, err := HttpRequest(ctx, "POST", url, endpoint, token, bytes.NewReader(body))
	return err
}
