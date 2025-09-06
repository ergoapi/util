// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package expass

import (
	cryptorand "crypto/rand"
	"math/big"

	"github.com/cockroachdb/errors"
	"github.com/ergoapi/util/common"
)

// Character sets for password generation
var (
	AlphaNum        = common.DIGITS + common.Alpha
	AlphaNumSymbols = AlphaNum + common.Symbols
)

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

// PasswordStrength represents the strength level of a password
type PasswordStrength int

const (
	StrengthVeryWeak PasswordStrength = iota
	StrengthWeak
	StrengthFair
	StrengthStrong
	StrengthVeryStrong
)

// String returns the string representation of password strength
func (ps PasswordStrength) String() string {
	switch ps {
	case StrengthVeryWeak:
		return "Very Weak"
	case StrengthWeak:
		return "Weak"
	case StrengthFair:
		return "Fair"
	case StrengthStrong:
		return "Strong"
	case StrengthVeryStrong:
		return "Very Strong"
	default:
		return "Unknown"
	}
}

// CheckPasswordStrength analyzes password strength based on length and character variety
func CheckPasswordStrength(password string) PasswordStrength {
	length := len(password)
	if length == 0 {
		return StrengthVeryWeak
	}

	var (
		hasLower   bool
		hasUpper   bool
		hasDigit   bool
		hasSpecial bool
	)

	for _, ch := range password {
		switch {
		case ch >= 'a' && ch <= 'z':
			hasLower = true
		case ch >= 'A' && ch <= 'Z':
			hasUpper = true
		case ch >= '0' && ch <= '9':
			hasDigit = true
		default:
			hasSpecial = true
		}
	}

	// Count character types
	charTypes := 0
	if hasLower {
		charTypes++
	}
	if hasUpper {
		charTypes++
	}
	if hasDigit {
		charTypes++
	}
	if hasSpecial {
		charTypes++
	}

	// Evaluate strength based on length and character variety
	switch {
	case length < 6:
		return StrengthVeryWeak
	case length < 8:
		if charTypes >= 3 {
			return StrengthWeak
		}
		return StrengthVeryWeak
	case length < 12:
		if charTypes >= 4 {
			return StrengthFair
		} else if charTypes >= 3 {
			return StrengthWeak
		}
		return StrengthVeryWeak
	case length < 16:
		if charTypes >= 4 {
			return StrengthStrong
		} else if charTypes >= 3 {
			return StrengthFair
		}
		return StrengthWeak
	default:
		if charTypes >= 4 {
			return StrengthVeryStrong
		} else if charTypes >= 3 {
			return StrengthStrong
		}
		return StrengthFair
	}
}

// GenerateSecurePassword generates a cryptographically secure password with specified requirements
func GenerateSecurePassword(length int, requireUpper, requireLower, requireDigit, requireSpecial bool) (string, error) {
	if length < 4 && (requireUpper && requireLower && requireDigit && requireSpecial) {
		return "", errors.New("length too short for all requirements")
	}

	// Build character set based on requirements
	var chars string
	var required []string

	if requireUpper {
		upperChars := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		chars += upperChars
		required = append(required, upperChars)
	}
	if requireLower {
		lowerChars := "abcdefghijklmnopqrstuvwxyz"
		chars += lowerChars
		required = append(required, lowerChars)
	}
	if requireDigit {
		digitChars := "0123456789"
		chars += digitChars
		required = append(required, digitChars)
	}
	if requireSpecial {
		specialChars := "!@#$%^&*()_+-=[]{}|;:,.<>?"
		chars += specialChars
		required = append(required, specialChars)
	}

	// Default to alphanumeric if no requirements
	if chars == "" {
		chars = AlphaNum
	}

	// Generate password ensuring all required character types
	for {
		password, err := NewPwGen(length, chars)
		if err != nil {
			return "", err
		}

		// Check if all requirements are met
		allMet := true
		for _, reqChars := range required {
			found := false
			for _, ch := range password {
				if containsRune(reqChars, ch) {
					found = true
					break
				}
			}
			if !found {
				allMet = false
				break
			}
		}

		if allMet || len(required) == 0 {
			return password, nil
		}
		// Retry if requirements not met (rare for reasonable length passwords)
	}
}

// containsRune checks if a string contains a specific rune
func containsRune(s string, r rune) bool {
	for _, ch := range s {
		if ch == r {
			return true
		}
	}
	return false
}
