package users_sql

import (
	"database/sql"
	"errors"
	"github.com/kiasuo/bot/internal/helpers"
	_ "github.com/lib/pq"
)

type UserState int

const (
	Unknown UserState = iota
	Ready
	Pending
	Blacklisted
)

type User struct {
	ID                 int
	TelegramID         int64
	DiscordID          string
	AccessToken        string
	RefreshToken       string
	StudentID          int
	StudentNameAcronym string
	State              UserState
}

var db *sql.DB

func init() {
	uri := "user=" + helpers.GetEnv("POSTGRES_USER") +
		" dbname=" + helpers.GetEnv("POSTGRES_DB") +
		" password=" + helpers.GetEnv("POSTGRES_PASSWORD") +
		" sslmode=disable"

	var err error
	db, err = sql.Open("postgres", uri)

	if err != nil {
		panic(err)
	}
}

func createTable() {
	query(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			telegram_id BIGINT NOT NULL UNIQUE,
			discord_id TEXT UNIQUE,
			access_token TEXT,
			refresh_token TEXT,
			student_id INTEGER,
			student_name_acronym TEXT,
			state INTEGER NOT NULL
		)
	`)
}

func query(query string, args ...any) {
	_, err := db.Query(query, args...)

	if err != nil {
		panic(err)
	}
}

func queryRow(query string, args ...any) *User {
	rows := db.QueryRow(query, args...)

	var user User

	err := rows.Scan(&user.ID, &user.TelegramID, &user.DiscordID, &user.AccessToken, &user.RefreshToken, &user.StudentID, &user.StudentNameAcronym, &user.State)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		panic(err)
	}

	return &user
}

func GetByID(id string) *User {
	return queryRow("SELECT * FROM users WHERE id = $1", id)
}

func GetByTelegramID(telegramID int64) *User {
	return queryRow("SELECT * FROM users WHERE telegram_id = $1", telegramID)
}

func GetByDiscordID(discordID string) *User {
	return queryRow("SELECT * FROM users WHERE discord_id = $1", discordID)
}

func UpdateToken(user User, accessToken string, refreshToken string) {
	query("UPDATE users SET access_token = $1, refresh_token = $2 WHERE id = $3", accessToken, refreshToken, user.ID)
}

func UpdateState(user User, state UserState) {
	query("UPDATE users SET state = $1 WHERE id = $2", state, user.ID)
}

func UpdateStudent(user User, studentID int, studentNameAcronym string) {
	query("UPDATE users SET student_id = $1, student_name_acronym = $2 WHERE id = $3", studentID, studentNameAcronym, user.ID)
}

func UpdateDiscord(user User, discordID string) {
	query("UPDATE users SET discord_id = $1 WHERE id = $2", discordID, user.ID)
}
