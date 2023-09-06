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
	api := api.New(srv)
	apiAuth := api.Auth.Authorization
	router := chi.NewRouter()

	router.Route("/api", func(r chi.Router) {
		r.Post("/sign-up", api.Auth.SignUp)
		r.Post("/sign-in", api.Auth.SignIn)

		authRouter := r.With(apiAuth)
		authRouter.Delete("/logout", api.Auth.Logout)
		authRouter.Post("/create", api.Task.Create)
		authRouter.Put("/edit", api.Task.Edit)
		authRouter.Get("/task/{id}", api.Task.TaskByID)
		authRouter.Delete("/delete/{id}", api.Task.Delete)
		authRouter.Get("/tasks", api.Task.GetTasks)
	})

	// UI Routes
	ui := ui2.New(srv)
	auth := ui.Auth.Authorization

	router.Route("/", func(r chi.Router) {
		r.Get("/sign-up", ui.Auth.SignUp)
		r.Post("/sign-up", ui.Auth.SignUpPost)
		r.Get("/sign-in", ui.Auth.SignIn)
		r.Post("/sign-in", ui.Auth.SignInPost)

		authRouter := r.With(auth)
		authRouter.Delete("/logout", ui.Auth.Logout)
		authRouter.Get("/", ui.HomePage)
		authRouter.Get("/create", ui.Task.Create)
		authRouter.Post("/create", ui.Task.CreatePost)
		authRouter.Get("/edit/{id}", ui.Task.Edit)
		authRouter.Put("/edit/{id}", ui.Task.EditPost)
		authRouter.Delete("/delete/{id}", ui.Task.Delete)
		authRouter.Put("/mark/{id}/{status}", ui.Task.MarkAsDone)
		authRouter.Get("/search/{search}", ui.Task.Search)
	})

	if err = http.ListenAndServe(":"+cfg.HTTP.Port, router); err != nil {
		log.Fatal(err)
	}

}
