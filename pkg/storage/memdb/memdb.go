package memdb

import (
	"club/pkg/storage"
	"context"
	"errors"
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

//GetUser - получение пользователя по емейлу
func (s *Store) GetUser(ctx context.Context, email string) (storage.User, error) {
	s.db.mutex.RLock()
	user, ok := s.db.dbase[email]
	s.db.mutex.RUnlock()
	if !ok {
		return storage.User{}, errors.New("memdb GetUser error: no data")
	}
	return user, nil
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
