package database

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/kiasuo/bot/internal/crypto"
	"github.com/redis/go-redis/v9"
)

type User struct {
	db                 *DB
	TelegramID         int64     `redis:"telegram_id"`
	AccessToken        string    `redis:"access_token"`
	RefreshToken       string    `redis:"refresh_token"`
	StudentID          int       `redis:"student_id"`
	StudentNameAcronym string    `redis:"student_name_acronym"`
	State              UserState `redis:"state"`
	Flags              UserFlag  `redis:"flags"`
}

func getUserKey(telegramID int64) string {
	return "users:" + strconv.FormatInt(telegramID, 10)
}

func (db *DB) GetUser(ctx context.Context, telegramID int64) (*User, error) {
	var user User

	if err := db.client.HGetAll(ctx, getUserKey(telegramID)).Scan(&user); err != nil {
		return nil, err
	}

	if user.IsAnonymous() {
		return nil, nil
	}

	user.db = db

	return &user, nil
}

func (db *DB) NewUser(ctx context.Context, telegramID int64) (*User, error) {
	user := User{
		TelegramID: telegramID,
		State:      UserStatePending,
	}

	if err := db.client.HSet(ctx, getUserKey(telegramID), user).Err(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *User) IsAnonymous() bool {
	return u.TelegramID == 0
}

func (u *User) Save(ctx context.Context, keys ...string) error {
	v := reflect.ValueOf(u).Elem()
	t := v.Type()

	values := make([]interface{}, 0, len(keys)*2)

	for _, k := range keys {
		field, ok := t.FieldByName(k)

		if !ok {
			return fmt.Errorf("field %q does not exist in User", k)
		}

		redisKey := field.Tag.Get("redis")

		if redisKey == "" {
			return fmt.Errorf("field %q does not have a 'redis:...' tag", k)
		}

		fieldValue := v.FieldByName(k)
		values = append(values, redisKey, fieldValue.Interface())
	}

	return u.db.client.HSet(ctx, getUserKey(u.TelegramID), values...).Err()
}

func (u *User) SetState(ctx context.Context, state UserState) error {
	u.State = state

	return u.Save(ctx, "State")
}

func (u *User) SetStudent(ctx context.Context, studentID int, studentNameAcronym string) error {
	u.StudentID = studentID
	u.StudentNameAcronym = crypto.Encrypt(studentNameAcronym).Encrypted

	return u.Save(ctx, "StudentID", "StudentNameAcronym")
}

func (u *User) SetToken(ctx context.Context, accessToken, refreshToken string) error {
	u.AccessToken = crypto.Encrypt(accessToken).Encrypted
	u.RefreshToken = crypto.Encrypt(refreshToken).Encrypted

	return u.Save(ctx, "AccessToken", "RefreshToken")
}

func (u *User) SetFlag(ctx context.Context, flag UserFlag, value bool) error {
	if value {
		u.Flags |= flag
	} else {
		u.Flags &^= flag
	}

	return u.Save(ctx, "Flags")
}

func (u *User) HasFlag(flag UserFlag) bool {
	return u.Flags&flag != 0
}

// GetAccessToken TODO: somehow remove ts
func (u *User) GetAccessToken() string {
	return (&crypto.Crypt{Encrypted: u.AccessToken}).Decrypt()
}

func (u *User) GetRefreshToken() string {
	return (&crypto.Crypt{Encrypted: u.RefreshToken}).Decrypt()
}

func (u *User) GetStudentNameAcronym() string {
	return (&crypto.Crypt{Encrypted: u.StudentNameAcronym}).Decrypt()
}

func getLastMarksCommandKey(telegramID int64, studyPeriodID int) string {
	return getUserKey(telegramID) + ":marks_command:" + strconv.FormatInt(int64(studyPeriodID), 10)
}

func (u *User) GetLastMarksCommand(ctx context.Context, studyPeriodID int) (time.Time, error) {
	var t time.Time

	if err := u.db.client.Get(ctx, getLastMarksCommandKey(u.TelegramID, studyPeriodID)).Scan(&t); err != nil {
		if errors.Is(err, redis.Nil) {
			return t, nil
		}

		return t, err
	}

	return t, nil
}

func (u *User) SetLastMarksCommand(ctx context.Context, studyPeriodID int, t time.Time) error {
	return u.db.client.Set(ctx, getLastMarksCommandKey(u.TelegramID, studyPeriodID), t.Format(time.RFC3339), 0).Err()
}

func (u *User) Delete(ctx context.Context) error {
	if err := u.db.client.Del(ctx, getUserKey(u.TelegramID)).Err(); err != nil {
		return err
	}

	keys, err := u.db.client.Keys(ctx, getUserKey(u.TelegramID)+":").Result()

	if err != nil {
		return err
	}

	for _, key := range keys {
		_ = u.db.client.Del(ctx, key)
	}

	return nil
}
