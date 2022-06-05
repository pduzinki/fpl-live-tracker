package mongo

import (
	"fmt"
	"fpl-live-tracker/pkg/config"
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/storage"
	"log"
	"os"
	"testing"

	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
)

var (
	jimsTeam  = domain.Team{ID: 1}
	janesTeam = domain.Team{ID: 2}
	joelsTeam = domain.Team{ID: 3}
	jacksTeam = domain.Team{ID: 4}
)

var tr domain.TeamRepository

func TestMain(m *testing.M) {
	fmt.Println("testmain")
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("err 1")
	}

	opts := dockertest.RunOptions{
		Repository:   "mongo",
		Tag:          "latest",
		ExposedPorts: []string{"27017"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"27017": {
				{HostIP: "0.0.0.0", HostPort: "27017"},
			},
		},
	}

	resource, err := pool.RunWithOptions(&opts)
	if err != nil {
		log.Fatalf("err 2 %v", err)
	}

	if err = pool.Retry(func() error {
		config := config.MongoConfig{
			Host:     "localhost",
			Port:     27017,
			Database: "fpl-live-tracker",
		}

		tr, err = NewTeamRepository(config)
		if err != nil {
			log.Fatalf("err 5")
		}

		_ = tr

		err = tr.Add(domain.Team{ID: 1239, ActiveChip: "3xc"})
		log.Println("err:", err)

		team, err := tr.GetByID(1239)
		log.Println(team, err)

		return nil
	}); err != nil {
		log.Fatalf("err 3")
	}

	// seed data
	tr.Add(jimsTeam)

	code := m.Run()

	if err = pool.Purge(resource); err != nil {
		log.Fatalf("err 4")
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
	fmt.Println("b")
	// TODO add test
}

func TestTeamGetByID(t *testing.T) {
	fmt.Println("c")
	// TODO add test
}

func TestTeamGetCount(t *testing.T) {
	fmt.Println("d")
	// TODO add test
}
