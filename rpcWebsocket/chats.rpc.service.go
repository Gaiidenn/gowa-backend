package rpcWebsocket

import (
	"github.com/Gaiidenn/gowa-backend/chats"
	"github.com/Gaiidenn/gowa-backend/users"
	"errors"
)

// UserRPCService for jsonRPC requests
type ChatRPCService struct {
}


func (cs *ChatRPCService) OpenPrivateChat(users []*users.User, reply *chats.Chat) error {
	if len(users) < 2 {
		return errors.New("Not enough parameters")
	}
	if len(users) > 2 {
		return errors.New("Too many parameters")
	}
	chat, err := chats.OpenPrivateChat(users[0], users[1])
	if err != nil {
		return err
	}

	for _, user := range users {
		pm, err := user.GetPeopleMet()
		if err != nil {
			return err
		}
		var r *bool
		call := RpcCall{
			Method: "UsersService.updatePeopleMet",
			Args: []interface{}{user.ID,pm},
			Reply: r,
		}
		h.broadcast <- &call
	}

	*reply = *chat
	return nil
}

func (cs *ChatRPCService) GetChat(chatID string, reply *chats.Chat) error {
	if len(chatID) == 0 {
		return errors.New("No such ID")
	}
	chat, err := chats.GetByID(chatID)
	if err != nil {
		return err
	}

	*reply = *chat
	return nil
}

func (cs *ChatRPCService) NewMessage(m *chats.Message, r *bool) error {
	//log.Println("RPC NewMessage : ", m)

	err := m.Save()
	if err != nil {
		return err
	}
	//log.Println("Message saved")
	chat, err := chats.GetByID(m.ChatID)
	if err != nil {
		return err
	}
	//log.Println("Chat getted")
	for _, c := range h.connections {
		for _, u := range chat.Users {
			if c.user.Username == u.Username {
				var rr *bool
				call := RpcCall{
					Method: "ChatService.msgReceived",
					Args: m,
					Reply: rr,
				}
				c.call <- &call
			}
		}
	}

	*r = true
	return nil
}

// -------------------- OLD CODE -----------------

/*
// Save the user in database
func (cs *ChatRPCService) NewChat(users []*users.User, reply *chats.Chat) error {
	chat, err := chats.NewChat(users)
	if err != nil {
		return err
	}

	for _, u := range chat.Users {
		key := u.Token
		if c, ok := h.connections[key]; ok {
			var rr *bool
			call := RpcCall{
				Method: "ChatService.registerChat",
				Args: chat,
				Reply: rr,
			}
			c.call <- &call
		}
	}

	*reply = *chat
	return nil
}

func (cs *ChatRPCService) OpenChat(key *string, reply *chats.Chat) error {
	chat, err := chats.GetByKey(*key)
	if err != nil {
		return err
	}

	for _, u := range chat.Users {
		key := u.Token
		if c, ok := h.connections[key]; ok {
			var rr *bool
			call := RpcCall{
				Method: "ChatService.registerChat",
				Args: chat,
				Reply: rr,
			}
			c.call <- &call
		}
	}

	*reply = *chat
	return nil
}

func (cs *ChatRPCService) NewMessage(m *chats.Message, r *bool) error {
	//log.Println("RPC NewMessage : ", m)
	chat, err := chats.NewMessage(m)
	if err != nil {
		return err
	}

	for _, c := range h.connections {
		for _, u := range chat.Users {
			if c.user.Username == u.Username {
				var rr *bool
				call := RpcCall{
					Method: "ChatService.msgReceived",
					Args: m,
					Reply: rr,
				}
				c.call <- &call
			}
		}
	}

	*r = true
	return nil
}*/