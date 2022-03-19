package server

import (
	"time"

	"github.com/lodthe/mafia/internal/server/player"
)

type Config struct {
	RoleDistribution map[player.Role]uint

	DayDuration   time.Duration
	NightDuration time.Duration
}

var DefaultConfig = Config{
	RoleDistribution: map[player.Role]uint{
		player.RoleInnocent: 2,
		player.RoleMafiosi:  1,
		player.RoleSheriff:  1,
	},
	DayDuration:   20 * time.Second,
	NightDuration: 15 * time.Second,
}

func (c *Config) Players() uint {
	var players uint
	for _, cnt := range c.RoleDistribution {
		players += cnt
	}

	return players
}
