package mongo

import (
	"fmt"
	"fpl-live-tracker/pkg/config"
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/storage"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
)

const hostPort = 27018

var (
	jimsTeam  = domain.Team{ID: 1}
	janesTeam = domain.Team{ID: 2}
	joelsTeam = domain.Team{ID: 3}
	jacksTeam = domain.Team{ID: 4}
)

var tr domain.TeamRepository

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("failed to create new pool: %v", err)
	}

	opts := dockertest.RunOptions{
		Repository:   "mongo",
		Tag:          "latest",
		ExposedPorts: []string{"27017"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"27017": {
				{HostIP: "0.0.0.0", HostPort: fmt.Sprint(hostPort)},
			},
		},
	}

	resource, err := pool.RunWithOptions(&opts)
	if err != nil {
		log.Fatalf("failed to start the container: %v", err)
	}

	if err = pool.Retry(func() error {
		config := config.MongoConfig{
			Host:     "localhost",
			Port:     hostPort,
			Database: "fpl-live-tracker",
		}

		tr, err = NewTeamRepository(config)
		if err != nil {
			log.Fatalf("failed to create new team repository: %v", err)
		}

		mr, err = NewManagerRepository(config)
		if err != nil {
			log.Fatalf("failed to create new manager repository: %v", err)
		}

		return nil
	}); err != nil {
		log.Fatalf("failed to connect to the container: %v", err)
	}

	// seed data
	tr.Add(jimsTeam)

	mr.Add(john)

	code := m.Run()

	if err = pool.Purge(resource); err != nil {
		log.Fatalf("failed to remove the container: %v", err)
	}

	os.Exit(code)
}

func TestTeamAdd(t *testing.T) {
	testcases := []struct {
		team domain.Team
		want error
	}{
		{joelsTeam, nil},
		{jimsTeam, storage.ErrTeamAlreadyExists},
	}

	for _, test := range testcases {
		got := tr.Add(test.team)
		if got != test.want {
			t.Errorf("error: for %v, got err '%v', want '%v'", test.team, got, test.want)
		}
	}
}

func TestTeamUpdate(t *testing.T) {
	testcases := []struct {
		team domain.Team
		want error
	}{
		{
			team: domain.Team{
				ID: jimsTeam.ID,
			},
			want: nil,
		},
		{
			team: domain.Team{
				ID: 420,
			},
			want: storage.ErrTeamNotFound,
		},
	}

	for _, test := range testcases {
		got := tr.Update(test.team)
		if got != test.want {
			t.Errorf("error: got err '%v', want '%v'", got, test.want)
		}
	}
}

func TestTeamGetByID(t *testing.T) {
	testcases := []struct {
		id      int
		want    domain.Team
		wantErr error
	}{
		{jimsTeam.ID, jimsTeam, nil},
		{jacksTeam.ID, domain.Team{}, storage.ErrTeamNotFound},
	}

	for _, test := range testcases {
		got, gotErr := tr.GetByID(test.id)
		if gotErr != test.wantErr {
			t.Errorf("error: for %v, got err '%v', want err '%v'", test.id, gotErr, test.wantErr)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("error: for %v, got '%v', want '%v'", test.id, got, test.want)
		}
	}
}

func TestTeamGetCount(t *testing.T) {
	testcases := []struct {
		want    int
		wantErr error
	}{
		{
			want:    2,
			wantErr: nil,
		},
	}

	for _, test := range testcases {
		got, gotErr := tr.GetCount()
		if gotErr != test.wantErr {
			t.Errorf("error: got err '%v', want err '%v'", gotErr, test.wantErr)
		}
		if got != test.want {
			t.Errorf("error: got '%v', want '%v'", got, test.want)
		}
	}
}
