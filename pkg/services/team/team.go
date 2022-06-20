package team

import (
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/services/gameweek"
	"fpl-live-tracker/pkg/services/player"
	"fpl-live-tracker/pkg/storage"
	"fpl-live-tracker/pkg/wrapper"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"
)

const tripleCaptainActive = "3xc"
const benchBoostActive = "bboost"

// TeamService is an interface for interacting with teams
type TeamService interface {
	UpdateTeams() error
	UpdatePoints() error
	GetByID(id int) (domain.Team, error)
}

// teamService implements TeamService interface
type teamService struct {
	tr domain.TeamRepository
	gs gameweek.GameweekService
	ps player.PlayerService
	wr wrapper.Wrapper
}

// NewTeamService returns new instance of TeamService
func NewTeamService(tr domain.TeamRepository, gs gameweek.GameweekService,
	ps player.PlayerService, wr wrapper.Wrapper) TeamService {
	rand.Seed(time.Now().UnixNano())
	return &teamService{
		tr: tr,
		gs: gs,
		ps: ps,
		wr: wr,
	}
}

// UpdateTeams queries FPL API and updates all teams data
// (e.g. gameweek picks, chips played, hits taken), in its underlying team storage
func (ts *teamService) UpdateTeams() error {
	log.Println("team service: UpdateTeams started")
	// keep in mind, that number of teams can be lower than number of managers
	// due to gameweek deadlines
	inFplManagers, err := ts.wr.GetManagersCount()
	if err != nil {
		log.Println("team service:", err)
		return err
	}

	gameweek, err := ts.gs.GetCurrentGameweek()
	if err != nil {
		log.Println("team service:", err)
		return err
	}

	chanSize := runtime.NumCPU() * 4
	workerCount := runtime.NumCPU() * 16

	ids := make(chan int, chanSize)
	failed := make(chan int, chanSize)
	teams := make(chan wrapper.Team, chanSize)
	var workerWg sync.WaitGroup
	var innerWg sync.WaitGroup

	workerWg.Add(inFplManagers)

	for i := 0; i < workerCount; i++ {
		go func() {
			for id := range ids {
				wt, err := ts.wr.GetTeam(id, gameweek.ID)

				if herr, ok := err.(wrapper.ErrorHttpNotOk); ok {
					statusCode := herr.GetHttpStatusCode()
					switch statusCode {
					case http.StatusTooManyRequests:
						log.Println("team service: too many requests!")
						failed <- id
						time.Sleep(duration())
						continue
					case http.StatusServiceUnavailable:
						failed <- id
						time.Sleep(10 * time.Minute)
						continue
					case http.StatusNotFound:
						wt = wrapper.Team{
							ID:         id,
							ActiveChip: "not found",
						}
					default:
						failed <- id
						time.Sleep(10 * time.Minute)
						continue
					}
				} else if err != nil {
					failed <- id
					time.Sleep(10 * time.Minute)
					continue
				}

				teams <- wt
				workerWg.Done()
			}
		}()
	}

	innerWg.Add(1)
	go func() {
		// if team needs an update, send id to ids chan
		for id := 1; id <= inFplManagers; id++ {
			team, err := ts.tr.GetByID(id)
			if err != nil || team.GameweekID < gameweek.ID {
				ids <- id
			} else {
				// team already up-to-date, skipping
				workerWg.Done()
			}
		}
		innerWg.Done()
	}()

	innerWg.Add(1)
	go func() {
		// receive from failed chan, send to ids chan
		for id := range failed {
			ids <- id
		}
		innerWg.Done()
	}()

	innerWg.Add(1)
	go func() {
		// receive from teams chan
		for wt := range teams {
			dt, err := ts.convertToDomainTeam(wt)
			if err != nil {
				log.Println("team service: failed to convert team data")
			}

			err = ts.tr.Update(dt)
			if err == storage.ErrTeamNotFound {
				err = ts.tr.Add(dt)
				if err != nil {
					log.Println("team service: failed to add new team", err)
				}
			} else if err != nil {
				log.Println("team service: failed to update team", err)

			}
		}
		innerWg.Done()
	}()

	workerWg.Wait()

	close(ids)
	close(failed)
	close(teams)

	innerWg.Wait()

	log.Println("team service: UpdateTeams returned")
	return nil
}

