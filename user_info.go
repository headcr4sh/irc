package irc

// Each user is distinguished from other users by a unique nickname.
type UserInfo struct {
	Nickname string
	Operator bool
}

// Creates a new user.
func NewUserInfo(Nickname string, Operator bool) *UserInfo {
	return &UserInfo{
		Nickname: Nickname,
		Operator: Operator,
	}
}

func isValidNickname(nick string) bool {
	return nickNameRegexp.MatchString(nick)
}
