package main

import (
	"context"

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

	pageFile, err := os.ReadFile("page.html")
	if err != nil {
		log.Fatal("page file reading error")
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
	srv.api = api.New(ctx, srv.db, pageFile)

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
