package client

import (
	"testing"

	"github.com/kiasuo/bot/internal/crypto"
	"github.com/kiasuo/bot/internal/users"
)

func TestIsTokenExpired(t *testing.T) {
	client := New(&users.User{
		// This is a fake token that is expired
		AccessToken: crypto.Encrypt("{.eyJleHAiOjE3MjY3MzQwNzh9.}"),
	})

	if client.isTokenExpired() != true {
		t.Errorf("Client.isTokenExpired() = %t; want true", client.isTokenExpired())
	}

	// This is a fake token that is not expired
	client.User.AccessToken = crypto.Encrypt("{.eyJleHAiOjc5NjE3NjAwMDB9.}")

	if client.isTokenExpired() != false {
		t.Errorf("Client.isTokenExpired() = %t; want false", client.isTokenExpired())
	}
}
