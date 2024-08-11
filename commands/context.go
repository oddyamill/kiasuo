package commands

import (
	"github.com/kiasuo/bot/client"
	"github.com/kiasuo/bot/users"
)

type Context struct {
	Command string
	User    users.User
}

func (c Context) GetClient() client.Client {
	return client.Client{User: c.User}
}
