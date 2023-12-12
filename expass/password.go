package expass

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/pbkdf2"
)

// GenerateHash generates a bcrypt hash of the password using the default cost
func GenerateHash(password []byte) ([]byte, error) {
	return bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
}

// CompareHash compares a bcrypt hashed password with its possible plaintext equivalent.
func CompareHash(hashedPassword, password []byte) bool {
	return bcrypt.CompareHashAndPassword(hashedPassword, password) == nil
}

func SaltPbkdf2Pass(sl, password string) string {
	pwd := []byte(password)
	salt := []byte(sl)
	iterations := 120000
	digest := sha256.New
	dk := pbkdf2.Key(pwd, salt, iterations, 32, digest)
	str := base64.StdEncoding.EncodeToString(dk)
	return "pbkdf2_sha256" + "$" + strconv.FormatInt(int64(iterations), 10) + "$" + string(salt) + "$" + str
}

func validCodeSize(code string) string {
	if len(code) < 16 {
		paddingLength := 16 - len(code)
		paddedString := code + strings.Repeat("6", paddingLength)
		return paddedString
	} else if len(code) < 24 {
		paddingLength := 24 - len(code)
		paddedString := code + strings.Repeat("4", paddingLength)
		return paddedString
	} else if len(code) < 32 {
		paddingLength := 32 - len(code)
		paddedString := code + strings.Repeat("2", paddingLength)
		return paddedString
	} else {
		return code[:32]
	}
}

func AesEncryptCBC(data, code string) (string, error) {
	code = validCodeSize(code)
	origData := []byte(data)
	key := []byte(code)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	origData = paddingPKCS7(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize])
	encrypted := make([]byte, len(origData))
	blockMode.CryptBlocks(encrypted, origData)
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func AesDecryptCBC(data, code string) (string, error) {
	code = validCodeSize(code)
	encrypted, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}
	key := []byte(code)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize])
	origData := make([]byte, len(encrypted))
	if len(origData)%blockMode.BlockSize() != 0 {
		return "", errors.New("crypto/cipher: input not full blocks")
	}
	blockMode.CryptBlocks(origData, encrypted)
	origData = unPaddingPKCS7(origData)
	return string(origData), nil
}

// paddingPKCS7 补码
func paddingPKCS7(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// unPaddingPKCS7 去码
func unPaddingPKCS7(origData []byte) []byte {
	if origData == nil {
		return nil
	}
	if len(origData) > 0 {
		length := len(origData)
		unpadding := int(origData[length-1])
		return origData[:(length - unpadding)]
	}
	return nil
}
