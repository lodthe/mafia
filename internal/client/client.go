package client

import (
	"context"
	"log"

	"github.com/lodthe/mafia/pkg/mafiapb"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	ctx    context.Context
	cli    mafiapb.MafiaClient
	stream mafiapb.Mafia_JoinGameClient
	events chan *mafiapb.GameEvent
}

func NewClient(ctx context.Context, username string, conn *grpc.ClientConn) (*Client, error) {
	cli := mafiapb.NewMafiaClient(conn)

	stream, err := cli.JoinGame(ctx, &mafiapb.JoinGameRequest{Username: username})
	if err != nil {
		return nil, errors.Wrap(err, "failed to join a game")
	}

	md, err := stream.Header()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get gRPC headers")
	}

	sessionID, err := mafiapb.FetchSessionID(md)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get session ID")
	}

	ctx = metadata.NewOutgoingContext(ctx, mafiapb.WithSessionID(sessionID))

	return &Client{
		ctx:    ctx,
		cli:    cli,
		stream: stream,
		events: make(chan *mafiapb.GameEvent),
	}, nil
}

func (c *Client) Events() <-chan *mafiapb.GameEvent {
	return c.events
}

func (c *Client) ForwardEvents() {
	for {
		event, err := c.stream.Recv()
		if err != nil {
			log.Fatalf("\n\nServer closed\n")
			return
		}

		c.events <- event
	}
}

func (c *Client) SendMessage(content string) (receivers uint, err error) {
	resp, err := c.cli.SendMessage(c.ctx, &mafiapb.SendMessageRequest{Content: content})
	if err != nil {
		return 0, err
	}

	return uint(resp.GetReceiverCount()), nil
}
