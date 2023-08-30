package sessions

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Repo interface {
	Upsert(ctx context.Context, userID, session string, now time.Time) error
	GetUserID(ctx context.Context, session string) (string, error)
	GetSessionTime(ctx context.Context, session string) (time.Time, error)

	Delete(ctx context.Context, session string) error
}

type sessions struct {
	db *mongo.Collection
}

func New(database *mongo.Database) Repo {
	return sessions{
		db: database.Collection("sessions"),
	}
}

type Session struct {
	SessionID string    `bson:"session"`
	UserID    string    `bson:"user_id"`
	CreatedAt time.Time `bson:"created_at"`
}

func (r sessions) Upsert(ctx context.Context, userID, session string, now time.Time) error {
	filter := bson.M{"user_id": userID}
	update := bson.M{"$set": bson.M{"session": session, "created_at": now}}
	opts := options.Update().SetUpsert(true)

	_, err := r.db.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	return nil
}

func (r sessions) GetUserID(ctx context.Context, session string) (string, error) {
	var id string

	err := r.db.FindOne(ctx, bson.M{"session": session}).Decode(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (r sessions) GetSessionTime(ctx context.Context, session string) (time.Time, error) {
	var sessionInfo Session
	filter := bson.M{"session": session}

	err := r.db.FindOne(ctx, filter).Decode(&sessionInfo)
	if err != nil {
		return time.Time{}, err
	}

	return sessionInfo.CreatedAt, nil
}

func (r sessions) Delete(ctx context.Context, id string) error {
	_, err := r.db.DeleteMany(ctx, bson.M{"user_id": id})
	if err != nil {
		return err
	}

	return nil
}
