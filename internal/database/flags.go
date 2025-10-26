package database

import "encoding/json"

type UserFlag int8

const (
	UserFlagCache UserFlag = 1 << iota
	UserFlagShowPasses
	UserFlagShowEmptyLessons
)

func (uf UserFlag) MarshalBinary() ([]byte, error) {
	return json.Marshal(uf)
}

func (uf UserFlag) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &uf)
}
