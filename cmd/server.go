package main

import (
	"context"
	"fmt"
	"html/template"
	"regexp"

	"club/pkg/api"
	"club/pkg/storage"
	"club/pkg/storage/memdb"
	"club/pkg/storage/mongodb"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Создаём объект базы данных MongoDB.
	pwd := os.Getenv("Cloud0pass")
	connstr := fmt.Sprintf(
		"mongodb+srv://sup:%s@cloud0.wspoq.mongodb.net/clubusers?retryWrites=true&w=majority",
		pwd)
	db, err := mongodb.New("clubusers", connstr)
	if err != nil {
		log.Fatalf("mongo.New error: %s", err)
	}

	// Создаём объект базы данных в памяти
	_ = memdb.New()

	// Создаём объект сервера
	srv := server{}

	// Инициализируем БД
	srv.db = db

	// Освобождаем ресурс
	defer srv.db.Close()

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(ctx, srv.db, tmpl, regexp.MustCompile(`^[a-zA-Z\s\.]+$`))

	// Запускаем веб-сервер, передаём серверу маршрутизатор запросов.
	port := os.Getenv("PORT")
	go func() {
		log.Fatal(http.ListenAndServe("0.0.0.0:"+port, srv.api.Router()))
	}()
	log.Printf("HTTP server is started on port %s\n", port)
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)
	<-signalCh
	log.Println("HTTP server has been stopped")
}
