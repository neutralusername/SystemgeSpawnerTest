package app

import (
	"SystemgeSpawnerTest/topics"

	"github.com/neutralusername/Systemge/Message"
	"github.com/neutralusername/Systemge/Node"
)

func (app *App) GetSyncMessageHandlers() map[string]Node.SyncMessageHandler {
	return map[string]Node.SyncMessageHandler{
		topics.SYNC: func(node *Node.Node, message *Message.Message) (string, error) {
			return "", nil
		},
	}
}

func (app *App) GetAsyncMessageHandlers() map[string]Node.AsyncMessageHandler {
	return map[string]Node.AsyncMessageHandler{
		topics.ASYNC: func(node *Node.Node, message *Message.Message) error {
			return nil
		},
	}
}
