package player

import (
	"sync"

	"github.com/lodthe/mafia/internal/server/session"
	"github.com/lodthe/mafia/pkg/mafiapb"
)

type Player struct {
	mu sync.RWMutex

	session *session.Session

	username string
	role     Role

	alive bool
}

func New(username string, role Role, s *session.Session) *Player {
	return &Player{
		session:  s,
		username: username,
		role:     role,
		alive:    true,
	}
}

func (p *Player) Session() *session.Session {
	return p.session
}

func (p *Player) Username() string {
	return p.username
}

func (p *Player) Role() Role {
	return p.role
}

func (p *Player) Alive() bool {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return p.alive
}

func (p *Player) Kill() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.alive = false
}

func (p *Player) Proto() *mafiapb.Player {
	p.mu.RLock()
	defer p.mu.RUnlock()

	return &mafiapb.Player{
		Username: p.username,
		Alive:    p.alive,
	}
}
