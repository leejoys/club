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

//API - Программный интерфейс сервиса
type API struct {
	t   *template.Template
	db  storage.Interface
	r   *mux.Router
	ctx context.Context
}

//New - Конструктор объекта API
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
	//метод добавления пользователя
	api.r.HandleFunc("/", api.storeUser).Methods(http.MethodPost)
}

//Router - Получение маршрутизатора запросов.
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
		err = api.t.ExecuteTemplate(w, "error", err.Error())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	err = api.t.ExecuteTemplate(w, "index", re)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// // метод добавления пользователя
func (api *API) storeUser(w http.ResponseWriter, r *http.Request) {
	//получаем имя из тела запроса
	userName := r.FormValue("name")
	//todo regex name check
	if userName == "" {
		err := api.t.Execute(w, "wrong name")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	//получаем адрес из тела запроса
	userEmail := r.FormValue("email")
	if !isEmailValid(userEmail) {
		err := api.t.ExecuteTemplate(w, "error", "wrong email")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	//считаем количество пользователей с таким адресом в базе
	i, err := api.db.CountUser(api.ctx, userEmail)
	if err != nil {
		err = api.t.ExecuteTemplate(w, "error", err.Error())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	//если пользователей больше нуля - ошибка
	if i > 0 {
		err = api.t.ExecuteTemplate(w, "error", "email already in use")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	//сохраняем пользователя
	err = api.db.StoreUser(api.ctx,
		storage.User{Name: userName,
			Email: userEmail,
			Date:  time.Now().UTC().Format("2006-02-01")})
	if err != nil {
		err = api.t.ExecuteTemplate(w, "error", err.Error())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	//формируем ответ для отправки в шаблон
	re := resp{}
	re.List, err = api.db.Users(api.ctx)
	if err != nil {
		err = api.t.ExecuteTemplate(w, "error", err.Error())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		return
	}
	//рендерим шаблон с ответом
	err = api.t.ExecuteTemplate(w, "index", re)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
