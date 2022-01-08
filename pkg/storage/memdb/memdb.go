package memdb

import (
	"club/pkg/storage"
	"context"
	"sync"
)

type inmemory struct {
	mutex sync.RWMutex
	dbase map[string]storage.User
}

// Хранилище данных.
type Store struct {
	db *inmemory
}

//todo context.WithTimeout (ctx, time.Second*30)
//New - Конструктор объекта хранилища.
func New() *Store {
	return &Store{db: &inmemory{sync.RWMutex{},
		make(map[string]storage.User)}}
}

//Close - освобождение ресурса. Заглушка для реализации интерфейса.
func (s *Store) Close() {}

//DropDB - очистка базы.
func (s *Store) DropDB() error {
	s.db.dbase = make(map[string]storage.User)
	return nil
}

//Users - получение пользователей
func (s *Store) Users(ctx context.Context) ([]storage.User, error) {
	usersList := []storage.User{}
	s.db.mutex.RLock()
	for _, el := range s.db.dbase {
		usersList = append(usersList, el)
	}
	s.db.mutex.RUnlock()
	return usersList, nil
}

//StoreUser - сохранение нового пользователя
func (s *Store) StoreUser(ctx context.Context, u storage.User) error {
	s.db.mutex.Lock()
	s.db.dbase[u.Email] = u
	s.db.mutex.Unlock()
	return nil
}

//CountUser - проверка наличия пользователя
func (s *Store) CountUser(ctx context.Context, email string) (int, error) {
	s.db.mutex.RLock()
	defer s.db.mutex.RUnlock()
	if _, ok := s.db.dbase[email]; !ok {
		return 0, nil
	}
	return 1, nil
}
