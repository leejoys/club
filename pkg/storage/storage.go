package storage

import "context"

//User - хранимый пользователь
type User struct {
	Name  string
	Email string
	Date  int
}

// Interface задаёт контракт на работу с БД.
type Interface interface {
	GetUser(context.Context, string) (User, error)  // получение пользователя по емейлу
	CountUser(context.Context, string) (int, error) // проверка наличия пользователя
	StoreUser(context.Context, User) error          // сохранение нового пользователя
	Close()                                         // освобождение ресурса
}
