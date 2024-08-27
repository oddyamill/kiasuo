package client

type Token struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

// TODO decrypt token and check expiration date
type TokenPayload struct {
	Expiration int64 `json:"exp"`
}

type User struct {
	ID       int     `json:"id"`
	Username string  `json:"username"`
	Parent   bool    `json:"parent"`
	Children []Child `json:"children"`
}

type Child struct {
	ID          int    `json:"id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	MiddleName  string `json:"middle_name"`
	SchoolClass string `json:"school_class"`
	Age         int    `json:"age"`
}

type RawRecipient map[int]Recipients

type Recipients struct {
	Staff    map[string]map[string]int `json:"staff"`
	Students map[string]Student        `json:"students"`
}

type Student struct {
	Parents any  `json:"parents"`
	ID      *int `json:"id"`
}
