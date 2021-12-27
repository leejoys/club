package main

import (
	"context"
	"html/template"
	"regexp"

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
		"templates/success.html", "templates/header.html", "templates/footer.html",
		"templates/about.html", "templates/contacts.html")
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
	srv.api = api.New(ctx, srv.db, tmpl, regexp.MustCompile(`^[a-zA-Z\s\.]+$`))

	// Запускаем веб-сервер, передаём серверу маршрутизатор запросов.
	go func() {
		port := os.Getenv("PORT")
		log.Fatal(http.ListenAndServe("0.0.0.0:"+port, srv.api.Router()))
	}()
	log.Println("HTTP server is started")
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	<-signalCh
	log.Println("HTTP server has been stopped")
}
