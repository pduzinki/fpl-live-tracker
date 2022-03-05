package mongo

import (
	"fpl-live-tracker/pkg/config"
	"log"
	"testing"
)

func TestManagerAdd(t *testing.T) {
	// TODO add test

	mr, _ := NewManagerRepository(config.MongoConfig{})
	log.Println(mr)

}

func TestManagerAddMany(t *testing.T) {
	// TODO add test
}

func TestManagerUpdateInfo(t *testing.T) {
	// TODO add test
}

func TestManagerUpdateTeam(t *testing.T) {
	// TODO add test
}

func TestManagerGetByID(t *testing.T) {
	// TODO add test
}

func TestManagerGetCount(t *testing.T) {
	// TODO add test
}
