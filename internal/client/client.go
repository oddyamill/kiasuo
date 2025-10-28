package client

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/kiasuo/bot/internal/database"
)

type Client struct {
	User *database.User
}

func New(user *database.User) *Client {
	return &Client{user}
}

func (c *Client) RefreshToken() error {
	return refreshToken(c)
}

func (c *Client) GetUser() (*User, error) {
	return requestWithClient[User](c, userURL, http.MethodGet)
}

func (c *Client) GetRecipients() (*Recipients, error) {
	rawRecipients, err := requestWithClient[RawRecipient](c, recipientsURL, http.MethodGet)

	if err != nil {
		return nil, err
	}

	recipients := (*rawRecipients)[c.User.StudentID]
	return &recipients, nil
}

func (c *Client) GetStudyPeriods() (*[]StudyPeriod, error) {
	return requestWithClient[[]StudyPeriod](c, studyPeriodsURL, http.MethodGet)
}

func (c *Client) GetLessons(id int) (*[]Lesson, error) {
	rawMarks, err := requestWithClient[RawLessons](c, lessonMarksURL(id), http.MethodGet)

	if err != nil {
		return nil, err
	}

	return &rawMarks.Lessons, nil
}

func (c *Client) GetSchedule(time time.Time) (*RawSchedule, error) {
	year, week := time.ISOWeek()

	return requestWithClient[RawSchedule](c, scheduleURL(year, week), http.MethodGet)
}

func (c *Client) RevokeToken() error {
	_, err := requestWithClient[any](c, revokeURL, http.MethodDelete)
	return err
}

func (c *Client) PurgeCache() bool {
	return requestPurgeCache(c.User.StudentID)
}

func (c *Client) isTokenExpired() bool {
	segments := strings.Split(c.User.GetAccessToken(), ".")

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

	var tokenPayload TokenPayload

	if err = json.Unmarshal(plain, &tokenPayload); err != nil {
		return true
	}

	return time.Unix(tokenPayload.Expiration, 0).Before(time.Now())
}
