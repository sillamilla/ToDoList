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
	r.HandleFunc("/user/login", hnd.UserHandler.Login).Methods(http.MethodPost)
	r.HandleFunc("/user", auth(hnd.UserHandler.Logout)).Methods(http.MethodDelete)

	//todo fix delete таски інших людей(готовов) / Перевіряти чи це таска цього юзера
	//todo ? Я зміминв принцип роботи лог аут, там був json і потім доставали юзера і по юзер id удаляли сесію
	//todo додати логи

	//task
	r.HandleFunc("/task", auth(hnd.TaskHandler.Create)).Methods(http.MethodPost)
	r.HandleFunc("/tasks", auth(hnd.TaskHandler.GetTasksByUserID)).Methods(http.MethodGet)
	r.HandleFunc("/task/{id}", auth(hnd.TaskHandler.GetTaskByID)).Methods(http.MethodGet)
	r.HandleFunc("/task/{id}", auth(hnd.TaskHandler.DeleteByTaskID)).Methods(http.MethodDelete)
	http.ListenAndServe(":8080", r)
}
