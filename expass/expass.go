package expass

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	cryptorand "crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/ergoapi/util/common"

	"golang.org/x/crypto/pbkdf2"
)

var (
	AlphaNum        = common.DIGITS + common.Alpha
	AlphaNumSymbols = AlphaNum + common.Symbols
)

var CHARS = fmt.Sprintf("%s%s%s%s", common.DIGITS, common.LETTERS, strings.ToUpper(common.LETTERS), common.PUNC)

func SaltPbkdf2Pass(sl, password string) string {
	pwd := []byte(password)
	salt := []byte(sl)
	iterations := 120000
	digest := sha256.New
	dk := pbkdf2.Key(pwd, salt, iterations, 32, digest)
	str := base64.StdEncoding.EncodeToString(dk)
	return "pbkdf2_sha256" + "$" + strconv.FormatInt(int64(iterations), 10) + "$" + string(salt) + "$" + str
}

func Encrypt(code, p string) string {
	if len(code) == 16 {
		// 转成字节数组
		origData := []byte(p)
		k := []byte(code)
		// 分组秘钥
		block, _ := aes.NewCipher(k)
		// 获取秘钥块的长度
		blockSize := block.BlockSize()
		// 补全码
		origData = PKCS7Padding(origData, blockSize)
		// 加密模式
		blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
		// 创建数组
		cryted := make([]byte, len(origData))
		// 加密
		blockMode.CryptBlocks(cryted, origData)

		return base64.StdEncoding.EncodeToString(cryted)
	}
	return ""
}

func Decrypt(code, cryted string) string {
	// 转成字节数组
	crytedByte, _ := base64.StdEncoding.DecodeString(cryted)
	k := []byte(code)

	// 分组秘钥
	block, _ := aes.NewCipher(k)
	// 获取秘钥块的长度
	blockSize := block.BlockSize()
	// 加密模式
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	// 创建数组
	orig := make([]byte, len(crytedByte))
	// 解密
	if len(orig)%blockMode.BlockSize() != 0 {
		return ""
	}
	blockMode.CryptBlocks(orig, crytedByte)
	// 去补全码
	orig = PKCS7UnPadding(orig)
	if orig == nil {
		log.Println("无法获得传入密码")
		return ""
	}
	return string(orig)
}

// PKCS7Padding 补码
func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

// PKCS7UnPadding 去码
func PKCS7UnPadding(origData []byte) []byte {
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
