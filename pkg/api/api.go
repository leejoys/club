package api

import (
	"club/pkg/storage"
	"context"
	"net"
	"net/http"
	"net/mail"
	"strings"

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
	pageFile []byte
	db       storage.Interface
	r        *mux.Router
	ctx      context.Context
}

// Конструктор объекта API
func New(ctx context.Context, db storage.Interface, pageFile []byte) *API {
	api := API{
		db:       db,
		pageFile: pageFile,
		ctx:      ctx,
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
	api.r.HandleFunc("/user", api.storeUser).Methods(http.MethodGet)
}

// Получение маршрутизатора запросов.
// Требуется для передачи маршрутизатора веб-серверу.
func (api *API) Router() *mux.Router {
	return api.r
}

//отображение страницы
func (api *API) page(w http.ResponseWriter, r *http.Request) {
	w.Write(api.pageFile)
}

// // метод получения данных для таблицы
// func (api *API) getData(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte())
// }

// // метод добавления пользователя
func (api *API) storeUser(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query().Get("name")
	if userName == "" {
		http.Error(w, "wrong name", http.StatusBadRequest)
		return
	}
	userEmail := r.URL.Query().Get("email")

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
			Email: userEmail})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(api.pageFile)
}