package main

import (
	"context"
	"time"
)

func main() {
	loadConfig()
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	database := connect(ctx)

	yesterday, today := period()
	activities := lookupActivities(ctx, database, yesterday, today)
	users := lookupUsers(ctx, database)
	boards := lookupBoards(ctx, database)
	msgs := loadMessages(activities, users)
	getMail(msgs, yesterday, today, users, boards)
}

func period() (time.Time, time.Time) {
	today := time.Now().Truncate(24 * time.Hour)
	yesterday := time.Now().Truncate(24 * time.Hour).Add(-24 * 30 * time.Hour)
	return yesterday, today
}
