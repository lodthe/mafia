package server

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/lodthe/mafia/internal/server/game"
	"github.com/lodthe/mafia/internal/server/player"
	"github.com/lodthe/mafia/internal/server/session"
	"github.com/lodthe/mafia/pkg/mafiapb"
	"github.com/pkg/errors"
)

type Engine struct {
	ctx    context.Context
	config Config

	sessionsLocker sync.RWMutex
	sessions       map[uuid.UUID]*session.Session

	gamesLocker sync.RWMutex
	games       map[uuid.UUID]*game.Game

	latestGameLocker sync.Mutex
	latestGame       *game.Game
	playersLeft      uint
	rolesToSelect    map[player.Role]uint
}

func NewEngine(ctx context.Context, config Config) (*Engine, error) {
	if config.RoleDistribution[player.RoleMafiosi] == 0 {
		return nil, errors.New("impossible to host games without mafiosi")
	}

	return &Engine{
		ctx:      ctx,
		config:   config,
		sessions: make(map[uuid.UUID]*session.Session),
		games:    make(map[uuid.UUID]*game.Game),
	}, nil
}

func (e *Engine) AddPlayer(username string, events chan<- *mafiapb.GameEvent) (*session.Session, error) {
	e.latestGameLocker.Lock()
	defer e.latestGameLocker.Unlock()

	if e.latestGame == nil {
		e.hostNewGame()
	}

	// Select a role.
	var role player.Role

	rnd := uint(rand.Intn(int(e.playersLeft)))
	for r, cntLeft := range e.rolesToSelect {
		if rnd < cntLeft {
			role = r
			break
		}

		rnd -= cntLeft
	}

	sess := session.New(e.latestGame.ID(), username, events)
	p := player.New(username, role, sess)
	err := e.latestGame.AddPlayer(p)
	if err != nil {
		return nil, err
	}

	e.addSession(sess)

	e.rolesToSelect[role] = e.rolesToSelect[role] - 1
	e.playersLeft--

	defer e.broadcast(e.latestGame, &mafiapb.GameEvent{
		Type: mafiapb.GameEvent_EVENT_PLAYER_JOINED,
		Payload: &mafiapb.GameEvent_PayloadPlayerJoined_{
			PayloadPlayerJoined: &mafiapb.GameEvent_PayloadPlayerJoined{
				Player: p.Proto(),
			},
		},
	})

	// If there are no free slots, start a game.
	if e.playersLeft == 0 {
		g := e.latestGame
		e.latestGame = nil

		go e.startGame(g)
	}

	return sess, err
}

func (e *Engine) RemovePlayer(s *session.Session) {
	p, err := e.FindPlayerBySessionID(s.ID)
	if err != nil {
		return
	}

	g, err := e.FindGameByID(s.GameID)
	if err != nil {
		return
	}

	p.Kill()
	g.CheckStatus()

	e.broadcast(e.latestGame, &mafiapb.GameEvent{
		Type: mafiapb.GameEvent_EVENT_PLAYER_LEFT,
		Payload: &mafiapb.GameEvent_PayloadPlayerLeft_{
			PayloadPlayerLeft: &mafiapb.GameEvent_PayloadPlayerLeft{
				Player: p.Proto(),
			},
		},
	})
}

func (e *Engine) hostNewGame() {
	e.latestGame = game.New()
	e.playersLeft = 0
	e.rolesToSelect = make(map[player.Role]uint)

	for role, cnt := range e.config.RoleDistribution {
		e.playersLeft += cnt
		e.rolesToSelect[role] = cnt
	}

	e.addGame(e.latestGame)
}

func (e *Engine) broadcast(g *game.Game, event *mafiapb.GameEvent) {
	for _, p := range g.Players() {
		p.Session().SendNonBlocking(event)
	}
}

func (e *Engine) startGame(g *game.Game) {
	log.Printf("%s game stared\n", g)

	for {
		killedPlayer, dayID := g.NewDay()

		log.Printf("%s day %d started\n", g, dayID)
		if killedPlayer != nil {
			log.Printf("%s %s was killed\n", g, killedPlayer.Username())
		}

		e.broadcast(g, &mafiapb.GameEvent{
			Type: mafiapb.GameEvent_EVENT_DAY_STARTED,
			Payload: &mafiapb.GameEvent_PayloadDayStarted_{
				PayloadDayStarted: &mafiapb.GameEvent_PayloadDayStarted{
					DayId:        uint64(dayID),
					KilledPlayer: killedPlayer.Proto(),
				},
			},
		})

		g.CheckStatus()
		if g.Winners() != nil {
			e.endGame(g)
			return
		}

		time.Sleep(e.config.DayDuration)

		kickedPlayer := g.NewNight()

		log.Printf("%s night %d started\n", g, dayID)
		if kickedPlayer != nil {
			log.Printf("%s %s was kicked\n", g, kickedPlayer.Username())
		}

		e.broadcast(g, &mafiapb.GameEvent{
			Type: mafiapb.GameEvent_EVENT_NIGHT_STARTED,
			Payload: &mafiapb.GameEvent_PayloadNightStarted_{
				PayloadNightStarted: &mafiapb.GameEvent_PayloadNightStarted{
					DayId:        uint64(dayID),
					KickedPlayer: kickedPlayer.Proto(),
				},
			},
		})

		g.CheckStatus()
		if g.Winners() != nil {
			e.endGame(g)
			return
		}

		time.Sleep(e.config.NightDuration)
	}
}

