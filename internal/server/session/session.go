package session

import (
	"github.com/google/uuid"
	"github.com/lodthe/mafia/pkg/mafiapb"
)

type Session struct {
	ID uuid.UUID

	GameID   uuid.UUID
	Username string

	events chan<- *mafiapb.GameEvent
}

func New(gameID uuid.UUID, username string, events chan<- *mafiapb.GameEvent) *Session {
	return &Session{
		ID:       uuid.New(),
		GameID:   gameID,
		Username: username,
		events:   events,
	}
}

func (s *Session) SendNonBlocking(event *mafiapb.GameEvent) bool {
	select {
	case s.events <- event:
		return true

	default:
		return false
	}
}
