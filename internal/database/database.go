package database

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Database struct {
	client *redis.Client
}

//users:(id) => User
//users:(id):marks_command:(study period id) => time.Time

func New() *Database {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	return &Database{client}
}

func (db *Database) Ping(ctx context.Context) error {
	return db.client.Ping(ctx).Err()
}
