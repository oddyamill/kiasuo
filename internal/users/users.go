package users

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/kiasuo/bot/internal/crypto"
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
	AccessToken        crypto.Crypt
	RefreshToken       crypto.Crypt
	StudentID          int
	StudentNameAcronym crypto.Crypt
	State              UserState
	LastMarksUpdate    time.Time
	Cache              bool
	VkCookie           crypto.Crypt
}

var db *sql.DB

func init() {
	if helpers.IsTesting() {
		return
	}

	uri := "host=" + helpers.GetEnv("POSTGRES_HOST") +
		" user=" + helpers.GetEnv("POSTGRES_USER") +
		" dbname=" + helpers.GetEnv("POSTGRES_DB") +
		" password=" + helpers.GetEnv("POSTGRES_PASSWORD") +
		" sslmode=disable"

	var err error
	db, err = sql.Open("postgres", uri)

	if err != nil {
		panic(err)
	}

	log.Println("Connected to database")
	createTable()
	createIndex()
	migrate()
}

func createTable() {
	query(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			telegram_id BIGINT NOT NULL UNIQUE,
			discord_id TEXT UNIQUE,
			access_token TEXT,
			refresh_token VARCHAR(96),
			student_id INTEGER,
			student_name_acronym TEXT,
			state INTEGER NOT NULL,
		  last_marks_update TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		  cache BOOLEAN DEFAULT TRUE,
			vk_cookie TEXT
		)
	`)
}

func createIndex() {
	query("CREATE INDEX IF NOT EXISTS telegram_id_index ON users (telegram_id)")
	query("CREATE INDEX IF NOT EXISTS discord_id_index ON users (discord_id)")
}

func migrate() {
	query("ALTER TABLE users ADD COLUMN IF NOT EXISTS cache BOOLEAN DEFAULT TRUE")
	query("ALTER TABLE users ADD COLUMN IF NOT EXISTS vk_cookie TEXT")
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

	err := rows.Scan(
		&user.ID,
		&user.TelegramID,
		&user.DiscordID,
		&user.AccessToken,
		&user.RefreshToken,
		&user.StudentID,
		&user.StudentNameAcronym,
		&user.State,
		&user.LastMarksUpdate,
		&user.Cache,
		&user.VkCookie,
	)

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

func (u *User) IsReady() bool {
	return u.State == Ready
}

func (u *User) UpdateToken(accessToken, refreshToken string) {
	u.AccessToken = crypto.Encrypt(accessToken)
	u.RefreshToken = crypto.Encrypt(refreshToken)

	query(
		"UPDATE users SET access_token = $1, refresh_token = $2 WHERE id = $3",
		u.AccessToken.Encrypted,
		u.RefreshToken.Encrypted,
		u.ID,
	)
}

func (u *User) UpdateState(state UserState) {
	query("UPDATE users SET state = $1 WHERE id = $2", state, u.ID)
}

func (u *User) UpdateStudent(studentID int, studentNameAcronym string) {
	query(
		"UPDATE users SET student_id = $1, student_name_acronym = $2 WHERE id = $3",
		studentID,
		crypto.Encrypt(studentNameAcronym).Encrypted,
		u.ID,
	)
}

func (u *User) UpdateDiscord(discordID string) {
	query("UPDATE users SET discord_id = $1 WHERE id = $2", discordID, u.ID)
}

func (u *User) UpdateLastMarksUpdate() {
	query("UPDATE users SET last_marks_update = CURRENT_TIMESTAMP WHERE id = $1", u.ID)
}

func (u *User) UpdateCache(cache bool) {
	query("UPDATE users SET cache = $1 WHERE id = $2", cache, u.ID)
}

func (u *User) UpdateVkCookie(cookie string) {
	query("UPDATE users SET vk_cookie = $1 WHERE id = $2", cookie, u.ID)
}

func (u *User) Delete() {
	query("DELETE FROM users WHERE id = $1", u.ID)
}
