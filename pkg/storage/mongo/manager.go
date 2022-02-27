package mongo

import (
	"context"
	"fmt"
	"fpl-live-tracker/pkg/config"
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/storage"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

//
type managerRepository struct {
	db       *mongo.Database
	managers *mongo.Collection
}

//
func NewManagerRepository(config config.MongoConfig) (domain.ManagerRepository, error) {
	client, err := mongo.Connect(context.Background(),
		options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:27017", config.Host)))
	if err != nil {
		return nil, err
	}
	// defer client.Disconnect(context.TODO())

	db := client.Database(config.Database)
	managers := db.Collection("managers")

	return &managerRepository{
		db:       db,
		managers: managers,
	}, nil
}

//
func (mr *managerRepository) Add(manager domain.Manager) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := mr.managers.InsertOne(ctx, manager)
	if err != nil {
		return err
	}

	return nil
}

//
func (mr *managerRepository) AddMany(managers []domain.Manager) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mgrs := make([]interface{}, 0, len(managers))
	for _, m := range managers {
		mgrs = append(mgrs, m)
	}

	_, err := mr.managers.InsertMany(ctx, mgrs)
	if err != nil {
		return err
	}

	return nil
}

//
func (mr *managerRepository) UpdateInfo(managerID int, info domain.ManagerInfo) error {
	return storage.ErrManagerNotFound
}

//
func (mr *managerRepository) UpdateTeam(managerID int, team domain.Team) error {
	return nil // fmt.Errorf("not implemented")
}

//
func (mr *managerRepository) GetByID(id int) (domain.Manager, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := mr.managers.FindOne(ctx, bson.M{"_id": id})

	var manager domain.Manager
	err := result.Decode(&manager)
	if err != nil {
		return domain.Manager{}, err
	}

	return manager, nil
}
