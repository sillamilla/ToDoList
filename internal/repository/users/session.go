package users

import (
	"ToDoWithKolya/internal/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

func (r userRepo) CreateSession(userID string, session string) error {
	sessionDoc := bson.M{"user_id": userID, "session": session}
	_, err := r.db.InsertOne(context.TODO(), sessionDoc)
	return err
}

func (r userRepo) UpsertSession(userID string, session string) error {
	filter := bson.M{"user_id": userID}
	update := bson.M{"$set": bson.M{"session": session}}
	opts := options.Update().SetUpsert(true)
	_, err := r.db.UpdateOne(context.TODO(), filter, update, opts)
	return err
}

func (r userRepo) GetUserBySession(session string) (models.User, error) {
	var user models.User
	lookupStage := bson.D{{"$lookup", bson.D{
		{"from", "users"},
		{"localField", "user_id"},
		{"foreignField", "id"},
		{"as", "user"},
	}}}
	matchStage := bson.D{{"$match", bson.D{{"session", session}}}}
	unwindStage := bson.D{{"$unwind", "$user"}}
	projectStage := bson.D{{"$project", bson.D{
		{"_id", 0},
		{"user.id", 1},
		{"user.login", 1},
		{"user.password", 1},
		{"user.email", 1},
	}}}

	cursor, err := r.db.Aggregate(context.TODO(), mongo.Pipeline{
		lookupStage, matchStage, unwindStage, projectStage,
	})
	if err != nil {
		return models.User{}, err
	}
	defer cursor.Close(context.TODO())

	if cursor.Next(context.TODO()) {
		if err := cursor.Decode(&user); err != nil {
			return models.User{}, err
		}
		return user, nil
	}

	return models.User{}, mongo.ErrNoDocuments
}

func (r userRepo) DeleteSession(session string) error {
	_, err := r.db.DeleteMany(context.TODO(), bson.M{"session": session})
	return err
}

func (r userRepo) GetSessionLastActive(session string) (time.Time, error) {
	var sessionTime time.Time
	projection := options.FindOne().SetProjection(bson.M{"created_at": 1})
	filter := bson.M{"session": session}

	result := r.db.FindOne(context.TODO(), filter, projection)
	if result.Err() != nil {
		if result.Err() == mongo.ErrNoDocuments {
			return time.Time{}, models.DBErr(models.ErrNotFound)
		}
		return time.Time{}, result.Err()
	}

	if err := result.Decode(&sessionTime); err != nil {
		return time.Time{}, err
	}

	return sessionTime, nil
}
