package users

import (
	"context"
	"github.com/kiasuo/bot/internal/helpers"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserState int

const (
	Unknown UserState = iota
	Ready
	Pending
	Blacklisted
)

type User struct {
	ID                 primitive.ObjectID `bson:"_id"`
	TelegramID         int64              `bson:"telegramID"`
	DiscordID          string             `bson:"discordID,omitempty"`
	AccessToken        string             `bson:"accessToken"`
	RefreshToken       string             `bson:"refreshToken"`
	StudentID          int                `bson:"studentID"`
	StudentNameAcronym string             `bson:"studentNameAcronym"`
	State              UserState          `bson:"state"`
}

var collection *mongo.Collection

func init() {
	uri := helpers.GetEnv("DATABASE")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))

	if err != nil {
		panic(err)
	}

	collection = client.Database("app").Collection("users")
}

func get(filter bson.D) *User {
	var user User

	err := collection.FindOne(context.TODO(), filter).Decode(&user)

	if err != nil {
		return nil
	}

	return &user
}

func GetByID(id string) *User {
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil
	}

	return get(bson.D{{Key: "_id", Value: objectID}})
}

func GetByTelegramID(telegramID int64) *User {
	return get(bson.D{{Key: "telegramID", Value: telegramID}})
}

func GetByDiscordID(discordID string) *User {
	return get(bson.D{{Key: "discordID", Value: discordID}})
}

func update(user User, update bson.D) {
	_, err := collection.UpdateOne(context.TODO(), bson.D{{Key: "_id", Value: user.ID}}, bson.D{{Key: "$set", Value: update}})

	if err != nil {
		panic(err)
	}
}

func (u User) UpdateToken(accessToken string, refreshToken string) {
	update(u, bson.D{{Key: "accessToken", Value: accessToken}, {Key: "refreshToken", Value: refreshToken}})
}

func (u User) UpdateState(state UserState) {
	update(u, bson.D{{Key: "state", Value: state}})
}

func (u User) UpdateStudent(studentID int, studentNameAcronym string) {
	update(u, bson.D{{Key: "studentID", Value: studentID}, {Key: "studentNameAcronym", Value: studentNameAcronym}})
}

func (u User) UpdateDiscord(discordID string) {
	update(u, bson.D{{Key: "discordID", Value: discordID}})
}
