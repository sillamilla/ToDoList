package task

import (
	models "ToDoWithKolya/internal/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type TaskRepo interface {
	Create(task models.Task) error
	Update(task models.Task, userID string) error

	GetTasksByUserID(userID string) ([]models.Task, error)
	GetByUserID(id string) (models.Task, error)
	GetByID(id string) (models.Task, error)

	MarkValueSet(taskID string, status int) error

	SearchTask(taskName string) ([]models.Task, error)

	DeleteByTaskID(id string, userID string) error
}

type taskRepo struct {
	db *mongo.Collection
}

func NewTaskRepo(database mongo.Database) taskRepo {
	return taskRepo{
		db: database.Collection("tasks"),
	}
}

func (r taskRepo) Create(task models.Task) error {
	_, err := r.db.InsertOne(context.TODO(), task)
	return err
}

func (r taskRepo) GetByUserID(userID string) (models.Task, error) {
	var task models.Task
	err := r.db.FindOne(context.TODO(), bson.M{"user_id": userID}).Decode(&task)
	return task, err
}

func (r taskRepo) DeleteByTaskID(taskID string, userID string) error {
	_, err := r.db.DeleteOne(context.TODO(), bson.M{"_id": taskID, "user_id": userID})
	return err
}

func (r taskRepo) Update(task models.Task, userID string) error {
	filter := bson.M{"_id": task.ID, "user_id": userID}
	update := bson.M{"$set": bson.M{"title": task.Title, "description": task.Description}}
	_, err := r.db.UpdateOne(context.TODO(), filter, update)
	return err
}

func (r taskRepo) GetTasksByUserID(userID string) ([]models.Task, error) {
	filter := bson.M{"user_id": userID}
	cur, err := r.db.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	var tasks []models.Task
	for cur.Next(context.TODO()) {
		var task models.Task
		err := cur.Decode(&task)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r taskRepo) GetByID(id string) (models.Task, error) {
	var task models.Task
	err := r.db.FindOne(context.TODO(), bson.M{"_id": id}).Decode(&task)
	return task, err
}

func (r taskRepo) MarkValueSet(taskID string, status int) error {
	filter := bson.M{"_id": taskID}
	update := bson.M{"$set": bson.M{"isDone": status}}
	_, err := r.db.UpdateOne(context.TODO(), filter, update)
	return err
}

func (r taskRepo) SearchTask(taskName string) ([]models.Task, error) {
	filter := bson.M{"title": taskName}
	cur, err := r.db.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	var tasks []models.Task
	for cur.Next(context.TODO()) {
		var task models.Task
		err := cur.Decode(&task)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}
