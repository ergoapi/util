// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package expass

import (
	cryptorand "crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"github.com/ergoapi/util/common"

	"github.com/cockroachdb/errors"
)

var (
	AlphaNum        = common.DIGITS + common.Alpha
	AlphaNumSymbols = AlphaNum + common.Symbols
)

var CHARS = fmt.Sprintf("%s%s%s%s", common.DIGITS, common.LETTERS, strings.ToUpper(common.LETTERS), common.PUNC)

// NewPwGen generates a cryptographically secure random string with uniform distribution
func NewPwGen(length int, chars string) (string, error) {
	if length <= 0 {
		return "", errors.New("length must be positive")
	}
	if len(chars) == 0 {
		return "", errors.New("chars cannot be empty")
	}

	result := make([]byte, length)
	charLen := big.NewInt(int64(len(chars)))

	for i := 0; i < length; i++ {
		// Use crypto/rand.Int for uniform distribution
		n, err := cryptorand.Int(cryptorand.Reader, charLen)
		if err != nil {
			return "", errors.Wrap(err, "failed to generate random number")
		}
		result[i] = chars[n.Int64()]
	}
	return string(result), nil
}

// PwGenNum generates a random string of the given length out of numeric characters
func PwGenNum(length int) (string, error) {
	return NewPwGen(length, common.DIGITS)
}

// PwGenAlpha generates a random string of the given length out of alphabetic characters
func PwGenAlpha(length int) (string, error) {
	return NewPwGen(length, common.Alpha)
}

// PwGenSymbols generates a random string of the given length out of symbols
func PwGenSymbols(length int) (string, error) {
	return NewPwGen(length, common.Symbols)
}

// PwGenAlphaNum generates a random string of the given length out of alphanumeric characters
func PwGenAlphaNum(length int) (string, error) {
	return NewPwGen(length, AlphaNum)
}

// PwGenAlphaNumSymbols generates a random string of the given length out of alphanumeric characters and
// symbols
func PwGenAlphaNumSymbols(length int) (string, error) {
	return NewPwGen(length, AlphaNumSymbols)
}
