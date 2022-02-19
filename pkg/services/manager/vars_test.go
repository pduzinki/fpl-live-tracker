package manager

import "fpl-live-tracker/pkg/domain"

var (
	john  = domain.Manager{ID: 1, Info: domain.ManagerInfo{Name: "John", TeamName: "John FC"}}
	picks = []domain.TeamPlayer{
		{
			Player: domain.Player{
				Stats: domain.PlayerStats{
					FixturesInfo: []domain.FixtureInfo{{Started: true, FinishedProvisional: true, Finished: true}},
					Minutes:      90,
					TotalPoints:  2,
				},
			},
		},
		{
			Player: domain.Player{
				Stats: domain.PlayerStats{
					FixturesInfo: []domain.FixtureInfo{{Started: true, FinishedProvisional: true, Finished: true}},
					Minutes:      90,
					TotalPoints:  15,
				},
			},
			IsCaptain: true,
		},
		{
			Player: domain.Player{
				Stats: domain.PlayerStats{
					FixturesInfo: []domain.FixtureInfo{{Started: true, FinishedProvisional: true, Finished: true}},
					Minutes:      90,
					TotalPoints:  2,
				},
			},
		},
		{
			Player: domain.Player{
				Stats: domain.PlayerStats{
					FixturesInfo: []domain.FixtureInfo{{Started: true, FinishedProvisional: true, Finished: true}},
					Minutes:      90,
					TotalPoints:  2,
				},
			},
		},
		{
			Player: domain.Player{
				Stats: domain.PlayerStats{
					FixturesInfo: []domain.FixtureInfo{{Started: true, FinishedProvisional: true, Finished: true}},
					Minutes:      90,
					TotalPoints:  8,
				},
			},
		},
		{
			Player: domain.Player{
				Stats: domain.PlayerStats{
					FixturesInfo: []domain.FixtureInfo{{Started: true, FinishedProvisional: true, Finished: true}},
					Minutes:      90,
					TotalPoints:  9,
				},
			},
		},
		{
			Player: domain.Player{
				Stats: domain.PlayerStats{
					FixturesInfo: []domain.FixtureInfo{{Started: true, FinishedProvisional: true, Finished: true}},
					Minutes:      90,
					TotalPoints:  5,
				},
			},
		},
		{
			Player: domain.Player{
				Stats: domain.PlayerStats{
					FixturesInfo: []domain.FixtureInfo{{Started: true, FinishedProvisional: true, Finished: true}},
					Minutes:      90,
					TotalPoints:  -2,
				},
			},
		},
		{
			Player: domain.Player{
				Stats: domain.PlayerStats{
					FixturesInfo: []domain.FixtureInfo{{Started: true, FinishedProvisional: true, Finished: true}},
					Minutes:      90,
					TotalPoints:  1,
				},
			},
		},
		{
			Player: domain.Player{
				Stats: domain.PlayerStats{
					FixturesInfo: []domain.FixtureInfo{{Started: true, FinishedProvisional: true, Finished: true}},
					Minutes:      90,
					TotalPoints:  2,
				},
			},
		},
		{
			Player: domain.Player{
				Stats: domain.PlayerStats{
					FixturesInfo: []domain.FixtureInfo{{Started: true, FinishedProvisional: true, Finished: true}},
					Minutes:      90,
					TotalPoints:  7,
				},
			},
		},
		// bench
		{
			Player: domain.Player{
				Stats: domain.PlayerStats{
					FixturesInfo: []domain.FixtureInfo{{Started: true, FinishedProvisional: true, Finished: true}},
					Minutes:      90,
					TotalPoints:  8,
				},
			},
		},
		{
			Player: domain.Player{
				Stats: domain.PlayerStats{
					FixturesInfo: []domain.FixtureInfo{{Started: true, FinishedProvisional: true, Finished: true}},
					Minutes:      90,
					TotalPoints:  8,
				},
			},
		},
		{
			Player: domain.Player{
				Stats: domain.PlayerStats{
					FixturesInfo: []domain.FixtureInfo{{Started: true, FinishedProvisional: true, Finished: true}},
					Minutes:      0,
					TotalPoints:  0,
				},
			},
		},
		{
			Player: domain.Player{
				Stats: domain.PlayerStats{
					FixturesInfo: []domain.FixtureInfo{{Started: true, FinishedProvisional: true, Finished: true}},
					Minutes:      90,
					TotalPoints:  1,
				},
			},
		},
	}

	noChipTeam = domain.Team{
		Picks:      picks,
		ActiveChip: "",
	}
	tripleCaptainTeam = domain.Team{
		Picks:      picks,
		ActiveChip: tripleCaptainActive,
	}
	benchBoostTeam = domain.Team{
		Picks:      picks,
		ActiveChip: benchBoostActive,
	}
)
