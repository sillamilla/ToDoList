package ui

import (
	"ToDoWithKolya/internal/ctxpkg"
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
	TaskHandler tasks.Handler
	UserHandler users.Handler
	Home        *template.Template
	srv         *service.Service
}

func New(srv *service.Service) Handler {
	home, err := template.ParseFiles("./internal/templates/home.html")
	if err != nil {
		panic(err)
	}

	return Handler{
		TaskHandler: tasks.NewHandler(srv.TaskSrv),
		UserHandler: users.NewHandler(srv.UserSrv),
		Home:        home,
		srv:         srv,
	}
}

func (h Handler) HomePage(w http.ResponseWriter, r *http.Request) {
	user, ok := ctxpkg.UserFromContext(r.Context())
	if !ok {
		http.Redirect(w, r, "/sign-in", http.StatusSeeOther)
		return
	}

	tasks, err := h.srv.TaskSrv.GetTasksByUserID(user.ID)
	if err != nil {
		errs.HandleError(w, err, http.StatusInternalServerError)
		return
	}

	sort.Slice(tasks, func(i, j int) bool {
		if tasks[i].CreatedDate.After(tasks[j].CreatedDate) {
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
