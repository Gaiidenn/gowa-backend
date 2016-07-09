package rpcWebsocket

import (
	"github.com/Gaiidenn/gowa-backend/chats"
	"github.com/Gaiidenn/gowa-backend/users"
	"log"
)

// UserRPCService for jsonRPC requests
type ChatRPCService struct {
}

// Save the user in database
func (cs *ChatRPCService) NewChat(users []*users.User, reply *chats.Chat) error {
	chat, err := chats.NewChat(users)
	if err != nil {
		return err
	}

	*reply = *chat
	return nil
}

func (cs *ChatRPCService) OpenChat(key *string, reply *chats.Chat) error {
	chat, err := chats.GetByKey(*key)
	if err != nil {
		return err
	}
	*reply = *chat
	return nil
}

func (cs *ChatRPCService) NewMessage(m *chats.Message, r *bool) error {
	log.Println("RPC NewMessage : ", m)
	chat, err := chats.NewMessage(m)
	if err != nil {
		return err
	}

	for _, u := range chat.Users {
		key := u.Token
		if c, ok := h.connections[key]; ok {
			var rr *bool
			call := RpcCall{
				Method: "ChatService.msgReceived",
				Args: m,
				Reply: rr,
			}
			c.call <- &call
		}
	}

	*r = true
	return nil
}