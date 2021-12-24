package api

import (
	"club/pkg/storage"
	"context"
	"html/template"
	"net"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

//проверка адреса по написанию и домена
func isEmailValid(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	valid := func() bool {
		_, err := mail.ParseAddress(e)
		return err == nil
	}
	if !valid() {
		return false
	}
	parts := strings.Split(e, "@")
	mx, err := net.LookupMX(parts[1])
	if err != nil || len(mx) == 0 {
		return false
	}
	return true
}

// Программный интерфейс сервиса
type API struct {
	t   *template.Template
	db  storage.Interface
	r   *mux.Router
	ctx context.Context
}

// Конструктор объекта API
func New(ctx context.Context, db storage.Interface, t *template.Template) *API {
	api := API{
		db:  db,
		t:   t,
		ctx: ctx,
	}
	api.r = mux.NewRouter()
	api.endpoints()
	return &api
}

// Регистрация обработчиков API.
func (api *API) endpoints() {
	//метод отображения страницы
	api.r.HandleFunc("/", api.page).Methods(http.MethodGet)
	// //метод получения данных для таблицы
	// api.r.HandleFunc("/users", api.getData).Methods(http.MethodGet)
	//метод добавления пользователя
	api.r.HandleFunc("/", api.storeUser).Methods(http.MethodPost)
}

// Получение маршрутизатора запросов.
// Требуется для передачи маршрутизатора веб-серверу.
func (api *API) Router() *mux.Router {
	return api.r
}

type resp struct {
	List []storage.User
}

//отображение страницы
func (api *API) page(w http.ResponseWriter, r *http.Request) {
	var err error
	re := resp{}
	re.List, err = api.db.Users(api.ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = api.t.Execute(w, re)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// // метод получения данных для таблицы
// func (api *API) getData(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte())
// }

// // метод добавления пользователя
func (api *API) storeUser(w http.ResponseWriter, r *http.Request) {
	userName := r.FormValue("name")
	//todo regex name check
	if userName == "" {
		http.Error(w, "wrong name", http.StatusBadRequest)
		return
	}
	userEmail := r.FormValue("email")

	if !isEmailValid(userEmail) {
		http.Error(w, "wrong email", http.StatusBadRequest)
		return
	}
	i, err := api.db.CountUser(api.ctx, userEmail)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if i > 0 {
		http.Error(w, "email already in use", http.StatusBadRequest)
		return
	}
	err = api.db.StoreUser(api.ctx,
		storage.User{Name: userEmail,
			Email: userEmail,
			Date:  time.Now().UTC().Format(time.UnixDate)})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	re := resp{}
	re.List, err = api.db.Users(api.ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = api.t.Execute(w, re)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
