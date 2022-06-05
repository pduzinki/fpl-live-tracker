package mongo

import (
	"context"
	"errors"
	"fmt"
	"fpl-live-tracker/pkg/config"
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/storage"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

//
type teamRepository struct {
	db    *mongo.Database
	teams *mongo.Collection
}

//
func NewTeamRepository(config config.MongoConfig) (domain.TeamRepository, error) {
	client, err := mongo.Connect(context.Background(),
		options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", config.Host, config.Port)))
	if err != nil {
		return nil, err
	}
	// defer client.Disconnect(context.TODO())

	db := client.Database(config.Database)
	teams := db.Collection("team")

	return &teamRepository{
		db:    db,
		teams: teams,
	}, nil
}

//
func (tr *teamRepository) Add(team domain.Team) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	_, err := tr.teams.InsertOne(ctx, team)
	if err != nil {
		var werr mongo.WriteException
		if errors.As(err, &werr) {
			if len(werr.WriteErrors) > 0 && werr.WriteErrors[0].Code == 11000 {
				return storage.ErrTeamAlreadyExists
			}
		}

		return storage.ErrAddRecordFailed
	}

	return nil
}

//
func (tr *teamRepository) Update(ID int, team domain.Team) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	filter := bson.M{"_id": ID}

	result, err := tr.teams.ReplaceOne(ctx, filter, team)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return storage.ErrTeamNotFound
	}

	return nil
}

//
func (tr *teamRepository) GetByID(ID int) (domain.Team, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	result := tr.teams.FindOne(ctx, bson.M{"_id": ID})

	var team domain.Team
	err := result.Decode(&team)
	if err != nil {
		return domain.Team{}, err
	}

	return team, nil
}

//
func (tr *teamRepository) GetCount() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	count, err := tr.teams.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, err
	}

	return int(count), nil
}
