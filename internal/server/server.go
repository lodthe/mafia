package server

import (
	"context"
	"log"

	"github.com/lodthe/mafia/internal/server/player"
	"github.com/lodthe/mafia/pkg/mafiapb"
	"github.com/pkg/errors"
	"google.golang.org/grpc/metadata"
)

const EventsBufferSize = 8

type Server struct {
	mafiapb.UnimplementedMafiaServer

	ctx    context.Context
	engine *Engine
}

func NewServer(ctx context.Context, engine *Engine) *Server {
	return &Server{
		ctx:    ctx,
		engine: engine,
	}
}

func (s *Server) JoinGame(in *mafiapb.JoinGameRequest, stream mafiapb.Mafia_JoinGameServer) error {
	if in.GetUsername() == "" {
		return errors.New("empty username")
	}

	events := make(chan *mafiapb.GameEvent, EventsBufferSize)
	sess, err := s.engine.AddPlayer(in.GetUsername(), events)
	if err != nil {
		return err
	}

	err = stream.SendHeader(mafiapb.WithSessionID(sess.ID))
	if err != nil {
		log.Printf("failed to send stream header: %v\n", err)
		return errors.New("failed to send metadata header")
	}

	for {
		var event *mafiapb.GameEvent
		select {
		case <-s.ctx.Done():
			return nil

		case event = <-events:
		}

		err := stream.Send(event)
		if err != nil {
			log.Printf("failed to send an event: %v\n%v\n", err, event)

			s.engine.RemovePlayer(sess)

			break
		}
	}

	return nil
}

func (s *Server) fetchPlayer(ctx context.Context) (*player.Player, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("missed metadata")
	}

	sessionID, err := mafiapb.FetchSessionID(md)
	if err != nil {
		return nil, err
	}

	p, err := s.engine.FindPlayerBySessionID(sessionID)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (s *Server) SendMessage(ctx context.Context, in *mafiapb.SendMessageRequest) (*mafiapb.SendMessageResponse, error) {
	p, err := s.fetchPlayer(ctx)
	if err != nil {
		return nil, err
	}

	if in.GetContent() == "" {
		return nil, errors.New("missed content")
	}

	receivers, err := s.engine.SendMessage(p, in.GetContent())
	if err != nil {
		return nil, err
	}

	return &mafiapb.SendMessageResponse{
		ReceiverCount: uint64(len(receivers)),
	}, nil
}

func (s *Server) DayVote(ctx context.Context, in *mafiapb.DayVoteRequest) (*mafiapb.DayVoteResponse, error) {
	p, err := s.fetchPlayer(ctx)
	if err != nil {
		return nil, err
	}

	if in.GetUsername() == "" {
		return nil, errors.New("missed username")
	}

	err = s.engine.VoteKick(p, in.GetUsername())
	if err != nil {
		return nil, err
	}

	return &mafiapb.DayVoteResponse{}, nil
}

func (s *Server) NightVote(ctx context.Context, in *mafiapb.NightVoteRequest) (*mafiapb.NightVoteResponse, error) {
	p, err := s.fetchPlayer(ctx)
	if err != nil {
		return nil, err
	}

	if in.GetUsername() == "" {
		return nil, errors.New("missed username")
	}

	err = s.engine.KillVote(p, in.GetUsername())
	if err != nil {
		return nil, err
	}

	return &mafiapb.NightVoteResponse{}, nil
}
