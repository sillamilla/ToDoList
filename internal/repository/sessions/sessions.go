package sessions

import (
	"ToDoWithKolya/internal/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Repo interface {
	Upsert(ctx context.Context, userID, session string, now time.Time) error
	GetUserID(ctx context.Context, session string) (string, error)

	Delete(ctx context.Context, session string) error
	SessionInfo(ctx context.Context, session string) (models.SessionInfo, error)
}

type sessions struct {
	db *mongo.Collection
}

func New(database *mongo.Database) Repo {
	return sessions{
		db: database.Collection("sessions"),
	}
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

func (r sessions) SessionInfo(ctx context.Context, session string) (models.SessionInfo, error) {
	var sessionInfo models.SessionInfo

	filter := bson.M{"session": session}

	err := r.db.FindOne(ctx, filter).Decode(&sessionInfo)
	if err != nil {
		return models.SessionInfo{}, err
	}

	return sessionInfo, nil
}

func (r sessions) Delete(ctx context.Context, session string) error {
	_, err := r.db.DeleteMany(ctx, bson.M{"session": session})
	if err != nil {
		return err
	}

	return nil
}
