package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/lodthe/mafia/pkg/mafiapb"
)

type Engine struct {
	ctx context.Context

	client    *Client
	messenger *Messenger
}

func NewEngine(ctx context.Context, client *Client, messenger *Messenger) *Engine {
	return &Engine{
		ctx:       ctx,
		client:    client,
		messenger: messenger,
	}
}

func (e *Engine) Start() {
	go e.handleEvents()
	go e.handleUserInput()

	e.sendHelp(false)
	e.sendState()
}

func (e *Engine) handleEvents() {
	for {
		var event *mafiapb.GameEvent
		select {
		case <-e.ctx.Done():
			return
		case event = <-e.client.Events():
		}

		switch event.GetType() {
		case mafiapb.GameEvent_EVENT_PLAYER_JOINED:
			e.handlePlayerJoinedEvent(event.GetPayloadPlayerJoined())

		case mafiapb.GameEvent_EVENT_PLAYER_LEFT:
			e.handlePlayerLeftEvent(event.GetPayloadPlayerLeft())

		case mafiapb.GameEvent_EVENT_MESSAGE:
			e.handleMessageEvent(event.GetPayloadMessage())

		case mafiapb.GameEvent_EVENT_DAY_STARTED:
			e.handleDayStartedEvent(event.GetPayloadDayStarted())

		case mafiapb.GameEvent_EVENT_NIGHT_STARTED:
			e.handleNightStartedEvent(event.GetPayloadNightStarted())

		case mafiapb.GameEvent_EVENT_GAME_FINISHED:
			e.handleGameFinishedEvent(event.GetPayloadGameFinished())

		default:
			e.sendWithPrompt(fmt.Sprintf("Received an event with type %s", event.GetType()))
		}
	}
}

func (e *Engine) handlePlayerJoinedEvent(payload *mafiapb.GameEvent_PayloadPlayerJoined) {
	e.sendWithPrompt(fmt.Sprintf("%s joined the game", payload.GetPlayer().GetUsername()))
}

func (e *Engine) handlePlayerLeftEvent(payload *mafiapb.GameEvent_PayloadPlayerLeft) {
	e.sendWithPrompt(fmt.Sprintf("%s left the game", payload.GetPlayer().GetUsername()))
}

func (e *Engine) handleMessageEvent(payload *mafiapb.GameEvent_PayloadMessage) {
	e.sendWithPrompt(fmt.Sprintf("[%s]: %s", payload.GetSender().GetUsername(), payload.GetContent()))
}

func (e *Engine) handleDayStartedEvent(payload *mafiapb.GameEvent_PayloadDayStarted) {
	var text string
	if payload.GetKilledPlayer() != nil {
		text += fmt.Sprintf("%s was murdered that night\n\n", payload.GetKilledPlayer().GetUsername())
	}

	text += fmt.Sprintf("Day No. %d started", payload.GetDayId())

	e.sendWithPrompt(text)
}

func (e *Engine) handleNightStartedEvent(payload *mafiapb.GameEvent_PayloadNightStarted) {
	var text string
	if payload.GetKickedPlayer() != nil {
		text += fmt.Sprintf("The majority voted to kick %s today\n\n", payload.GetKickedPlayer().GetUsername())
	}

	text += fmt.Sprintf("Night No. %d started", payload.GetDayId())

	e.sendWithPrompt(text)
}

func (e *Engine) handleGameFinishedEvent(payload *mafiapb.GameEvent_PayloadGameFinished) {
	e.sendWithPrompt(fmt.Sprintf("Game finished! Winners: %s", payload.GetWinners().String()))
}

func (e *Engine) handleUserInput() {
	for {
		var input string
		select {
		case <-e.ctx.Done():
			return

		case input = <-e.messenger.Input():
		}

		switch {
		case strings.HasPrefix(input, "help"):
			e.sendHelp(true)

		case strings.HasPrefix(input, "state"):
			e.sendState()

		case strings.HasPrefix(input, "send"):
			e.handleSendCommand(input)

		case strings.HasPrefix(input, "votekick"):
			e.handleKickCommand(input)

		case strings.HasPrefix(input, "votekill"):
			e.handleKillCommand(input)

		default:
			e.sendError("Invalid command, please see help.")
		}
	}
}

func (e *Engine) sendHelp(withPrompt bool) {
	text := `===========

Command list:

> help - print this message.

> state - print the current state.

> send [msg] - send a text message.

> votekick [username] - vote for kick someone (available during the day).

> votekill [username] - vote for kill someone (available for mafiosi during the night).

===========`

	if withPrompt {
		e.sendWithPrompt(text)
	} else {
		e.messenger.Output() <- text
	}
}

func (e *Engine) sendState() {
	state := "state"

	e.sendWithPrompt(state)
}

func (e *Engine) handleSendCommand(input string) {
	content := input[len("send "):]
	if content == "" {
		e.sendError("The message cannot be empty")
		return
	}

	receivers, err := e.client.SendMessage(content)
	if err != nil {
		e.sendError(fmt.Sprintf("Failed to send message: %v", err))
		return
	}

	e.sendWithPrompt(fmt.Sprintf("The message has been sent to %d players", receivers))
}

func (e *Engine) handleKickCommand(input string) {
	username := input[len("votekick "):]

	err := e.client.VoteKick(username)
	if err != nil {
		e.sendError(err.Error())
		return
	}

	e.sendWithPrompt(fmt.Sprintf("You cast your vote for %s", username))
}

func (e *Engine) handleKillCommand(input string) {
	username := input[len("votekill "):]

	err := e.client.VoteKill(username)
	if err != nil {
		e.sendError(err.Error())
		return
	}

	e.sendWithPrompt(fmt.Sprintf("You cast your vote for %s", username))
}

func (e *Engine) sendError(text string) {
	text = strings.ReplaceAll(text, "rpc error: code = Unknown desc =", "")

	e.sendWithPrompt("[ERROR] " + text)
}

func (e *Engine) sendWithPrompt(text string) {
	msg := "\n\n" + text + "\n\n"

	msg += "Send me something: "

	e.send(msg)
}

func (e *Engine) send(text string) {
	e.messenger.Output() <- text
}
