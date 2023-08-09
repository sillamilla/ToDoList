package main

import (
	"ToDoWithKolya/internal/config"
	"ToDoWithKolya/internal/repository"
	"ToDoWithKolya/internal/service"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func main() {
	cfg := config.GetConfig()

	dbMongo, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.Mongo.Address))
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err = dbMongo.Disconnect(context.Background()); err != nil {
			log.Println(err)
		}
	}()

	//defer func(dbMongo *mongo.Client, ctx context.Context) {
	//	err = dbMongo.Disconnect(ctx)
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}(dbMongo, context.Background())

	err = dbMongo.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Connect error Mongo:", err)
	}

	repo := repository.New(dbMongo.Database(cfg.Mongo.DBName, nil))
	srv := service.New(repo)
	fmt.Print(srv)
	//r := mux.NewRouter()
	//
	//{
	//	apiHnd := api.New(srv)
	//	api := r.PathPrefix("/api").Subrouter()
	//	auth := apiHnd.UserHandler.Authorization
	//
	//	api.HandleFunc("/user", apiHnd.UserHandler.Register).Methods(http.MethodPost)
	//	api.HandleFunc("/user/login", apiHnd.UserHandler.Login).Methods(http.MethodPost)
	//	api.HandleFunc("/user", auth(apiHnd.UserHandler.Logout)).Methods(http.MethodDelete)
	//
	//	api.HandleFunc("/tasks", auth(apiHnd.TaskHandler.Create)).Methods(http.MethodPost)
	//	api.HandleFunc("/tasks/edit/{id}", auth(apiHnd.TaskHandler.Edit)).Methods(http.MethodPost)
	//	api.HandleFunc("/tasks", auth(apiHnd.TaskHandler.GetTasksByUserID)).Methods(http.MethodGet)
	//	api.HandleFunc("/tasks/{id}", auth(apiHnd.TaskHandler.GetTaskByID)).Methods(http.MethodGet)
	//	api.HandleFunc("/tasks/{id}", auth(apiHnd.TaskHandler.DeleteByTaskID)).Methods(http.MethodDelete)
	//
	//}
	//
	//{
	//	uiHnd := ui.New(srv)
	//	auth := uiHnd.User.Authorization
	//
	//	{
	//		r.HandleFunc("/sign-up", uiHnd.User.SignUpPost).Methods(http.MethodPost)
	//		r.HandleFunc("/sign-up", uiHnd.User.SignUp).Methods(http.MethodGet)
	//		r.HandleFunc("/sign-in", uiHnd.User.SignInPost).Methods(http.MethodPost)
	//		r.HandleFunc("/sign-in", uiHnd.User.SignIn).Methods(http.MethodGet)
	//		r.HandleFunc("/logout", auth(uiHnd.User.Logout)).Methods(http.MethodGet)
	//	}
	//
	//	{
	//		r.HandleFunc("/", auth(uiHnd.HomePage)).Methods(http.MethodGet)
	//		r.HandleFunc("/create", auth(uiHnd.Task.CreatePost)).Methods(http.MethodPost)
	//		r.HandleFunc("/create", auth(uiHnd.Task.Create)).Methods(http.MethodGet)
	//		r.HandleFunc("/edit/{id}", auth(uiHnd.Task.EditPost)).Methods(http.MethodPost)
	//		r.HandleFunc("/edit/{id}", auth(uiHnd.Task.Edit)).Methods(http.MethodGet)
	//		r.HandleFunc("/mark/{taskID}/{status}", auth(uiHnd.Task.MarkAsDone)).Methods(http.MethodGet)
	//		r.HandleFunc("/search/{search}", auth(uiHnd.Task.Search)).Methods(http.MethodGet)
	//		r.HandleFunc("/tasks/delete/{id}", auth(uiHnd.Task.Delete)).Methods(http.MethodGet)
	//		r.HandleFunc("/tasks/deleteAll/{id}", auth(uiHnd.Task.Delete)).Methods(http.MethodGet)
	//	}
	//}

	//if err = http.ListenAndServe(":8080", r); err != nil {
	//	log.Fatal(err)
	//}
}
