package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/kiasuo/bot/users"
	"net/http"
)

const BaseUrl = "https://kiasuo-proxy.oddya.ru/diary"

type Client struct {
	User users.User
}

func NewClient(user users.User) Client {
	return Client{User: user}
}

func Request[T any](client Client, request *http.Request) (*http.Response, *T, error) {
	request.Header.Set("Authorization", "Bearer "+client.User.AccessToken)
	response, err := http.DefaultClient.Do(request)

	println(response.StatusCode, response.Status, request.URL)

	if err != nil {
		return response, nil, err
	}

	defer response.Body.Close()

	var result *T
	err = json.NewDecoder(response.Body).Decode(&result)

	if err != nil {
		return response, nil, err
	}

	return response, result, nil
}

func RefreshToken(client Client) error {
	body := []byte(`{"refresh-token":"` + client.User.RefreshToken + `"}`)

	request, err := http.NewRequest("POST", BaseUrl+"/refresh", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")

	if err != nil {
		return err
	}

	response, result, err := Request[Token](client, request)

	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return errors.New(response.Status)
	}

	if result == nil {
		return errors.New("empty response")
	}

	users.UpdateToken(client.User, result.AccessToken, result.RefreshToken)
	client.User.AccessToken = result.AccessToken
	client.User.RefreshToken = result.RefreshToken
	return nil
}

func ClientRequest[T any](client Client, pathname string, method string) (*T, error) {
	request, err := http.NewRequest(method, BaseUrl+pathname, nil)

	if err != nil {
		return nil, err
	}

	response, result, err := Request[T](client, request)

	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusUnauthorized {
		err = RefreshToken(client)

		if err != nil {
			return nil, err
		}

		_, result, err = Request[T](client, request)

		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (client Client) GetUser() (*User, error) {
	return ClientRequest[User](client, "/api/user", "GET")
}

func (client Client) GetRecipients() (*Recipients, error) {
	rawRecipients, err := ClientRequest[RawRecipient](client, "/api/recipients", "GET")

	if err != nil {
		return nil, err
	}

	recipients := (*rawRecipients)[client.User.StudentID]
	return &recipients, nil
}
