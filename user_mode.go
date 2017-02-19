package irc

import "fmt"

type UserMode rune

const (
	UserModeAway                    UserMode = 'a'
	UserModeInvisible               UserMode = 'i'
	UserModeReceivesWallops         UserMode = 'w'
	UserModeRestrictedConnection    UserMode = 'r'
	UserModeOperator                UserMode = 'o'
	UserModeLocalOperator           UserMode = 'O'
	UserModeReceiptForServerNotices UserMode = 's'
)

// Numeric returns the numeric representation of the user mode.
// The only use of this numeric representation is during emission of the
// USER message, whereas according to the specification only modes 'i' and 'w'
// have any meaning.
func (um UserMode) Numeric() int {
	switch um {
	case UserModeReceivesWallops:
		return 4
	case UserModeInvisible:
		return 8
	default:
		panic(fmt.Errorf("no numeric/bitmask representation for UserMode %v", um))
	}
}

func UserModesFromBitmask(mask int) UserModes {
	var um UserModes
	if mask&4 == 4 {
		um = append(um, UserModeReceivesWallops)
	}
	if mask&8 == 8 {
		um = append(um, UserModeInvisible)
	}
	return um
}

type UserModes []UserMode

func (modes UserModes) Bitmask() int {
	m := 0
	for _, um := range modes {
		m |= um.Numeric()
	}
	return m
}
