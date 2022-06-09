package team

import (
	"fpl-live-tracker/pkg/domain"
	"log"
)

//
func calculateTotalPoints(team *domain.Team) int {
	if len(team.Picks) == 0 {
		log.Println("calculateTotalPoints: team has no players!")
		return 0
	}

	captainMultiplier := 2
	playersCount := 11

	if team.ActiveChip == tripleCaptainActive {
		captainMultiplier = 3
	} else if team.ActiveChip == benchBoostActive {
		playersCount = 15
	}

	var totalPoints int
	for i := 0; i < playersCount; i++ {
		if team.Picks[i].IsCaptain {
			totalPoints += team.Picks[i].Stats.TotalPoints * captainMultiplier
		} else {
			totalPoints += team.Picks[i].Stats.TotalPoints
		}
	}

	return totalPoints
}

//
func calculateSubPoints(team *domain.Team) int {
	subPoints := 0

	if len(team.Picks) == 0 {
		log.Println("calculateSubPoints: team has no players!")
		return subPoints
	}
	if team.ActiveChip == "bboost" {
		return subPoints
	}

	// clear sub flags
	for i := 0; i < 15; i++ {
		team.Picks[i].SubIn = false
		// team.Picks[i].SubOut = false
	}

	liveFormation := getLiveFormation(team)
	subsIn := make([]domain.TeamPlayer, 0, 4)

	// check if goalkeeper needs a substitution
	if liveFormation[0] == 0 {
		benchGk := &team.Picks[11]
		if played(benchGk) {
			benchGk.SubIn = true
			subsIn = append(subsIn, *benchGk)
		}
	}

	// check if outfield players need substitutions
	bench := team.Picks[12:]
	subsNeeded := 10 - (liveFormation[1] + liveFormation[2] + liveFormation[3])

	for i := 0; i < subsNeeded; i++ {
		for j := 0; j < len(bench); j++ {
			b := &bench[j]
			if b.SubIn {
				continue
			}

			pos := b.Info.Position

			if liveFormation[1] < 3 { // too few defs, add only if b is def
				if pos == "DEF" && played(b) {
					b.SubIn = true
					subsIn = append(subsIn, *b)
					liveFormation[1]++
					break
				}
				continue
			}

			if liveFormation[3] < 1 { // too few fwds, add only if b is fwd
				if pos == "FWD" && played(b) {
					b.SubIn = true
					subsIn = append(subsIn, *b)
					liveFormation[3]++
					break
				}
				continue
			}

			if played(b) {
				b.SubIn = true
				subsIn = append(subsIn, *b)
				break
			}
		}
	}

	if !captainPlayed(team) {
		for i := 0; i < 11; i++ {
			if team.Picks[i].IsViceCaptain {
				if team.ActiveChip == tripleCaptainActive {
					subPoints += team.Picks[i].Stats.TotalPoints * 2
				} else {
					subPoints += team.Picks[i].Stats.TotalPoints
				}
			}
		}
	}

	for _, s := range subsIn {
		subPoints += s.Stats.TotalPoints
	}

	return subPoints
}

//
func getLiveFormation(team *domain.Team) [4]int {
	var gkps, defs, mids, fwds int

	for _, p := range team.Picks[:11] {
		pos := p.Info.Position
		if pos == "GKP" && played(&p) {
			gkps++
		} else if pos == "DEF" && played(&p) {
			defs++
		} else if pos == "MID" && played(&p) {
			mids++
		} else if pos == "FWD" && played(&p) {
			fwds++
		}
	}

	return [4]int{gkps, defs, mids, fwds}
}

// played returns boolean value, indicating if given player
// had played some minutes in his fixtures
func played(player *domain.TeamPlayer) bool {
	stats := player.Stats

	var allFixturesStarted bool = true
	for _, f := range stats.FixturesInfo {
		if !f.Started {
			allFixturesStarted = false
			break
		}
	}

	if stats.Minutes == 0 && allFixturesStarted {
		return false
	}
	return true
}

//
func captainPlayed(team *domain.Team) bool {
	for _, p := range team.Picks {
		if p.IsCaptain {
			return played(&p)
		}
	}
	return false
}
