package client

import (
	"encoding/base64"
	"encoding/json"
	"github.com/kiasuo/bot/internal/users"
	"strings"
	"time"
)

type Client struct {
	User *users.User
}

func (c *Client) RefreshToken() error {
	return refreshToken(c)
}

func (c *Client) GetUser() (*User, error) {
	return requestWithClient[User](c, userURL, "GET")
}

func (c *Client) GetRecipients() (*Recipients, error) {
	rawRecipients, err := requestWithClient[RawRecipient](c, recipientsURL, "GET")

	if err != nil {
		return nil, err
	}

	recipients := (*rawRecipients)[c.User.StudentID]
	return &recipients, nil
}

func (c *Client) GetStudyPeriods() (*[]StudyPeriod, error) {
	return requestWithClient[[]StudyPeriod](c, studyPeriodsURL, "GET")
}

func (c *Client) GetLessons(id int) (*[]Lesson, error) {
	rawMarks, err := requestWithClient[RawLessons](c, lessonMarksURL(id), "GET")

	if err != nil {
		return nil, err
	}

	return &rawMarks.Lessons, nil
}

func (c *Client) GetSchedule(time time.Time) (*RawSchedule, error) {
	year, week := time.ISOWeek()

	return requestWithClient[RawSchedule](c, scheduleURL(year, week), "GET")
}

func (c *Client) isTokenExpired() bool {
	segments := strings.Split(c.User.AccessToken.Decrypt(), ".")

	if len(segments) != 3 {
		return true
	}

	raw := segments[1]
	padding := len(raw) % 4

	if padding > 0 {
		raw += strings.Repeat("=", 4-padding)
	}

	plain, err := base64.StdEncoding.DecodeString(raw)

	if err != nil {
		return true
	}

	var tokenPayload struct {
		Exp int `json:"exp"`
	}

	err = json.Unmarshal(plain, &tokenPayload)

	if err != nil {
		return true
	}

	return time.Unix(int64(tokenPayload.Exp), 0).Before(time.Now())
}
