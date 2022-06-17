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
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"gopkg.in/mgo.v2/bson"
)

// teamRepository implements domain.TeamRepository interface
type teamRepository struct {
	db    *mongo.Database
	teams *mongo.Collection
}

// NewTeamRepository returns new instance of domain.TeamRepository
func NewTeamRepository(config config.MongoConfig) (domain.TeamRepository, error) {
	client, err := mongo.Connect(context.Background(),
		options.Client().ApplyURI(fmt.Sprintf("mongodb://%s:%d", config.Host, config.Port)))
	if err != nil {
		return nil, err
	}
	// defer client.Disconnect(context.TODO())

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("team repository: failed to ping mongodb %w", err)
	}

	db := client.Database(config.Database)
	teams := db.Collection("team")

	return &teamRepository{
		db:    db,
		teams: teams,
	}, nil
}

// Add saves given team into mongo collection, or returns an error on failure
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

		return fmt.Errorf("storage: add record failed: %w", err)
	}

	return nil
}

// Update updates team with matching ID in mongo collection, of return an error on failure
func (tr *teamRepository) Update(team domain.Team) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	filter := bson.M{"_id": team.ID}

	result, err := tr.teams.ReplaceOne(ctx, filter, team)
	if err != nil {
		return fmt.Errorf("storage: update record failed: %w", err)
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
	if result.Err() == mongo.ErrNoDocuments {
		return domain.Team{}, storage.ErrTeamNotFound
	}

	var team domain.Team
	err := result.Decode(&team)
	if err != nil {
		return domain.Team{}, fmt.Errorf("storage: get record failed: %w", err)
	}

	return team, nil
}

// GetCount returns number of team records in mongo collection
func (tr *teamRepository) GetCount() (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	count, err := tr.teams.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, fmt.Errorf("storage: get record count failed: %w", err)
	}

	return int(count), nil
}
