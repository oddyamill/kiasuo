package commands

import (
	"github.com/kiasuo/bot/internal/client"
	"github.com/kiasuo/bot/internal/users"
)

type Context struct {
	Command string
	User    users.User
}

func (c *Context) GetClient() *client.Client {
	return &client.Client{User: &c.User}
}
