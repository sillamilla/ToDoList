package tasks

import (
	models "ToDoWithKolya/internal/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repo interface {
	Create(ctx context.Context, task models.Task) error
	Update(ctx context.Context, task models.Task) error
	GetAll(ctx context.Context, userID string) ([]models.Task, error)
	Get(ctx context.Context, id string) (models.Task, error)
	MarkValueSet(ctx context.Context, taskID string, status int) error
	Search(ctx context.Context, taskName, userID string) ([]models.Task, error)
	Delete(ctx context.Context, id string) error
}

type tasks struct {
	db *mongo.Collection
}

func New(database *mongo.Database) Repo {
	return tasks{
		db: database.Collection("tasks"),
	}
}

func (r tasks) Create(ctx context.Context, task models.Task) error {
	data := bson.M{
		"id":          task.ID,
		"user_id":     task.UserID,
		"title":       task.Title,
		"description": task.Description,
		"is_done":     task.IsDone,
		"created_at":  task.CreatedAt,
	}

	_, err := r.db.InsertOne(ctx, data)
	if err != nil {
		return err
	}

	return nil
}

func (r tasks) Delete(ctx context.Context, id string) error {
	filter := bson.M{"id": id}

	_, err := r.db.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (r tasks) Update(ctx context.Context, task models.Task) error {
	filter := bson.M{"id": task.ID}
	update := bson.M{"$set": bson.M{"title": task.Title, "description": task.Description}}

	_, err := r.db.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r tasks) GetAll(ctx context.Context, userID string) ([]models.Task, error) {
	filter := bson.M{"user_id": userID}
	opts := options.Find().SetSort(bson.D{{"created_at", -1}})

	cur, err := r.db.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer func(cur *mongo.Cursor, ctx context.Context) {
		err = cur.Close(ctx)
		if err != nil {
			return
		}
	}(cur, ctx)

	var tasksDoc []models.Task
	for cur.Next(ctx) {
		var task models.Task
		err = cur.Decode(&task)
		if err != nil {
			return nil, err
		}

		tasksDoc = append(tasksDoc, task)
	}

	return tasksDoc, nil
}

func (r tasks) Get(ctx context.Context, id string) (models.Task, error) {
	var task models.Task
	filter := bson.M{"id": id}

	err := r.db.FindOne(ctx, filter).Decode(&task)
	if err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func (r tasks) MarkValueSet(ctx context.Context, taskID string, status int) error {
	filter := bson.M{"id": taskID}
	update := bson.M{"$set": bson.M{"is_done": status}}

	_, err := r.db.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return nil
}

func (r tasks) Search(ctx context.Context, taskName, userID string) ([]models.Task, error) {
	// Создаем регулярное выражение для поиска частичных совпадений
	regexPattern := primitive.Regex{Pattern: taskName, Options: "i"} // "i" - регистронезависимый поиск

	// Создаем фильтр с использованием $regex оператора
	filter := bson.M{"title": bson.M{"$regex": regexPattern}, "user_id": userID}

	cur, err := r.db.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	var tasksDoc []models.Task
	for cur.Next(ctx) {
		var task models.Task
		err = cur.Decode(&task)
		if err != nil {
			return nil, err
		}

		tasksDoc = append(tasksDoc, task)
	}

	return tasksDoc, nil
}
