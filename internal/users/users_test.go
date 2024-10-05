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
