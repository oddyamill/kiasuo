package database

import "encoding/json"

type UserState int8

const (
	UserStateUnknown UserState = iota
	UserStateReady
	UserStatePending
	UserStateBlacklisted
)

func (u UserState) String() string {
	switch u {
	case UserStateUnknown:
		return "неизвестно"
	case UserStateReady:
		return "готов"
	case UserStatePending:
		return "ожидает"
	case UserStateBlacklisted:
		return "заблокирован"
	default:
		panic("bad user state")
	}
}

func (u UserState) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}

func (u UserState) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &u)
}
