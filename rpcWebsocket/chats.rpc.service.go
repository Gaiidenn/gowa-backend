package rpcWebsocket

import (
	"github.com/Gaiidenn/gowa-backend/chats"
)

// UserRPCService for jsonRPC requests
type ChatRPCService struct {
}

// Save the user in database
func (cs *ChatRPCService) OpenChat(key string, reply *chats.Chat) error {


	//*reply = *chat
	return nil
}

