package commands

import (
	"github.com/kiasuo/bot/internal/client"
	"github.com/kiasuo/bot/internal/users"
)

type Context struct {
	Command   string
	Arguments string
	User      users.User
}

func NewContext(command, arguments string, user users.User) Context {
	return Context{command, arguments, user}
}

func (c *Context) GetClient() *client.Client {
	return client.New(&c.User)
}
