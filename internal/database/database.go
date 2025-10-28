package database

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
)

type DB struct {
	client *redis.Client
}

//users:(id) => User
//users:(id):marks_command:(study period id) => time.Time

func New(url string) *DB {
	options, err := redis.ParseURL(url)

	if err != nil {
		log.Panic(err)
	}

	return &DB{redis.NewClient(options)}
}

func (db *DB) Ping(ctx context.Context) error {
	return db.client.Ping(ctx).Err()
}
