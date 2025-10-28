package commands

import (
	"context"

	"github.com/kiasuo/bot/internal/client"
	"github.com/kiasuo/bot/internal/database"
)

type Context struct {
	Command   string
	Arguments string
	User      database.User
}

func NewContext(command, arguments string, user database.User) Context {
	return Context{command, arguments, user}
}

func (c *Context) Context() context.Context {
	return context.Background()
}

func (c *Context) GetClient() *client.Client {
	return client.New(&c.User)
}
