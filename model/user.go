package model

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"lamartire/cryptoratesbot/storage"
)

type UserMongoStorage struct {
	Collection *mongo.Collection
}

func (u UserMongoStorage) GetUserByName(name string) (storage.User, error) {
	var user storage.User

	ctx := context.Background()
	err := u.Collection.FindOne(ctx, bson.M{"name": name}).Decode(&user)

	if err != nil && err == mongo.ErrNoDocuments {
		return user, storage.NoResultErr
	}

	if err != nil {
		return user, err
	}

	return user, nil
}

func (u UserMongoStorage) GetNotSuspendedUsers() ([]storage.User, error) {
	var users = make([]storage.User, 0)

	ctx := context.Background()
	cursor, err := u.Collection.Find(ctx, bson.D{{"suspended", false}})

	if err != nil && err == mongo.ErrNoDocuments || cursor == nil {
		return users, storage.NoResultErr
	}

	for cursor.Next(ctx) {
		var user storage.User

		if err = cursor.Decode(&user); err != nil {
			return users, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (u UserMongoStorage) GetUsersForNotification(now int64) ([]storage.User, error) {
	targetUsers := make([]storage.User, 0)
	users, err := u.GetNotSuspendedUsers()

	if err == storage.NoResultErr {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if len(user.Symbols) == 0 {
			continue
		}

		// interval stored in minutes
		interval := user.Interval * 60
		lastReply := user.LastReply

		if now-lastReply < interval {
			continue
		}

		targetUsers = append(targetUsers, user)
	}

	return targetUsers, nil
}

func (u UserMongoStorage) CreateUser(name string, chatID int64) (storage.User, error) {
	user := storage.User{
		Name:      name,
		ChatID:    chatID,
		Suspended: false,
		Interval:  60,
		LastReply: 0,
	}
	ctx := context.Background()
	_, err := u.Collection.InsertOne(ctx, user)

	if err != nil {
		return user, err
	}

	return user, nil
}

func (u UserMongoStorage) UpdateUserByName(name string, user storage.User) error {
	var updatedUser storage.User

	ctx := context.Background()
	err := u.Collection.FindOneAndUpdate(
		ctx,
		bson.M{"name": name},
		bson.D{{
			"$set",
			bson.M{
				"symbols":   user.Symbols,
				"interval":  user.Interval,
				"lastreply": user.LastReply,
				"suspended": user.Suspended,
			},
		}},
	).Decode(&updatedUser)

	if err != nil {
		return err
	}

	return nil
}
