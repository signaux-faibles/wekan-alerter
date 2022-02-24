package main

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type activity struct {
	ActivityType string `bson:"activityType"`
	BoardId      string `bson:"boardId"`
	ListId       string `bson:"listId"`
	UserId       string `bson:"userId"`
	Card         card   `bson:"card"`
}

type card struct {
	ID      string   `bson:"_id"`
	Members []string `bson:"members"`
	UserID  string   `bson:"userId"`
	Title   string   `bson:"title"`
}

type board struct {
	ID     string        `bson:"_id"`
	Title  string        `bson:"title"`
	Labels []interface{} `bson:"labels"`
	Slug   string        `bson:"slug"`
}

type user struct {
	ID       string `bson:"_id"`
	Username string `bson:"username"`
	Services struct {
		OIDC struct {
			Email string `bson:"email"`
		} `bson:"oidc"`
	} `bson:"services"`
	Profile struct {
		Fullname string `bson:"fullname"`
	} `bson:"profile"`
	Emails []userEmail `bson:"emails"`
}

type userEmail struct {
	Address  string `bson:"address"`
	Verified bool   `bson:"verified"`
}

func connect(ctx context.Context) *mongo.Database {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(MONGO))
	if err != nil {
		panic(err)
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		panic(err)
	}
	return client.Database(DB)
}

func lookupActivities(ctx context.Context, database *mongo.Database, from time.Time, to time.Time) []activity {
	cursor, err := database.Collection("activities").Aggregate(ctx, getActivitiesPipeline(from, to))
	if err != nil {
		panic(err)
	}
	var activities = make([]activity, 0)
	err = cursor.All(ctx, &activities)
	if err != nil {
		panic(err)
	}
	return activities
}

func lookupBoards(ctx context.Context, database *mongo.Database) map[string]board {
	var boards []board
	cursor, err := database.Collection("boards").Find(ctx, bson.M{})
	if err != nil {
		panic(err)
	}
	err = cursor.All(ctx, &boards)
	if err != nil {
		panic(err)
	}
	mapBoards := make(map[string]board)
	for _, b := range boards {
		mapBoards[b.ID] = b
	}
	return mapBoards
}

func lookupUsers(ctx context.Context, database *mongo.Database) map[string]user {
	var users []user
	cursor, err := database.Collection("users").Find(ctx, bson.M{
		"$or": bson.A{
			bson.M{"loginDisabled": false},
			bson.M{"loginDisabled": bson.M{"$exists": false}},
		},
	})
	if err != nil {
		panic(err)
	}
	err = cursor.All(ctx, &users)
	if err != nil {
		panic(err)
	}
	mapUsers := make(map[string]user)
	for _, u := range users {
		mapUsers[u.ID] = u
	}
	return mapUsers
}

func getActivitiesPipeline(from time.Time, to time.Time) bson.A {
	return bson.A{
		bson.M{"$match": bson.M{
			"createdAt": bson.M{
				"$gte": from,
				"$lt":  to,
			},
		}},
		bson.M{"$lookup": bson.M{
			"from":         "cards",
			"localField":   "cardId",
			"foreignField": "_id",
			"as":           "card",
		}},
		bson.M{"$unwind": "$card"},
	}
}
