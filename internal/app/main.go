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
		auth := api.Session.Authorization

		r := chi.NewRouter()
		r.Handle("/api", http.StripPrefix("/api", r))

		r.Post("/sign-up", api.User.Register)
		r.Post("/sign-in", api.User.Login)
		r.Delete("/logout", auth(api.Session.Logout))

		r.Post("/new-task", auth(api.Task.Create))
		r.Put("/task-edit", auth(api.Task.Edit))
		r.Get("/task-get", auth(api.Task.TaskByID))
		r.Delete("/task-delete", auth(api.Task.Delete))
		r.Get("/task-get-all", auth(api.Task.GetTasks))

		http.Handle("/api/", http.StripPrefix("/api", r))
	}

	// UI Routes
	{
		ui := ui2.New(srv)
		auth := ui.Session.Authorization

		r := chi.NewRouter()
		r.Handle("/", http.StripPrefix("/", r))

		r.Get("/sign-up", ui.User.SignUp)
		r.Post("/sign-up", ui.User.SignUpPost)
		r.Get("/sign-in", ui.User.SignIn)
		r.Post("/sign-in", ui.User.SignInPost)
		r.Delete("/logout", auth(ui.Session.Logout))

		r.Get("/", auth(ui.HomePage))

		r.Get("/create", auth(ui.Task.Create))
		r.Post("/create", auth(ui.Task.CreatePost))
		r.Get("/edit/{id}", auth(ui.Task.Edit))
		r.Put("/edit/{id}", auth(ui.Task.EditPost))
		r.Delete("/delete/{id}", auth(ui.Task.Delete))
		r.Put("/mark/{id}/{status}", auth(ui.Task.MarkAsDone))
		r.Get("/search/{search}", auth(ui.Task.Search))
	}

	if err = http.ListenAndServe(":"+cfg.HTTP.Port, nil); err != nil { //todo env port
		log.Fatal(err)
	}

}
