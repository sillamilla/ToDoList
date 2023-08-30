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
	"sort"
)

type Handler struct {
	Task tasks.Handler
	User users.Handler
	Auth auth.Handler

	Home *template.Template
	srv  *service.Service
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

		Home: home,
		srv:  srv,
	}
}

func (h Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(models.User)
	if !ok {
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return
	}

	tasks, err := h.srv.Task.GetTasks(r.Context(), user.ID)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	sort.Slice(tasks, func(i, j int) bool {
		if tasks[i].CreatedAt.After(tasks[j].CreatedAt) {
			return true
		}
		return false
	})

	userAndTask := models.UserAndTask{
		User:  user,
		Tasks: tasks,
	}

	err = h.Home.Execute(w, userAndTask)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}
}
