package expass

import (
	"bytes"
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

// deriveKey derives a key from password using PBKDF2
func deriveKey(password string, salt []byte, keySize int) []byte {
	if keySize <= 0 {
		keySize = 32 // Default to AES-256
	}
	return pbkdf2.Key([]byte(password), salt, DefaultConfig().PBKDF2Iterations, keySize, sha256.New)
}

// AesEncryptCBC encrypts data using AES-CBC with a random IV
func AesEncryptCBC(data, password string) (string, error) {
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
	key := deriveKey(password, salt, 32)
	
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errors.Wrap(err, "failed to create cipher")
	}
	
	// Generate random IV
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(cryptorand.Reader, iv); err != nil {
		return "", errors.Wrap(err, "failed to generate IV")
	}
	
	// Pad plaintext
	origData := []byte(data)
	origData = paddingPKCS7(origData, block.BlockSize())
	
	// Encrypt
	blockMode := cipher.NewCBCEncrypter(block, iv)
	encrypted := make([]byte, len(origData))
	blockMode.CryptBlocks(encrypted, origData)
	
	// Combine salt + iv + encrypted and encode
	result := append(salt, iv...)
	result = append(result, encrypted...)
	
	return base64.StdEncoding.EncodeToString(result), nil
}

// AesDecryptCBC decrypts data encrypted with AesEncryptCBC
func AesDecryptCBC(data, password string) (string, error) {
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
	
	// Extract salt, IV, and ciphertext
	if len(encrypted) < 48 { // 32 bytes salt + 16 bytes IV
		return "", errors.New("invalid encrypted data format")
	}
	
	salt := encrypted[:32]
	iv := encrypted[32:48]
	ciphertext := encrypted[48:]
	
	// Derive key from password
	key := deriveKey(password, salt, 32)
	
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errors.Wrap(err, "failed to create cipher")
	}
	
	if len(ciphertext)%block.BlockSize() != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}
	
	// Decrypt
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(ciphertext))
	blockMode.CryptBlocks(origData, ciphertext)
	
	// Remove padding
	origData, err = unPaddingPKCS7(origData)
	if err != nil {
		return "", errors.Wrap(err, "failed to remove padding")
	}
	
	return string(origData), nil
}

// paddingPKCS7 补码
func paddingPKCS7(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// unPaddingPKCS7 removes PKCS7 padding
func unPaddingPKCS7(origData []byte) ([]byte, error) {
	if origData == nil || len(origData) == 0 {
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
