package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/kiasuo/bot/internal/helpers"
	"github.com/kiasuo/bot/internal/users"
	"net/http"
	"strconv"
	"time"
)

const (
	PublicUrl = "https://dnevnik.kiasuo.ru"
	ApiUrl    = "https://kiasuo-proxy.oddya.ru/diary"
)

type Client struct {
	User *users.User
}

func httpRequest[T any](client Client, request *http.Request) (*http.Response, *T, error) {
	request.Header.Set("Authorization", "Bearer "+client.User.AccessToken.Decrypt())
	response, err := http.DefaultClient.Do(request)

	if err != nil {
		return response, nil, err
	}

	defer response.Body.Close()

	var result *T
	err = json.NewDecoder(response.Body).Decode(&result)

	if err != nil {
		if err.Error() == "EOF" {
			return response, nil, nil
		}

		return response, nil, err
	}

	return response, result, nil
}

func RefreshToken(client *Client) error {
	body := helpers.StringToBytes(`{"refresh-token":"` + client.User.RefreshToken.Decrypt() + `"}`)

	request, err := http.NewRequest("POST", ApiUrl+"/refresh", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")

	if err != nil {
		return err
	}

	response, result, err := httpRequest[Token](*client, request)

	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		return errors.New(response.Status)
	}

	if result == nil {
		return errors.New("empty response")
	}

	client.User.UpdateToken(result.AccessToken, result.RefreshToken)
	return nil
}

func requestWithClient[T any](client *Client, pathname string, method string) (*T, error) {
	request, err := http.NewRequest(method, ApiUrl+pathname, nil)

	if err != nil {
		return nil, err
	}

	if client.User.IsTokenExpired() {
		err = RefreshToken(client)

		if err != nil {
			return nil, err
		}
	}

	response, result, err := httpRequest[T](*client, request)

	if err != nil {
		return nil, err
	}

	if response.StatusCode == http.StatusUnauthorized {
		err = RefreshToken(client)

		if err != nil {
			return nil, err
		}

		_, result, err = httpRequest[T](*client, request)

		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (c *Client) GetUser() (*User, error) {
	// TODO:
	return requestWithClient[User](c, "/api/user?id="+strconv.Itoa(c.User.StudentID), "GET")
}

func (c *Client) GetRecipients() (*Recipients, error) {
	// TODO:
	rawRecipients, err := requestWithClient[RawRecipient](c, "/api/recipients?id="+strconv.Itoa(c.User.StudentID), "GET")

	if err != nil {
		return nil, err
	}

	recipients := (*rawRecipients)[c.User.StudentID]
	return &recipients, nil
}

func (c *Client) GetStudyPeriods() (*[]StudyPeriod, error) {
	// TODO:
	return requestWithClient[[]StudyPeriod](c, "/api/study_periods?id="+strconv.Itoa(c.User.StudentID), "GET")
}

func (c *Client) GetLessons(id int) (*[]Lesson, error) {
	rawMarks, err := requestWithClient[RawLessons](c, "/api/lesson_marks/"+strconv.Itoa(id), "GET")

	if err != nil {
		return nil, err
	}

	return &rawMarks.Lessons, nil
}

func (c *Client) GetSchedule(time time.Time) (*RawSchedule, error) {
	year, week := time.ISOWeek()

	return requestWithClient[RawSchedule](c, "/api/schedule?year="+strconv.Itoa(year)+"&week="+strconv.Itoa(week), "GET")
}
