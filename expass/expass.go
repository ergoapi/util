package expass

import (
	cryptorand "crypto/rand"
	"fmt"
	"strings"

	"github.com/ergoapi/util/common"
)

var (
	AlphaNum        = common.DIGITS + common.Alpha
	AlphaNumSymbols = AlphaNum + common.Symbols
)

var CHARS = fmt.Sprintf("%s%s%s%s", common.DIGITS, common.LETTERS, strings.ToUpper(common.LETTERS), common.PUNC)

func NewPwGen(length int, chars string) string {
	var bytes = make([]byte, length)
	var op = byte(len(chars))

	cryptorand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = chars[b%op]
	}
	return string(bytes)
}

// PwGenNum generates a random string of the given length out of numeric characters
func PwGenNum(length int) string {
	return NewPwGen(length, common.DIGITS)
}

// PwGenAlpha generates a random string of the given length out of alphabetic characters
func PwGenAlpha(length int) string {
	return NewPwGen(length, common.Alpha)
}

// PwGenSymbols generates a random string of the given length out of symbols
func PwGenSymbols(length int) string {
	return NewPwGen(length, common.Symbols)
}

// PwGenAlphaNum generates a random string of the given length out of alphanumeric characters
func PwGenAlphaNum(length int) string {
	return NewPwGen(length, AlphaNum)
}

// PwGenAlphaNumSymbols generates a random string of the given length out of alphanumeric characters and
// symbols
func PwGenAlphaNumSymbols(length int) string {
	return NewPwGen(length, AlphaNumSymbols)
}
