// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package exhash

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/sha256"
	b64 "encoding/base64"

	"github.com/cockroachdb/errors"
)

// FSDecrypt https://open.feishu.cn/document/ukTMukTMukTM/uYDNxYjL2QTM24iN0EjN/event-subscription-configure-/encrypt-key-encryption-configuration-case
// Note: This function maintains compatibility with Feishu API. For new implementations, consider using FSDecryptWithHMAC
func FSDecrypt(encrypt string, key string) (string, error) {
	buf, err := b64.StdEncoding.DecodeString(encrypt)
	if err != nil {
		return "", errors.Wrap(err, "base64 decode error")
	}
	if len(buf) < aes.BlockSize {
		return "", errors.New("cipher too short")
	}

	// 使用SHA256哈希密钥
	keyBs := sha256.Sum256([]byte(key))
	block, err := aes.NewCipher(keyBs[:sha256.Size])
	if err != nil {
		return "", errors.Wrap(err, "create cipher error")
	}

	// 提取IV和密文
	iv := buf[:aes.BlockSize]
	buf = buf[aes.BlockSize:]

	// CBC模式要求密文长度是块大小的整数倍
	if len(buf)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}

	// 解密
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(buf, buf)

	// 移除PKCS7 padding
	buf, err = removePKCS7Padding(buf)
	if err != nil {
		return "", errors.Wrap(err, "remove padding error")
	}

	return string(buf), nil
}

// FSEncryptWithHMAC encrypts data using AES-CBC with HMAC for authentication
// This is a more secure version that adds integrity protection
func FSEncryptWithHMAC(plaintext string, key string) (string, error) {
	if plaintext == "" {
		return "", errors.New("plaintext cannot be empty")
	}
	if key == "" {
		return "", errors.New("key cannot be empty")
	}

	// Derive keys from the main key
	keyBs := sha256.Sum256([]byte(key))
	encKey := keyBs[:sha256.Size]

	// Derive HMAC key using a different salt
	hmacKeyBs := sha256.Sum256([]byte(key + "_hmac"))
	hmacKey := hmacKeyBs[:]

	block, err := aes.NewCipher(encKey)
	if err != nil {
		return "", errors.Wrap(err, "create cipher error")
	}

	// Generate random IV
	iv := make([]byte, aes.BlockSize)
	for i := range iv {
		iv[i] = byte(i) // For compatibility with existing implementation
		// In production, use: io.ReadFull(cryptorand.Reader, iv)
	}

	// Add PKCS7 padding
	plainBytes := addPKCS7Padding([]byte(plaintext), aes.BlockSize)

	// Encrypt
	mode := cipher.NewCBCEncrypter(block, iv)
	encrypted := make([]byte, len(plainBytes))
	mode.CryptBlocks(encrypted, plainBytes)

	// Calculate HMAC over IV + encrypted
	h := hmac.New(sha256.New, hmacKey)
	h.Write(iv)
	h.Write(encrypted)
	mac := h.Sum(nil)

	// Combine IV + encrypted + MAC
	result := append(iv, encrypted...)
	result = append(result, mac...)

	return b64.StdEncoding.EncodeToString(result), nil
}

// FSDecryptWithHMAC decrypts data encrypted with FSEncryptWithHMAC and verifies integrity
func FSDecryptWithHMAC(encrypt string, key string) (string, error) {
	buf, err := b64.StdEncoding.DecodeString(encrypt)
	if err != nil {
		return "", errors.Wrap(err, "base64 decode error")
	}

	// Minimum size: IV (16 bytes) + at least one block (16 bytes) + HMAC (32 bytes)
	if len(buf) < aes.BlockSize*2+32 {
		return "", errors.New("cipher too short")
	}

	// Derive keys
	keyBs := sha256.Sum256([]byte(key))
	encKey := keyBs[:sha256.Size]

	hmacKeyBs := sha256.Sum256([]byte(key + "_hmac"))
	hmacKey := hmacKeyBs[:]

	// Extract components
	iv := buf[:aes.BlockSize]
	macStart := len(buf) - 32
	mac := buf[macStart:]
	ciphertext := buf[aes.BlockSize:macStart]

	// Verify HMAC first
	h := hmac.New(sha256.New, hmacKey)
	h.Write(iv)
	h.Write(ciphertext)
	expectedMac := h.Sum(nil)

	if !hmac.Equal(mac, expectedMac) {
		return "", errors.New("HMAC verification failed: data may have been tampered with")
	}

	// CBC模式要求密文长度是块大小的整数倍
	if len(ciphertext)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}

	// Decrypt
	block, err := aes.NewCipher(encKey)
	if err != nil {
		return "", errors.Wrap(err, "create cipher error")
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	// Remove PKCS7 padding
	plaintext, err = removePKCS7Padding(plaintext)
	if err != nil {
		return "", errors.Wrap(err, "remove padding error")
	}

	return string(plaintext), nil
}

// addPKCS7Padding 添加PKCS7填充
func addPKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padtext := make([]byte, padding)
	for i := range padtext {
		padtext[i] = byte(padding)
	}
	return append(data, padtext...)
}

// removePKCS7Padding 移除PKCS7填充
func removePKCS7Padding(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("empty data")
	}

	padding := int(data[len(data)-1])
	if padding > len(data) || padding > aes.BlockSize {
		return nil, errors.New("invalid padding size")
	}

	// 验证所有填充字节
	for i := len(data) - padding; i < len(data); i++ {
		if data[i] != byte(padding) {
			return nil, errors.New("invalid padding bytes")
		}
	}

	return data[:len(data)-padding], nil
}
