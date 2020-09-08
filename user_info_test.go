package irc

import (
	"testing"
)

func TestNewUserInfo(t *testing.T) {
	var testdata = []struct {
		nickname string
		operator bool
		valid    bool
	}{
		{
			"username",
			false,
			true,
		},
		{
			"John Doe",
			false,
			false,
		},
	}
	for _, tt := range testdata {
		_, err := NewUserInfo(tt.nickname, tt.operator)
		if err != nil && tt.valid {
			t.Errorf("Nickname \"%s\" with Operator=%v should produce a valid UserInfo struct.", tt.nickname, tt.operator)
		}
		if err == nil && !tt.valid {
			t.Errorf("Nickname \"%s\" with Operator=%v should *not* produce a valid UserInfo struct.", tt.nickname, tt.operator)
		}
	}
}
