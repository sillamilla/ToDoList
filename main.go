package main

import (
	"ToDoWithKolya/internal/handler"
	"ToDoWithKolya/internal/repository"
	"ToDoWithKolya/internal/service"
	"ToDoWithKolya/pkg/db"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	database, err := db.New()
	if err != nil {
		log.Fatal(err)
	}

	r := mux.NewRouter()

	repo := repository.New(database)
	srv := service.New(repo)
	hnd := handler.New(srv)

	auth := hnd.UserHandler.Authorization
	//user
	r.HandleFunc("/user", hnd.UserHandler.Register).Methods(http.MethodPost)
	r.HandleFunc("/user/edit/{id}", auth(hnd.UserHandler.Edit)).Methods(http.MethodPost)
	r.HandleFunc("/user/login", hnd.UserHandler.Login).Methods(http.MethodPost)
	r.HandleFunc("/user", auth(hnd.UserHandler.Logout)).Methods(http.MethodDelete)

	//todo Перевіряти чи сесія цього юзера ідентична сесії .зера якого ми хочемо модефікувати, інаккше кік
	//task
	r.HandleFunc("/task", auth(hnd.TaskHandler.Create)).Methods(http.MethodPost)
	r.HandleFunc("/task/edit/{id}", auth(hnd.TaskHandler.Edit)).Methods(http.MethodPost)
	r.HandleFunc("/tasks", auth(hnd.TaskHandler.GetTasksByUserID)).Methods(http.MethodGet)
	r.HandleFunc("/task/{id}", auth(hnd.TaskHandler.GetTaskByID)).Methods(http.MethodGet)
	r.HandleFunc("/task/{id}", auth(hnd.TaskHandler.DeleteByTaskID)).Methods(http.MethodDelete)
	http.ListenAndServe(":8080", r)
}
