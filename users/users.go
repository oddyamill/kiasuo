package users

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
)

type UserState int

const (
	Unknown UserState = iota
	Ready
	Pending
	Blacklisted
)

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	TelegramID   int64              `bson:"userId"`
	DiscordID    string             `bson:"discordId"`
	AccessToken  string             `bson:"accessToken"`
	RefreshToken string             `bson:"refreshToken"`
	StudentID    int                `bson:"studentId"`
	State        UserState          `bson:"state"`
}

var collection *mongo.Collection

func init() {
	uri, ok := os.LookupEnv("DATABASE")

	if !ok {
		panic("DATABASE not set")
	}

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}

	collection = client.Database("app").Collection("users")
}

func Get(filter bson.D) *User {
	var user User

	err := collection.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		return nil
	}

	return &user
}

func GetByTelegramID(telegramID int64) *User {
	return Get(bson.D{{Key: "userId", Value: telegramID}})
}

func GetByDiscordID(discordID string) *User {
	return Get(bson.D{{Key: "discordId", Value: discordID}})
}

func Update(user User, update bson.D) {
	_, err := collection.UpdateOne(context.TODO(), bson.D{{Key: "_id", Value: user.ID}}, bson.D{{Key: "$set", Value: update}})

	if err != nil {
		panic(err)
	}
}

func UpdateToken(user User, accessToken string, refreshToken string) {
	Update(user, bson.D{{Key: "accessToken", Value: accessToken}, {Key: "refreshToken", Value: refreshToken}})
}
