package sessions

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Sessions interface {
	Create(ctx context.Context, userID, session string) error
	GetUserID(ctx context.Context, session string) (string, error)
	Upsert(ctx context.Context, userID, session string) error
	GetSessionTime(ctx context.Context, session string) (time.Time, error)
	Delete(ctx context.Context, session string) error
}

type sessions struct {
	db *mongo.Collection
}

func New(database *mongo.Database) Sessions {
	return sessions{
		db: database.Collection("sessions"),
	}
}

func (r sessions) Create(ctx context.Context, userID, session string) error {
	sessionDoc := bson.M{"user_id": userID, "session": session, "created_at": time.Now()}

	_, err := r.db.InsertOne(ctx, sessionDoc)
	if err != nil {
		return err
	}

	return nil
}

func (r sessions) Upsert(ctx context.Context, userID, session string) error {
	filter := bson.M{"user_id": userID}
	update := bson.M{"$set": bson.M{"session": session, "created_at": time.Now()}}
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

func (r sessions) Delete(ctx context.Context, id string) error {
	_, err := r.db.DeleteMany(ctx, bson.M{"user_id": id})
	if err != nil {
		return err
	}

	return nil
}

func (r sessions) GetSessionTime(ctx context.Context, session string) (time.Time, error) {
	var last time.Time

	err := r.db.FindOne(ctx, bson.M{"session": session}).Decode(&last)
	if err != nil {
		return time.Time{}, err
	}

	return last, nil
}
