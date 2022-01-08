package storage

import "context"

//User - хранимый пользователь
type User struct {
	Name  string
	Email string
	Date  string
}

// Interface задаёт контракт на работу с БД.
type Interface interface {
	Users(context.Context) ([]User, error)          // получение пользователей
	CountUser(context.Context, string) (int, error) // проверка наличия пользователя
	StoreUser(context.Context, User) error          // сохранение нового пользователя
	Close()                                         // освобождение ресурса
	DropDB() error                                  // очистка базы
}
