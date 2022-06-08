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

const timeout = 10 * time.Second

// managerRepository implements domain.ManagerRepository interface
type managerRepository struct {
	db       *mongo.Database
	managers *mongo.Collection
}

// NewManagerRepository returns new instance of domain.ManagerRepository
func NewManagerRepository(config config.MongoConfig) (domain.ManagerRepository, error) {
	client, err := mongo.Connect(context.Background(),
		options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", config.Host, config.Port)))
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

// Add saves given manager into mongo collection, or returns error on failure
func (mr *managerRepository) Add(manager domain.Manager) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := mr.managers.InsertOne(ctx, manager)
	if err != nil {
		return err
	}

	return nil
}

// AddMany saves all given managers into mongo collection, or returns error on failure
func (mr *managerRepository) AddMany(managers []domain.Manager) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
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

// Update updates manager with matching ID in mongo collection, or returns error on failure
func (mr *managerRepository) UpdateInfo(managerID int, info domain.ManagerInfo) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	filter := bson.M{"_id": managerID}
	update := bson.M{
		"$set": bson.M{
			"ManagerInfo": info,
		},
	}

	result, err := mr.managers.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return storage.ErrManagerNotFound
	}

	return nil
}

// UpdateTeam updates team of manager with given ID, or returns error on failure
func (mr *managerRepository) UpdateTeam(managerID int, team domain.Team) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	filter := bson.M{"_id": managerID}
	update := bson.M{
		"$set": bson.M{
			"Team": team,
		},
	}

	result, err := mr.managers.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return storage.ErrManagerNotFound
	}

	return nil
}

// GetByID returns manager with given ID, or returns error on failure
func (mr *managerRepository) GetByID(id int) (domain.Manager, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result := mr.managers.FindOne(ctx, bson.M{"_id": id})

	var manager domain.Manager
	err := result.Decode(&manager)
	if err != nil {
		return domain.Manager{}, err
	}

	return manager, nil
}

// GetCount returns number of manager records in mongo collection
func (mr *managerRepository) GetCount() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	count, err := mr.managers.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}

	return int(count), nil
}
