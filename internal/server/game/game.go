package game

import (
	"fmt"
	"log"
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

	dayID     uint
	kickVotes map[uint]map[string]string
	killVotes map[uint]map[string]string

	// True when the current phase is day.
	dayPhase bool

	winners *Team
}

func New() *Game {
	return &Game{
		id:        uuid.New(),
		kickVotes: make(map[uint]map[string]string),
		killVotes: make(map[uint]map[string]string),
	}
}

func (g *Game) String() string {
	return fmt.Sprintf("[#%s]", g.id.String()[:6])
}

func (g *Game) ID() uuid.UUID {
	return g.id
}

func (g *Game) DayID() uint {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.dayID
}

func (g *Game) IsDayPhase() bool {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.dayPhase
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

// NewDay starts a new day. It returns who was killed the previous night and the next day id.
func (g *Game) NewDay() (*player.Player, uint) {
	g.mu.Lock()
	defer g.mu.Unlock()

	killed := g.findMostVoted(g.killVotes[g.dayID])
	if killed != nil {
		killed.Kill()
	}

	g.dayID++
	g.dayPhase = true

	return killed, g.dayID
}

// NewNight starts a night and returns who was kicked that day.
func (g *Game) NewNight() *player.Player {
	g.mu.Lock()
	defer g.mu.Unlock()

	kicked := g.findMostVoted(g.kickVotes[g.dayID])
	if kicked != nil {
		kicked.Kill()
	}

	g.dayPhase = false

	return kicked
}

func (g *Game) findMostVoted(votes map[string]string) *player.Player {
	cnt := make(map[string]uint)
	for _, candidate := range votes {
		cnt[candidate] = cnt[candidate] + 1
	}

	var winner string
	var topCount uint
	var isAbsolute bool
	for candidate, c := range cnt {
		log.Printf("%s has %d votes\n", candidate, c)

		if c == topCount {
			isAbsolute = false
		}
		if c > topCount {
			winner = candidate
			topCount = c
			isAbsolute = true
		}
	}

	if !isAbsolute {
		return nil
	}

	p, _ := g.findPlayerUnderLock(winner)

	return p
}

func (g *Game) AddKickVote(dayID uint, voter, target *player.Player) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.addVote(g.kickVotes, dayID, voter, target)

	log.Printf("%s %s voted to kick %s\n", g, voter.Username(), target.Username())
}

func (g *Game) AddKillVote(dayID uint, voter, target *player.Player) {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.addVote(g.killVotes, dayID, voter, target)

	log.Printf("%s %s voted to kill %s\n", g, voter.Username(), target.Username())
}

func (g *Game) addVote(m map[uint]map[string]string, dayID uint, voter, target *player.Player) {
	votes, exists := m[dayID]
	if !exists {
		votes = make(map[string]string)
		m[dayID] = votes
	}

	votes[voter.Username()] = target.Username()
}

func (g *Game) CheckStatus() {
	g.mu.Lock()
	defer g.mu.Unlock()

	if g.winners != nil {
		return
	}

	var aliveVillagers bool
	var aliveMafia bool
	for _, p := range g.players {
		if !p.Alive() {
			continue
		}

		switch p.Role() {
		case player.RoleInnocent, player.RoleSheriff:
			aliveVillagers = true
		case player.RoleMafiosi:
			aliveMafia = true
		}
	}

	if aliveVillagers && aliveMafia {
		return
	}

	var winners Team
	if aliveMafia {
		winners = TeamMafia
	} else {
		winners = TeamVillagers
	}

	g.winners = &winners
}

func (g *Game) Winners() *Team {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.winners
}

func (g *Game) FindPlayer(username string) (*player.Player, error) {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.findPlayerUnderLock(username)
}

func (g *Game) findPlayerUnderLock(username string) (*player.Player, error) {
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
		ok = ok || g.winners != nil
		ok = ok || !p.Alive()
		ok = ok || (g.dayPhase && sender.Alive())
		ok = ok || (!g.dayPhase && sender.Alive() && sender.Role() == player.RoleMafiosi && p.Role() == player.RoleMafiosi)

		if ok {
			receivers = append(receivers, p)
		}
	}

	return receivers
}

func (g *Game) Players() []*player.Player {
	g.mu.RLock()
	defer g.mu.RUnlock()

	return g.players
}
