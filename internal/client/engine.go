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
		case mafiapb.GameEvent_EVENT_MESSAGE:
			e.handleMessageEvent(event.GetPayloadMessage())

		default:
			e.sendWithPrompt(fmt.Sprintf("Received an event with type %s", event.GetType()))
		}
	}
}

func (e *Engine) handleMessageEvent(payload *mafiapb.GameEvent_PayloadMessage) {
	e.sendWithPrompt(fmt.Sprintf("[%s]: %s", payload.GetSender().GetUsername(), payload.GetContent()))
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

func (e *Engine) sendError(text string) {
	e.sendWithPrompt("[ERROR] " + text)
}

func (e *Engine) sendWithPrompt(text string) {
	msg := "\n\n" + text + "\n\n"

	msg += "Send me something: "

	e.messenger.Output() <- msg
}
