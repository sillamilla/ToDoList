package ui

import (
	"ToDoWithKolya/internal/ctxpkg"
	"ToDoWithKolya/internal/handler/ui/errs"
	"ToDoWithKolya/internal/handler/ui/tasks"
	"ToDoWithKolya/internal/handler/ui/users"
	"ToDoWithKolya/internal/service"
	"html/template"
	"net/http"
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
		http.Redirect(w, r, "/sign-in", http.StatusPermanentRedirect)
		return
	}

	tasks, err := h.srv.TaskSrv.GetTasksByUserID(user.ID)
	if err != nil {
		errs.ErrorWrap(w, err, http.StatusInternalServerError)
		return
	}

	err = h.Home.Execute(w, tasks)
	if err != nil {
		errs.ErrorWrap(w, err, http.StatusInternalServerError)
		return
	}
}
