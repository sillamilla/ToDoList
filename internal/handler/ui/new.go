package ui

import (
	"ToDoWithKolya/internal/handler/ui/auth"
	"ToDoWithKolya/internal/handler/ui/errs"
	"ToDoWithKolya/internal/handler/ui/tasks"
	"ToDoWithKolya/internal/handler/ui/users"
	"ToDoWithKolya/internal/models"
	"ToDoWithKolya/internal/service"
	"html/template"
	"net/http"
)

type Handler struct {
	Task tasks.Handler
	User users.Handler
	Auth auth.Handler

	HomeTemplate *template.Template
	srv          *service.Service
}

func New(srv *service.Service) Handler {
	home, err := template.ParseFiles("./internal/templates/home.html")
	if err != nil {
		panic(err)
	}

	return Handler{
		Task: tasks.New(srv.Task),
		User: users.New(srv.User),
		Auth: auth.New(srv.Auth),

		HomeTemplate: home,
		srv:          srv,
	}
}

func (h Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value("id").(string)
	if !ok {
		errs.HandleError(w, nil, http.StatusInternalServerError) //todo models.NoCtxContent
		return
	}

	tasks, err := h.srv.Task.GetTasks(r.Context(), id)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}
	user, err := h.srv.User.GetByID(r.Context(), id)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	userAndTask := models.UserAndTask{
		User:  user,
		Tasks: tasks,
	}

	err = h.HomeTemplate.Execute(w, userAndTask)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}
}
