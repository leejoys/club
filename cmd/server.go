package main

import (
	"context"
	"html/template"

	"club/pkg/api"
	"club/pkg/storage"
	"club/pkg/storage/memdb"
	"log"
	"net/http"
	"os"
	"os/signal"
)

// Сервер клуба.
type server struct {
	db  storage.Interface
	api *api.API
}

func main() {
	// загружаем шаблоны страницы
	tmpl, err := template.ParseFiles("templates/index.html", "templates/error.html",
		"templates/header.html", "templates/footer.html")
	if err != nil {
		log.Fatal("template files reading error")
	}

	//todo gracefull shutdown
	// Создаём объект сервера
	srv := server{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Инициализируем БД
	srv.db = memdb.New()

	// Освобождаем ресурс
	defer srv.db.Close()

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(ctx, srv.db, tmpl)

	// Запускаем веб-сервер на порту 8080 на всех интерфейсах.
	// Предаём серверу маршрутизатор запросов.
	go func() {
		log.Fatal(http.ListenAndServe("0.0.0.0:8080", srv.api.Router()))
	}()
	log.Println("HTTP server is started")
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	<-signalCh
	log.Println("HTTP server has been stopped")
}
