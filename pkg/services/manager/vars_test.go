package manager

import "fpl-live-tracker/pkg/domain"

var (
	john = domain.Manager{ID: 1, Info: domain.ManagerInfo{Name: "John", TeamName: "John FC"}}

	ramsdale = domain.TeamPlayer{Player: domain.Player{
		Info: domain.PlayerInfo{Position: "GKP"},
		Stats: domain.PlayerStats{
			FixturesInfo: []domain.FixtureInfo{{Started: true, FinishedProvisional: true, Finished: true}},
			Minutes:      90,
			TotalPoints:  2,
		},
	}}
	taa = domain.TeamPlayer{Player: domain.Player{
		Info: domain.PlayerInfo{Position: "DEF"},
		Stats: domain.PlayerStats{
			FixturesInfo: []domain.FixtureInfo{{Started: true, FinishedProvisional: true, Finished: true}},
			Minutes:      90,
			TotalPoints:  15,
		},
	}}
	cancelo = domain.TeamPlayer{Player: domain.Player{
		Info: domain.PlayerInfo{Position: "DEF"},
		Stats: domain.PlayerStats{
			FixturesInfo: []domain.FixtureInfo{{Started: true, FinishedProvisional: true, Finished: true}},
			Minutes:      90,
			TotalPoints:  2,
		},
	}}
	kilman = domain.TeamPlayer{Player: domain.Player{
		Info: domain.PlayerInfo{Position: "DEF"},
		Stats: domain.PlayerStats{
			FixturesInfo: []domain.FixtureInfo{{Started: true, FinishedProvisional: true, Finished: true}},
			Minutes:      90,
			TotalPoints:  2,
		},
	}}
	salah = domain.TeamPlayer{Player: domain.Player{
		Info: domain.PlayerInfo{Position: "MID"},
		Stats: domain.PlayerStats{
			FixturesInfo: []domain.FixtureInfo{{Started: true, FinishedProvisional: true, Finished: true}},
			Minutes:      90,
			TotalPoints:  8,
		},
	}}
	saka = domain.TeamPlayer{Player: domain.Player{
		Info: domain.PlayerInfo{Position: "MID"},
		Stats: domain.PlayerStats{
			FixturesInfo: []domain.FixtureInfo{{Started: true, FinishedProvisional: true, Finished: true}},
			Minutes:      90,
			TotalPoints:  9,
		},
	}}
	bowen = domain.TeamPlayer{Player: domain.Player{
		Info: domain.PlayerInfo{Position: "MID"},
		Stats: domain.PlayerStats{
			FixturesInfo: []domain.FixtureInfo{{Started: true, FinishedProvisional: true, Finished: true}},
			Minutes:      90,
			TotalPoints:  5,
		},
	}}
	son = domain.TeamPlayer{Player: domain.Player{
		Info: domain.PlayerInfo{Position: "MID"},
		Stats: domain.PlayerStats{
			FixturesInfo: []domain.FixtureInfo{{Started: true, FinishedProvisional: true, Finished: true}},
			Minutes:      90,
			TotalPoints:  -2,
		},
	}}
	dennis = domain.TeamPlayer{Player: domain.Player{
		Info: domain.PlayerInfo{Position: "FWD"},
		Stats: domain.PlayerStats{
			FixturesInfo: []domain.FixtureInfo{{Started: true, FinishedProvisional: true, Finished: true}},
			Minutes:      90,
			TotalPoints:  1,
		},
	}}
	broja = domain.TeamPlayer{Player: domain.Player{
		Info: domain.PlayerInfo{Position: "FWD"},
		Stats: domain.PlayerStats{
			FixturesInfo: []domain.FixtureInfo{{Started: true, FinishedProvisional: true, Finished: true}},
			Minutes:      90,
			TotalPoints:  2,
		},
	}}
	antonio = domain.TeamPlayer{Player: domain.Player{
		Info: domain.PlayerInfo{Position: "FWD"},
		Stats: domain.PlayerStats{
			FixturesInfo: []domain.FixtureInfo{{Started: true, FinishedProvisional: true, Finished: true}},
			Minutes:      90,
			TotalPoints:  7,
		},
	}}
	// bench
	sanchez = domain.TeamPlayer{Player: domain.Player{
		Info: domain.PlayerInfo{Position: "GKP"},
		Stats: domain.PlayerStats{
			FixturesInfo: []domain.FixtureInfo{{Started: true, FinishedProvisional: true, Finished: true}},
			Minutes:      90,
			TotalPoints:  8,
		},
	}}
	gilmour = domain.TeamPlayer{Player: domain.Player{
		Info: domain.PlayerInfo{Position: "MID"},
		Stats: domain.PlayerStats{
			FixturesInfo: []domain.FixtureInfo{{Started: true, FinishedProvisional: true, Finished: true}},
			Minutes:      90,
			TotalPoints:  8,
		},
	}}
	livramento = domain.TeamPlayer{Player: domain.Player{
		Info: domain.PlayerInfo{Position: "DEF"},
		Stats: domain.PlayerStats{
			FixturesInfo: []domain.FixtureInfo{{Started: true, FinishedProvisional: true, Finished: true}},
			Minutes:      90,
			TotalPoints:  1,
		},
	}}
	lamptey = domain.TeamPlayer{Player: domain.Player{
		Info: domain.PlayerInfo{Position: "DEF"},
		Stats: domain.PlayerStats{
			FixturesInfo: []domain.FixtureInfo{{Started: true, FinishedProvisional: true, Finished: true}},
			Minutes:      90,
			TotalPoints:  1,
		},
	}}

	teamA = domain.Team{
		Picks: []domain.TeamPlayer{
			ramsdale,
			cpt(taa), cancelo, kilman,
			salah, saka, bowen, son,
			dennis, broja, antonio,
			sanchez,
			gilmour, livramento, lamptey,
		},
	}

	teamB = domain.Team{
		Picks: []domain.TeamPlayer{
			sub(ramsdale),
			cpt(taa), cancelo, kilman,
			salah, saka, sub(bowen), son,
			dennis, broja, antonio,
			sanchez,
			gilmour, livramento, lamptey,
		},
	}

	teamC = domain.Team{
		Picks: []domain.TeamPlayer{
			ramsdale,
			sub(cpt(taa)), cancelo, kilman,
			salah, saka, bowen, son,
			dennis, broja, vcpt(antonio),
			sanchez,
			gilmour, livramento, lamptey,
		},
	}

	teamD = domain.Team{
		Picks: []domain.TeamPlayer{
			sub(ramsdale),
			sub(taa), cancelo, sub(kilman),
			cpt(salah), saka, sub(bowen), son,
			dennis, sub(broja), vcpt(antonio),
			sanchez,
			gilmour, livramento, lamptey,
		},
	}

	teamE = domain.Team{
		Picks: []domain.TeamPlayer{
			ramsdale,
			sub(taa), cancelo, sub(kilman),
			cpt(salah), saka, bowen, son,
			dennis, broja, vcpt(antonio),
			sanchez,
			gilmour, livramento, lamptey,
		},
	}

	teamF = domain.Team{
		Picks: []domain.TeamPlayer{
			ramsdale,
			taa, cancelo, kilman, livramento, lamptey,
			cpt(salah), saka, bowen, son,
			sub(broja),
			sanchez,
			gilmour, dennis, antonio,
		},
	}

	teamG = domain.Team{
		Picks: []domain.TeamPlayer{
			ramsdale,
			sub(taa), sub(cancelo), sub(kilman),
			cpt(salah), saka, bowen, son,
			dennis, broja, antonio,
			sanchez,
			gilmour, sub(livramento), sub(lamptey),
		},
	}

	teamH = domain.Team{
		Picks: []domain.TeamPlayer{
			ramsdale,
			cpt(taa), cancelo, kilman,
			salah, saka, bowen, son,
			sub(dennis), sub(broja), sub(antonio),
			sanchez,
			gilmour, livramento, lamptey,
		},
	}
)

func cpt(p domain.TeamPlayer) domain.TeamPlayer {
	p.IsCaptain = true
	return p
}

func vcpt(p domain.TeamPlayer) domain.TeamPlayer {
	p.IsViceCaptain = true
	return p
}

func sub(p domain.TeamPlayer) domain.TeamPlayer {
	p.Stats.Minutes = 0
	p.Stats.TotalPoints = 0
	return p
}

func chipTripleCpt(t domain.Team) domain.Team {
	t.ActiveChip = tripleCaptainActive
	return t
}

func chipBenchBoost(t domain.Team) domain.Team {
	t.ActiveChip = benchBoostActive
	return t
}
