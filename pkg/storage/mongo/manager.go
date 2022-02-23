package mongo

import (
	"context"
	"fpl-live-tracker/pkg/config"
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/storage"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
)

const uri = "mongodb://localhost:27017"

//
type managerRepository struct {
	db       *mongo.Database
	managers *mongo.Collection
}

//
type mongoManager struct {
	ID       int    `bson:"id"`
	Name     string `bson:"name"`
	TeamName string `bson:"teamname"`
}

//
func NewManagerRepository(config config.MongoConfig) (domain.ManagerRepository, error) {
	// TODO use passed config
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}
	// defer client.Disconnect(context.TODO())

	db := client.Database("fpl-live-tracker")
	managers := db.Collection("managers")

	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		return nil, err
	}

	return &managerRepository{
		db:       db,
		managers: managers,
	}, nil
}

//
func (mr *managerRepository) Add(manager domain.Manager) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	m := mongoManager{
		ID:       manager.ID,
		Name:     manager.Info.Name,
		TeamName: manager.Info.TeamName,
	}

	_, err := mr.managers.InsertOne(ctx, m)
	if err != nil {
		return err
	}
	return nil
}

//
func (mr *managerRepository) AddMany(managers []domain.Manager) error {
	return nil // fmt.Errorf("not implemented")
}

//
func (mr *managerRepository) UpdateInfo(managerID int, info domain.ManagerInfo) error {
	return storage.ErrManagerNotFound
	// return nil // fmt.Errorf("not implemented")
}

//
func (mr *managerRepository) UpdateTeam(managerID int, team domain.Team) error {
	return nil // fmt.Errorf("not implemented")
}

//
func (mr *managerRepository) GetByID(id int) (domain.Manager, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result := mr.managers.FindOne(ctx, bson.M{"id": id})

	var m mongoManager
	err := result.Decode(&m)
	if err != nil {
		return domain.Manager{}, err
	}

	return domain.Manager{
		ID: m.ID,
		Info: domain.ManagerInfo{
			Name:     m.Name,
			TeamName: m.TeamName,
		},
	}, nil
}
