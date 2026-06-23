package enums

type UserStatus int8

const (
	UserStatusActive UserStatus = iota + 1
	UserStatusBlocked
)