// UpdatePoints updates points gained for all of the teams in the team storage
func (ts *teamService) UpdatePoints() error {
	log.Println("team service: UpdatePoints started")

	inStorageTeams, err := ts.tr.GetCount()
	if err != nil {
		log.Println("team service:", err)
		return err
	}

	chanSize := runtime.NumCPU()
	workerCount := runtime.NumCPU()

	ids := make(chan int, chanSize)
	failed := make(chan int, chanSize)
	var workerWg sync.WaitGroup
	var innerWg sync.WaitGroup

	workerWg.Add(inStorageTeams)
	for i := 0; i < workerCount; i++ {
		go func() {
			for id := range ids {
				err := ts.updateTeamPoints(id)
				if err != nil {
					failed <- id
					log.Println("team service: failed to update points", err)
					continue
				}
				workerWg.Done()
			}
		}()
	}

	innerWg.Add(1)
	go func() {
		// send to ids chan
		for id := 1; id <= inStorageTeams; id++ {
			ids <- id
		}
		innerWg.Done()
	}()

	innerWg.Add(1)
	go func() {
		// receive from failed chan, send to ids chan
		for id := range failed {
			ids <- id
		}
		innerWg.Done()
	}()

	workerWg.Wait()

	close(ids)
	close(failed)

	innerWg.Wait()

	log.Println("team service: UpdatePoints returned")
	return nil
}

// GetByID returns team with given ID, or an error otherwise
func (ts *teamService) GetByID(id int) (domain.Team, error) {
	team := domain.Team{ID: id}

	err := runTeamValidations(&team, idHigherThanZero)
	if err != nil {
		return domain.Team{}, err
	}

	return ts.tr.GetByID(id)
}

// convertToDomainTeam returns domain.Team, consistent with given wrapper.Team
func (ts *teamService) convertToDomainTeam(wt wrapper.Team) (domain.Team, error) {
	team := domain.Team{
		ID:         wt.ID,
		GameweekID: wt.EntryHistory.GameweekID,
		Picks:      make([]domain.TeamPlayer, 0, 15),
		ActiveChip: wt.ActiveChip,
		HitPoints:  wt.EntryHistory.EventTransfersCost,
	}

	for _, pick := range wt.Picks {
		p, err := ts.ps.GetByID(pick.ID)
		if err != nil {
			log.Println("team service:", err)
			return domain.Team{}, err
		}

		dp := domain.TeamPlayer{
			Player:        p,
			IsCaptain:     pick.IsCaptain,
			IsViceCaptain: pick.IsViceCaptain,
		}

		team.Picks = append(team.Picks, dp)
	}

	return team, nil
}

// updateTeamPoints updates points gained by manager's team with given ID
func (ts *teamService) updateTeamPoints(teamID int) error {
	team, err := ts.tr.GetByID(teamID)
	if err != nil {
		return err
	}

	err = ts.updateTeamPlayersStats(&team)
	if err != nil {
		return err
	}

	totalPoints := calculateTotalPoints(&team)
	subPoints := calculateSubPoints(&team)

	team.Points = totalPoints - team.HitPoints
	team.PointsAfterSubs = totalPoints + subPoints - team.HitPoints

	err = ts.tr.Update(team)
	if err != nil {
		return err
	}

	return nil
}

// updateTeamPlayersStats updates players stats in the given team
func (ts *teamService) updateTeamPlayersStats(team *domain.Team) error {
	for i := 0; i < len(team.Picks); i++ {
		tp := team.Picks[i]
		p, err := ts.ps.GetByID(tp.ID)
		if err != nil {
			log.Println("team service: failed to update team stats", err)
			return err
		}
		tp.Stats = p.Stats
		team.Picks[i] = tp
	}

	return nil
}

func duration() time.Duration {
	// random duration between 30s to 5min
	return (time.Duration(rand.Intn(270)) * time.Second) + 30
}
