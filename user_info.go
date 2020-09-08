package irc

import "fmt"

// Each user is distinguished from other users by a unique nickname.
type UserInfo struct {
	Nickname string
	Operator bool
}

// Creates a new user.
func NewUserInfo(Nickname string, Operator bool) (*UserInfo, error) {
	var err error
	if !isValidNickname(Nickname) {
		err = fmt.Errorf("invalid nickname: \"%s\"", Nickname)
	}
	return &UserInfo{
		Nickname: Nickname,
		Operator: Operator,
	}, err
}

func isValidNickname(nickname string) bool {
	return nickNameRegexp.MatchString(nickname)
}

type UserInfoSlice []UserInfo

func (p UserInfoSlice) Len() int           { return len(p) }
func (p UserInfoSlice) Less(i, j int) bool { return p[i].Nickname < p[j].Nickname }
func (p UserInfoSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
