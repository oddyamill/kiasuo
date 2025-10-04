package users

type UserState int8

const (
	Unknown UserState = iota
	Ready
	Pending
	Blacklisted
)

func (u UserState) String() string {
	switch u {
	case Unknown:
		return "неизвестно"
	case Ready:
		return "готов"
	case Pending:
		return "ожидает"
	case Blacklisted:
		return "заблокирован"
	default:
		panic("bad user state")
	}
}
