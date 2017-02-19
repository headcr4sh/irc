package irc

import "strings"

// stringLowercaser can be used to convert a string used in the IRC protocol to it's lowercase representation.
// Background: Because of IRC's Scandinavian origin, the characters {}|^ are considered to be the lower case
// equivalents of the characters []\~, respectively.
var stringLowercaser = strings.NewReplacer("{", "[", "}", "]", "|", "\\", "^", "~")

// stringLowercaser can be used to convert a string used in the IRC protocol to it's uppercase representation.
// Background: Because of IRC's Scandinavian origin, the characters {}|^ are considered to be the lower case
// equivalents of the characters []\~, respectively.
var stringUppercaser = strings.NewReplacer("[", "{", "]", "}", "\\", "|", "~", "^")

func toLowercase(str string) string {
	return strings.ToLower(stringLowercaser.Replace(str))
}

func toUppercase(str string) string {
	return stringUppercaser.Replace(strings.ToUpper(str))
}