func (e *Engine) endGame(g *game.Game) {
	log.Printf("%s game finished\n", g)

	e.broadcast(g, &mafiapb.GameEvent{
		Type: mafiapb.GameEvent_EVENT_GAME_FINISHED,
		Payload: &mafiapb.GameEvent_PayloadGameFinished_{
			PayloadGameFinished: &mafiapb.GameEvent_PayloadGameFinished{
				Winners: g.Winners().Proto(),
				Players: nil,
			},
		},
	})
}

func (e *Engine) SendMessage(sender *player.Player, content string) ([]*player.Player, error) {
	g, err := e.FindGameByID(sender.Session().GameID)
	if err != nil {
		return nil, err
	}

	senderpb := sender.Proto()

	candidates := g.FindMessageReceivers(sender)
	receivers := make([]*player.Player, 0, len(candidates))
	for _, p := range candidates {
		msg := &mafiapb.GameEvent{
			Type: mafiapb.GameEvent_EVENT_MESSAGE,
			Payload: &mafiapb.GameEvent_PayloadMessage_{
				PayloadMessage: &mafiapb.GameEvent_PayloadMessage{
					Sender:  senderpb,
					Content: content,
				},
			},
		}
		if p.Session().SendNonBlocking(msg) {
			receivers = append(receivers, p)
		}
	}

	return receivers, nil
}

func (e *Engine) VoteKick(voter *player.Player, candidate string) error {
	g, err := e.FindGameByID(voter.Session().GameID)
	if err != nil {
		return err
	}

	if g.DayID() == 0 {
		return errors.New("game not started")
	}

	if g.Winners() != nil {
		return errors.New("game finished")
	}

	if !voter.Alive() {
		return errors.New("you are dead")
	}

	if !g.IsDayPhase() {
		return errors.New("it's night now")
	}

	target, err := g.FindPlayer(candidate)
	if err != nil {
		return errors.New("invalid username")
	}

	g.AddKickVote(g.DayID(), voter, target)

	return nil
}

func (e *Engine) KillVote(voter *player.Player, candidate string) error {
	g, err := e.FindGameByID(voter.Session().GameID)
	if err != nil {
		return err
	}

	if g.DayID() == 0 {
		return errors.New("game not started")
	}

	if g.Winners() != nil {
		return errors.New("game finished")
	}

	if !voter.Alive() {
		return errors.New("you are dead")
	}

	if voter.Role() != player.RoleMafiosi {
		return errors.New("only mafiosi can kill")
	}

	if g.IsDayPhase() {
		return errors.New("it's day now")
	}

	target, err := g.FindPlayer(candidate)
	if err != nil {
		return errors.New("invalid username")
	}

	g.AddKillVote(g.DayID(), voter, target)

	return nil
}

func (e *Engine) FindPlayerBySessionID(id uuid.UUID) (*player.Player, error) {
	s, err := e.findSession(id)
	if err != nil {
		return nil, err
	}

	g, err := e.FindGameByID(s.GameID)
	if err != nil {
		return nil, err
	}

	p, err := g.FindPlayer(s.Username)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (e *Engine) addSession(s *session.Session) {
	e.sessionsLocker.Lock()
	defer e.sessionsLocker.Unlock()

	e.sessions[s.ID] = s
}

func (e *Engine) findSession(id uuid.UUID) (*session.Session, error) {
	e.sessionsLocker.RLock()
	defer e.sessionsLocker.RUnlock()

	s, exists := e.sessions[id]
	if !exists {
		return nil, errors.New("session not found")
	}

	return s, nil
}

func (e *Engine) addGame(g *game.Game) {
	e.gamesLocker.Lock()
	defer e.gamesLocker.Unlock()

	e.games[g.ID()] = g
}

func (e *Engine) FindGameByID(id uuid.UUID) (*game.Game, error) {
	e.gamesLocker.RLock()
	defer e.gamesLocker.RUnlock()

	g, exists := e.games[id]
	if !exists {
		return nil, errors.New("game not found")
	}

	return g, nil
}
