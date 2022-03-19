package game

import "github.com/lodthe/mafia/pkg/mafiapb"

type Team string

const (
	TeamVillagers Team = "VILLAGERS"
	TeamMafia     Team = "MAFIA"
)

func (t Team) Proto() mafiapb.Team {
	switch t {
	case TeamVillagers:
		return mafiapb.Team_TEAM_VILLAGERS

	case TeamMafia:
		return mafiapb.Team_TEAM_MAFIA

	default:
		return mafiapb.Team_TEAM_UNKNOWN
	}
}
