package mongodb

import (
	"club/pkg/storage"
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ErrorDuplicateUser error = errors.New("MongoDB E11000")

// Хранилище данных.
type Store struct {
	c  *mongo.Client
	db *mongo.Database
}

//New - Конструктор объекта хранилища.
func New(name string, connstr string) (*Store, error) {
	client, err := mongo.Connect(context.Background(),
		options.Client().ApplyURI(connstr))
	if err != nil {
		return nil, err
	}
	// проверка связи с БД
	err = client.Ping(context.Background(), nil)
	if err != nil {
		client.Disconnect(context.Background())
		return nil, err
	}

	s := &Store{c: client, db: client.Database(name)}
	t := true
	_, err = s.db.Collection("Users").Indexes().CreateOne(
		context.Background(), mongo.IndexModel{
			Keys:    bson.D{{Key: "title", Value: 1}},
			Options: &options.IndexOptions{Unique: &t}})
	if err != nil {
		s.c.Disconnect(context.Background())
		return nil, err
	}

	return s, nil
}

//Close - освобождение ресурса
func (s *Store) Close() {
	s.c.Disconnect(context.Background())
}

func (s *Store) DropDB() error {
	return s.db.Drop(context.Background())
}

//Users - получение всех пользователей
func (s *Store) Users() ([]storage.User, error) {

	coll := s.db.Collection("Users")
	ctx := context.Background()
	filter := bson.D{}
	cur, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	Users := []storage.User{}
	for cur.Next(ctx) {
		var u storage.User
		err = cur.Decode(&u)
		if err != nil {
			return nil, err
		}
		Users = append(Users, u)
	}
	return Users, nil
}

//UsersN - получение n последних пользователей
func (s *Store) UsersN(n int) ([]storage.User, error) {

	coll := s.db.Collection("Users")
	ctx := context.Background()
	options := options.Find()
	options.SetLimit(int64(n))
	options.SetSort(bson.D{{Key: "$natural", Value: -1}})
	filter := bson.D{}
	cur, err := coll.Find(ctx, filter, options)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)

	Users := []storage.User{}
	for cur.Next(ctx) {
		var u storage.User
		err = cur.Decode(&u)
		if err != nil {
			return nil, err
		}
		Users = append(Users, u)
	}
	return Users, nil
}

//AddUser - создание нового пользователя
func (s *Store) AddUser(u storage.User) error {
	coll := s.db.Collection("Users")
	_, err := coll.InsertOne(context.Background(), u)

	if mongo.IsDuplicateKeyError(err) {
		return ErrorDuplicateUser
	}
	return err
}

//UpdateUser - обновление по email значения name
func (s *Store) UpdateUser(u storage.User) error {
	coll := s.db.Collection("Users")
	filter := bson.D{{Key: "email", Value: u.Email}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "name", Value: u.Name}}}}
	_, err := coll.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return nil
}

//DeleteUser - удаляет пользователя по email
func (s *Store) DeleteUser(u storage.User) error {
	coll := s.db.Collection("Users")
	filter := bson.D{{Key: "email", Value: u.Email}}
	_, err := coll.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return nil
}
