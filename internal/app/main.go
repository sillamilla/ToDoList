package main

import (
	"ToDoWithKolya/internal/config"
	"ToDoWithKolya/internal/handler/api"
	"ToDoWithKolya/internal/handler/ui"
	"ToDoWithKolya/internal/repository"
	"ToDoWithKolya/internal/service"
	"context"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

func main() {
	cfg := config.GetConfig()

	dbMongo, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.Mongo.Address))
	if err != nil {
		log.Fatal(err)
	}
	defer func(dbMongo *mongo.Client, ctx context.Context) {
		err = dbMongo.Disconnect(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(dbMongo, context.Background())

	err = dbMongo.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Connect error Mongo:", err)
	}

	repo := repository.New(*dbMongo.Database(cfg.Mongo.DBName, nil))
	srv := service.New(repo)

	r := mux.NewRouter()

	{
		apiHnd := api.New(srv)
		api := r.PathPrefix("/api").Subrouter()
		auth := apiHnd.UserHandler.Authorization

		api.HandleFunc("/user", apiHnd.UserHandler.Register).Methods(http.MethodPost)
		api.HandleFunc("/user/login", apiHnd.UserHandler.Login).Methods(http.MethodPost)
		api.HandleFunc("/user", auth(apiHnd.UserHandler.Logout)).Methods(http.MethodDelete)

		api.HandleFunc("/task", auth(apiHnd.TaskHandler.Create)).Methods(http.MethodPost)
		api.HandleFunc("/task/edit/{id}", auth(apiHnd.TaskHandler.Edit)).Methods(http.MethodPost)
		api.HandleFunc("/tasks", auth(apiHnd.TaskHandler.GetTasksByUserID)).Methods(http.MethodGet)
		api.HandleFunc("/task/{id}", auth(apiHnd.TaskHandler.GetTaskByID)).Methods(http.MethodGet)
		api.HandleFunc("/task/{id}", auth(apiHnd.TaskHandler.DeleteByTaskID)).Methods(http.MethodDelete)

	}

	{
		uiHnd := ui.New(srv)
		auth := uiHnd.UserHandler.Authorization

		r.HandleFunc("/sign-up", uiHnd.UserHandler.SignUpPost).Methods(http.MethodPost)
		r.HandleFunc("/sign-up", uiHnd.UserHandler.SignUp).Methods(http.MethodGet)
		r.HandleFunc("/sign-in", uiHnd.UserHandler.SignInPost).Methods(http.MethodPost)
		r.HandleFunc("/sign-in", uiHnd.UserHandler.SignIn).Methods(http.MethodGet)

		r.HandleFunc("/logout", auth(uiHnd.UserHandler.Logout)).Methods(http.MethodGet)

		r.HandleFunc("/", auth(uiHnd.HomePage)).Methods(http.MethodGet)

		r.HandleFunc("/create", auth(uiHnd.TaskHandler.CreatePost)).Methods(http.MethodPost)
		r.HandleFunc("/create", auth(uiHnd.TaskHandler.Create)).Methods(http.MethodGet)

		r.HandleFunc("/edit/{id}", auth(uiHnd.TaskHandler.EditPost)).Methods(http.MethodPost)
		r.HandleFunc("/edit/{id}", auth(uiHnd.TaskHandler.Edit)).Methods(http.MethodGet)

		r.HandleFunc("/mark/{taskID}/{status}", auth(uiHnd.TaskHandler.MarkAsDone)).Methods(http.MethodGet)

		r.HandleFunc("/search/{search}", auth(uiHnd.TaskHandler.Search)).Methods(http.MethodGet)

		r.HandleFunc("/task/delete/{id}", auth(uiHnd.TaskHandler.Delete)).Methods(http.MethodGet)
		r.HandleFunc("/task/deleteAll/{id}", auth(uiHnd.TaskHandler.Delete)).Methods(http.MethodGet)

	}

	if err = http.ListenAndServe(":8080", r); err != nil {
		log.Fatal(err)
	}
}
