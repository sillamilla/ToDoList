package tasks

import (
	models "ToDoWithKolya/internal/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Tasks interface {
	Create(ctx context.Context, task models.Task) error
	Update(ctx context.Context, task models.Task) error
	GetAll(ctx context.Context, userID string) ([]models.Task, error)
	Get(ctx context.Context, userID, id string) (models.Task, error)
	MarkValueSet(ctx context.Context, taskID string, status int) error
	Search(ctx context.Context, taskName, userID string) ([]models.Task, error)
	Delete(ctx context.Context, userID, id string) error
}

type tasks struct {
	db *mongo.Collection
}

func New(database *mongo.Database) Tasks {
	return tasks{
		db: database.Collection("tasks"),
	}
}

func (r tasks) Create(ctx context.Context, task models.Task) error {
	_, err := r.db.InsertOne(ctx, task)
	if err != nil {
		return err
	}

	return nil
}

func (r tasks) Delete(ctx context.Context, userID, id string) error {
	filter := bson.M{"user_id": userID, "id": id}

	_, err := r.db.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (r tasks) Update(ctx context.Context, task models.Task) error {
	filter := bson.M{"user_id": task.UserID, "id": task.ID}
	update := bson.M{"$set": bson.M{"title": task.Title, "description": task.Description}}

	_, err := r.db.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r tasks) GetAll(ctx context.Context, userID string) ([]models.Task, error) {
	filter := bson.M{"user_id": userID}
	cur, err := r.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var tasksDoc []models.Task
	for cur.Next(ctx) {
		var task models.Task
		err := cur.Decode(&task)
		if err != nil {
			return nil, err
		}

		tasksDoc = append(tasksDoc, task)
	}

	return tasksDoc, nil
}

func (r tasks) Get(ctx context.Context, userID, id string) (models.Task, error) {
	var task models.Task
	filter := bson.M{"user_id": userID, "id": id}

	err := r.db.FindOne(ctx, filter).Decode(&task)
	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func (r tasks) MarkValueSet(ctx context.Context, taskID string, status int) error {
	filter := bson.M{"id": taskID}
	update := bson.M{"$set": bson.M{"isDone": status}}

	_, err := r.db.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r tasks) Search(ctx context.Context, userID, taskName string) ([]models.Task, error) {
	filter := bson.M{"user_id": userID, "title": taskName}
	cur, err := r.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var tasksDoc []models.Task
	for cur.Next(ctx) {
		var task models.Task
		err := cur.Decode(&task)
		if err != nil {
			return nil, err
		}

		tasksDoc = append(tasksDoc, task)
	}

	return tasksDoc, nil
}
