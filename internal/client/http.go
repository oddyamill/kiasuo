package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

var workerAuth string

func init() {
	workerAuth = os.Getenv("WORKER_AUTH")
}

func request[T any](req *http.Request) (*T, error) {
	if workerAuth != "" {
		req.Header.Set("Worker-Authorization", workerAuth)
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusUnauthorized {
		return nil, ErrInvalidToken
	}

	if res.StatusCode == http.StatusInternalServerError {
		return nil, ErrServerError
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New(res.Status)
	}

	var result *T
	err = json.NewDecoder(res.Body).Decode(&result)

	if err != nil {
		if err.Error() == "EOF" {
			return nil, nil
		}

		return nil, err
	}

	return result, nil
}

func refreshToken(client *Client) error {
	req, err := http.NewRequest("POST", refreshURL, bytes.NewBufferString(
		"refresh-token="+client.User.RefreshToken.Decrypt(),
	))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		return err
	}

	result, err := request[Token](req)

	if err != nil {
		if errors.Is(err, ErrInvalidToken) {
			return ErrExpiredToken
		}

		return err
	}

	client.User.UpdateToken(result.AccessToken, result.RefreshToken)
	return nil
}

func requestWithAuth[T any](client *Client, req *http.Request) (*T, error) {
	req.Header.Set("Authorization", "Bearer "+client.User.AccessToken.Decrypt())
	return request[T](req)
}

func requestWithClient[T any](client *Client, url string, method string) (*T, error) {
	req, err := http.NewRequest(method, appendID(url, client.User.StudentID), nil)

	if err != nil {
		return nil, err
	}

	if client.isTokenExpired() {
		err = refreshToken(client)

		if err != nil {
			return nil, err
		}
	}

	result, err := requestWithAuth[T](client, req)

	if err == nil {
		return result, nil
	}

	if errors.Is(err, ErrInvalidToken) {
		err = refreshToken(client)

		if err != nil {
			return nil, err
		}

		result, err = requestWithAuth[T](client, req)

		if err != nil {
			return nil, err
		}

		return result, nil
	}

	return nil, err
}
