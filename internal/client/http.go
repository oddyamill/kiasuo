package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
)

const (
	cacheHeader = "Worker-Cache"
	authHeader  = "Worker-Authorization"
)

var workerAuth string

func init() {
	workerAuth = os.Getenv("WORKER_AUTH")
}

func request[T any](req *http.Request) (*T, error) {
	if workerAuth != "" {
		req.Header.Set(authHeader, workerAuth)
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	switch res.StatusCode {
	case http.StatusUnauthorized:
		return nil, ErrInvalidToken
	case http.StatusInternalServerError:
		return nil, ErrServerError
	case http.StatusNoContent:
		return nil, nil
	case http.StatusOK:
		break
	default:
		return nil, errors.New(res.Status)
	}

	var result *T

	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
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
		if err = refreshToken(client); err != nil {
			return nil, err
		}
	}

	if !client.User.Cache {
		req.Header.Set(cacheHeader, "no")
	}

	result, err := requestWithAuth[T](client, req)

	if err == nil {
		return result, nil
	}

	if errors.Is(err, ErrInvalidToken) {
		if err = refreshToken(client); err != nil {
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

func requestPurgeCache(id *int) bool {
	req, err := http.NewRequest("POST", appendID(purgeCacheURL, id), nil)

	if err != nil {
		return false
	}

	req.Header.Set(authHeader, workerAuth)
	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return false
	}

	defer res.Body.Close()
	return res.StatusCode == http.StatusOK
}
