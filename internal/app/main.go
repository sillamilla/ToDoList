package main

import (
	"ToDoWithKolya/internal/config"
	"ToDoWithKolya/internal/handler/api"
	ui2 "ToDoWithKolya/internal/handler/ui"
	"ToDoWithKolya/internal/repository"
	"ToDoWithKolya/internal/service"
	"context"
	"github.com/go-chi/chi/v5"
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

	defer func() {
		if err = dbMongo.Disconnect(context.Background()); err != nil {
			log.Println(err)
		}
	}()

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

	repo := repository.New(dbMongo.Database(cfg.Mongo.DBName, nil))
	srv := service.New(repo)

	// API Routes
	{
		api := api.New(srv)
		auth := api.Auth.Authorization

		r := chi.NewRouter()
		r.Group(func(r chi.Router) {
			http.Handle("/api/", http.StripPrefix("/api", r))
			r.Post("/sign-up", api.Auth.SignUp)
			r.Post("/sign-in", api.Auth.SignIn)
			r.With(auth).Delete("/logout", api.Auth.Logout)
			r.With(auth).Post("/create", api.Task.Create)
			r.With(auth).Put("/edit", api.Task.Edit)
			r.With(auth).Get("/task/{id}", api.Task.TaskByID)
			r.With(auth).Delete("/delete/{id}", api.Task.Delete)
			r.With(auth).Get("/tasks", api.Task.GetTasks)
		})
	}

	// UI Routes
	{
		ui := ui2.New(srv)
		auth := ui.Auth.Authorization

		r := chi.NewRouter()
		http.Handle("/", r)
		r.Group(func(r chi.Router) {
			r.Get("/sign-up", ui.Auth.SignUp)
			r.Post("/sign-up", ui.Auth.SignUpPost)
			r.Get("/sign-in", ui.Auth.SignIn)
			r.Post("/sign-in", ui.Auth.SignInPost)
			r.With(auth).Delete("/logout", ui.Auth.Logout)
			r.With(auth).Get("/", ui.HomePage)
			r.With(auth).Get("/create", ui.Task.Create)
			r.With(auth).Post("/create", ui.Task.CreatePost)
			r.With(auth).Get("/edit/{id}", ui.Task.Edit)
			r.With(auth).Put("/edit/{id}", ui.Task.EditPost)
			r.With(auth).Delete("/delete/{id}", ui.Task.Delete)
			r.With(auth).Put("/mark/{id}/{status}", ui.Task.MarkAsDone)
			r.With(auth).Get("/search/{search}", ui.Task.Search)
		})
	}

	if err = http.ListenAndServe(":"+cfg.HTTP.Port, nil); err != nil {
		log.Fatal(err)
	}

}
