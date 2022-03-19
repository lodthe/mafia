package game

import (
	"fmt"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/lodthe/mafia/internal/server/player"
	"github.com/pkg/errors"
)

type Game struct {
	id uuid.UUID

	mu sync.RWMutex

	players []*player.Player

	dayID      uint
	dayVotes   map[uint]map[string]string
	nightVotes map[uint]map[string]string

	// True when the current phase is day.
	dayPhase bool

	winners *Team
}

func New() *Game {
	return &Game{
		id:         uuid.New(),
		dayVotes:   make(map[uint]map[string]string),
		nightVotes: make(map[uint]map[string]string),
	}
}

func (g *Game) String() string {
	return fmt.Sprintf("#%s", string(g.id[:6]))
}

func (g *Game) ID() uuid.UUID {
	return g.id
}

func (g *Game) AddPlayer(p *player.Player) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	for _, other := range g.players {
		if strings.EqualFold(other.Username(), p.Username()) {
			return errors.New("player with username has already joined the game")
		}
	}

	g.players = append(g.players, p)

	return nil
}

func (g *Game) NewDay() uint {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.dayID++

	return g.dayID
}

func (g *Game) FindPlayer(username string) (*player.Player, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for _, p := range g.players {
		if strings.EqualFold(p.Username(), username) {
			return p, nil
		}
	}

	return nil, errors.New("player not found")
}

func (g *Game) FindMessageReceivers(sender *player.Player) (receivers []*player.Player) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	for _, p := range g.players {
		if strings.EqualFold(sender.Username(), p.Username()) {
			continue
		}

		var ok bool
		ok = ok || g.dayID == 0
		ok = ok || g.winners == nil
		ok = ok || !p.Alive()
		ok = ok || (g.dayPhase && sender.Alive())
		ok = ok || (!g.dayPhase && sender.Alive() && sender.Role() == player.RoleMafiosi && p.Role() == player.RoleMafiosi)

		if ok {
			receivers = append(receivers, p)
		}
	}

	return receivers
}
