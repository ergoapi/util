package expass

import (
	"fmt"
	"testing"
)

func TestGenerateHash(t *testing.T) {
	hashByte, err := GenerateHash([]byte("123456"))
	if err != nil {
		return
	}
	fmt.Println(string(hashByte))
	if CompareHash(hashByte, []byte("123456")) {
		fmt.Println("123456密码正确")
	}
	if !CompareHash(hashByte, []byte("1234567")) {
		fmt.Println("1234567密码错误")
	}
	code := "abcd"
	encrypted, err := AesEncryptCBC("123456", code)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(encrypted)
	decrypted, err := AesDecryptCBC(encrypted, code)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(decrypted)
}
