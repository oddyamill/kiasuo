package users

import (
	"testing"
)

func TestIsReady(t *testing.T) {
	user := User{
		State: Ready,
	}

	if user.IsReady() != true {
		t.Errorf("User.IsReady() = %t; want true", user.IsReady())
	}

	user.State = Pending

	if user.IsReady() != false {
		t.Errorf("User.IsReady() = %t; want false", user.IsReady())
	}
}

func TestIsTokenExpired(t *testing.T) {
	user := User{
		// This is a fake token that is expired
		AccessToken: "{.eyJleHAiOjE3MjY3MzQwNzh9.}",
	}

	if user.IsTokenExpired() != true {
		t.Errorf("User.IsTokenExpired() = %t; want true", user.IsTokenExpired())
	}

	// This is a fake token that is not expired
	user.AccessToken = "{.eyJleHAiOjc5NjE3NjAwMDB9.}"

	if user.IsTokenExpired() != false {
		t.Errorf("User.IsTokenExpired() = %t; want false", user.IsTokenExpired())
	}
}
