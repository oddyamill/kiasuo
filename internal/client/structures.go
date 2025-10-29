package client

import (
	"strconv"
	"strings"
	"time"

	"github.com/kiasuo/bot/internal/helpers"
)

type Token struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type TokenPayload struct {
	Expiration int64 `json:"exp"`
}

type User struct {
	ID       int     `json:"id"`
	Username string  `json:"username"`
	Parent   bool    `json:"parent"`
	Children []Child `json:"children"`
}

type Child struct {
	ID          int    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	MiddleName  string `json:"middle_name"`
	SchoolClass string `json:"school_class"`
	Age         int    `json:"age"`
}

type RawRecipient map[int]Recipients

type Recipients struct {
	Staff    map[string]map[string]int `json:"staff"`
	Students map[string]Student        `json:"students"`
}

type Student struct {
	Parents any  `json:"parents"`
	ID      *int `json:"id"`
}

type StudyPeriod struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
	From string `json:"from"`
	To   string `json:"to"`
}

func (p StudyPeriod) Match(t time.Time) bool {
	from, _ := time.Parse(time.DateOnly, p.From)

	if t.After(from) {
		to, _ := time.Parse(time.DateOnly, p.To)

		if t.Before(to) {
			return true
		}
	}

	return false
}

type RawLessons struct {
	Lessons []Lesson `json:"lessons"`
}

type Lesson struct {
	Subject  string         `json:"subject"`
	Slots    []Slot         `json:"slots"`
	Averages LessonAverages `json:"averages"`
}

func (l Lesson) String() string {
	return helpers.HumanizeLesson(l.Subject)
}

type LessonAverages struct {
	ForStudentRaw *[2]interface{} `json:"for_student"`
	ForClassRaw   *[2]interface{} `json:"for_class"`
}

func (l LessonAverages) ForStudent() string {
	return l.ForStudentRaw[0].(string)
}

func (l LessonAverages) ForStudentValue() int {
	return int(l.ForStudentRaw[1].(float64))
}

func (l LessonAverages) ForClass() string {
	return l.ForClassRaw[0].(string)
}

func (l LessonAverages) ForClassValue() int {
	return int(l.ForClassRaw[1].(float64))
}

type Slot struct {
	LessonDate string    `json:"lesson_date"`
	Mark       Mark      `json:"mark"`
	Text       string    `json:"text"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (s Slot) Date() string {
	// 2025-09-05
	return s.LessonDate[8:10] + "." + s.LessonDate[5:7]
}

type Mark struct {
	Value string `json:"value"`
}

func (m Mark) IsPass() bool {
	return m.Value == "Б" || m.Value == "Н" || m.Value == "У"
}

type RawSchedule struct {
	Schedule  []Event    `json:"schedule"`
	Homeworks []Homework `json:"homeworks"`
}

type Event struct {
	Subject    string `json:"subject"`
	LessonDate string `json:"lesson_date"`
	Number     int    `json:"lesson_number"`
	Homeworks  []int  `json:"homework_to_check_ids"`
	Slots      []Slot `json:"slots"`
}

func (e Event) Date() time.Time {
	date, _ := time.Parse(time.DateOnly, e.LessonDate)
	return date
}

func (e Event) String() string {
	return strconv.Itoa(e.Number) + ". " + helpers.HumanizeLesson(e.Subject)
}

type Homework struct {
	ID    int    `json:"id"`
	Text  string `json:"text"`
	Files []File `json:"files"`
	Links []Link `json:"links"`
}

func (h Homework) String() string {
	if h.Text == "Без задания" {
		return ""
	}

	return strings.TrimSpace(h.Text)
}

type File struct {
	Url   string `json:"url"`
	Title string `json:"title"`
}

func (f File) String(formatter helpers.Formatter) string {
	return formatter.Link(f.Title, PublicUrl+f.Url)
}

type Link struct {
	Url   string `json:"url"`
	Title string `json:"title"`
}

func (l Link) String(formatter helpers.Formatter) string {
	return formatter.Link(l.Title, l.Url)
}
