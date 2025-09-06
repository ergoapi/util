// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package expass

import (
	"crypto/aes"
	"crypto/cipher"
	cryptorand "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"strconv"

	"github.com/cockroachdb/errors"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
)

// Config holds configuration for password operations
type Config struct {
	PBKDF2Iterations int
	SaltLength       int
	BCryptCost       int
}

// DefaultConfig returns default configuration
func DefaultConfig() *Config {
	return &Config{
		PBKDF2Iterations: 600000, // OWASP 2024 recommendation
		SaltLength:       32,
		BCryptCost:       bcrypt.DefaultCost,
	}
}

// GenerateHash generates a bcrypt hash of the password using the default cost
func GenerateHash(password []byte) ([]byte, error) {
	return GenerateHashWithCost(password, bcrypt.DefaultCost)
}

// GenerateHashWithCost generates a bcrypt hash with specified cost
func GenerateHashWithCost(password []byte, cost int) ([]byte, error) {
	if len(password) == 0 {
		return nil, errors.New("password cannot be empty")
	}
	return bcrypt.GenerateFromPassword(password, cost)
}

// CompareHash compares a bcrypt hashed password with its possible plaintext equivalent.
func CompareHash(hashedPassword, password []byte) bool {
	return bcrypt.CompareHashAndPassword(hashedPassword, password) == nil
}

// SaltPbkdf2Pass generates a PBKDF2 hash with configurable iterations
func SaltPbkdf2Pass(salt, password string, iterations int) string {
	if iterations <= 0 {
		iterations = DefaultConfig().PBKDF2Iterations
	}
	pwd := []byte(password)
	saltBytes := []byte(salt)
	digest := sha256.New
	dk := pbkdf2.Key(pwd, saltBytes, iterations, 32, digest)
	str := base64.StdEncoding.EncodeToString(dk)
	return "pbkdf2_sha256" + "$" + strconv.FormatInt(int64(iterations), 10) + "$" + salt + "$" + str
}

// SaltPbkdf2PassDefault generates a PBKDF2 hash with default iterations
func SaltPbkdf2PassDefault(salt, password string) string {
	return SaltPbkdf2Pass(salt, password, DefaultConfig().PBKDF2Iterations)
}

// deriveKey derives a 32-byte key from password using PBKDF2 (for AES-256)
func deriveKey(password string, salt []byte) []byte {
	return pbkdf2.Key([]byte(password), salt, DefaultConfig().PBKDF2Iterations, 32, sha256.New)
}

// AesEncryptGCM encrypts data using AES-GCM for authenticated encryption
// GCM provides both confidentiality and authenticity in a single operation.
func AesEncryptGCM(data, password string) (string, error) {
	if data == "" {
		return "", errors.New("data cannot be empty")
	}
	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	// Generate random salt for key derivation
	salt := make([]byte, 32)
	if _, err := io.ReadFull(cryptorand.Reader, salt); err != nil {
		return "", errors.Wrap(err, "failed to generate salt")
	}

	// Derive key from password
	key := deriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errors.Wrap(err, "failed to create cipher")
	}

	// Create GCM mode
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.Wrap(err, "failed to create GCM")
	}

	// Generate random nonce
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(cryptorand.Reader, nonce); err != nil {
		return "", errors.Wrap(err, "failed to generate nonce")
	}

	// Encrypt and authenticate
	encrypted := aesgcm.Seal(nil, nonce, []byte(data), salt)

	// Combine salt + nonce + encrypted (which includes auth tag)
	result := append(salt, nonce...)
	result = append(result, encrypted...)

	return base64.StdEncoding.EncodeToString(result), nil
}

// AesDecryptGCM decrypts data encrypted by AesEncryptGCM using AES-GCM.
func AesDecryptGCM(data, password string) (string, error) {
	if data == "" {
		return "", errors.New("encrypted data cannot be empty")
	}
	if password == "" {
		return "", errors.New("password cannot be empty")
	}

	// Decode from base64
	encrypted, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", errors.Wrap(err, "failed to decode base64")
	}

	// Extract salt
	if len(encrypted) < 32 {
		return "", errors.New("invalid encrypted data format")
	}
	salt := encrypted[:32]

	// Derive key from password
	key := deriveKey(password, salt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errors.Wrap(err, "failed to create cipher")
	}

	// Create GCM mode
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.Wrap(err, "failed to create GCM")
	}

	// Extract nonce and ciphertext
	if len(encrypted) < 32+aesgcm.NonceSize()+aesgcm.Overhead() {
		return "", errors.New("invalid encrypted data format")
	}

	nonce := encrypted[32 : 32+aesgcm.NonceSize()]
	ciphertext := encrypted[32+aesgcm.NonceSize():]

	// Decrypt and verify (using salt as additional authenticated data)
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, salt)
	if err != nil {
		return "", errors.Wrap(err, "decryption failed")
	}

	return string(plaintext), nil
}

// Deprecated: Use AesEncryptGCM instead. The implementation was switched to
// AES-GCM; this wrapper preserves backward compatibility with the old name.
func AesEncryptCBC(data, password string) (string, error) {
	return AesEncryptGCM(data, password)
}

// Deprecated: Use AesDecryptGCM instead. The implementation was switched to
// AES-GCM; this wrapper preserves backward compatibility with the old name.
func AesDecryptCBC(data, password string) (string, error) {
	return AesDecryptGCM(data, password)
}

// unPaddingPKCS7 removes PKCS7 padding
// Note: This function is kept for testing purposes only
func unPaddingPKCS7(origData []byte) ([]byte, error) {
	if len(origData) == 0 {
		return nil, errors.New("invalid data: empty input")
	}

	length := len(origData)
	unpadding := int(origData[length-1])

	// Validate padding
	if unpadding > length || unpadding == 0 {
		return nil, errors.New("invalid PKCS7 padding")
	}

	// Check all padding bytes are the same
	for i := length - unpadding; i < length; i++ {
		if origData[i] != byte(unpadding) {
			return nil, errors.New("invalid PKCS7 padding")
		}
	}

	return origData[:(length - unpadding)], nil
}
