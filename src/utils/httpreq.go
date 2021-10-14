package utils

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
)

func HttpRequest(ctx context.Context, method, url, endpoint, token string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, method, url + endpoint, body)
	if err != nil {
		return nil, err
	}

	if token != "" {
		req.Header.Add("Authorization", "Bearer " + token)
	}
	req.Header.Add("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		return nil, errors.New("Got bad response " + response.Status)
	}
	return io.ReadAll(response.Body)
}

func HttpGetBody(ctx context.Context, url, endpoint, token string) ([]byte, error) {
	return HttpRequest(ctx, "GET", url, endpoint, token, nil)
}

func HttpPost(ctx context.Context, url, endpoint, token string, body []byte) error {
	_, err := HttpRequest(ctx, "POST", url, endpoint, token, bytes.NewReader(body))
	return err
}