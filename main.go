package main

import (
	"ToDoWithKolya/internal/handler/api"
	"ToDoWithKolya/internal/handler/ui"
	"ToDoWithKolya/internal/repository"
	"ToDoWithKolya/internal/service"
	"ToDoWithKolya/pkg/sql_lite"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	database, err := sql_lite.New()
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	repo := repository.New(database)
	srv := service.New(repo)

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

		r.HandleFunc("/", auth(uiHnd.HomePage)).Methods(http.MethodPost)
		r.HandleFunc("/task/{id}", auth(uiHnd.TaskHandler.Task)).Methods(http.MethodGet)
		//r.HandleFunc("/task/edit/{id}", uiHnd.TaskHandler.Edit).Methods(http.MethodPost)
	}

	http.ListenAndServe(":8080", r)
}
