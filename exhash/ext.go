// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package exhash

import (
	"crypto/aes"
	"crypto/cipher"
	cryptorand "crypto/rand"
	"crypto/sha256"
	b64 "encoding/base64"
	"io"

	"github.com/cockroachdb/errors"
)

// FSDecrypt https://open.feishu.cn/document/ukTMukTMukTM/uYDNxYjL2QTM24iN0EjN/event-subscription-configure-/encrypt-key-encryption-configuration-case
// Note: This function maintains compatibility with Feishu API. For new implementations, consider using FSDecryptGCM.
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

// FSEncryptGCM encrypts data using AES-GCM for authenticated encryption.
// GCM mode provides both confidentiality and authenticity.
func FSEncryptGCM(plaintext string, key string) (string, error) {
	if plaintext == "" {
		return "", errors.New("plaintext cannot be empty")
	}
	if key == "" {
		return "", errors.New("key cannot be empty")
	}

	// Derive key from the main key
	keyBs := sha256.Sum256([]byte(key))
	encKey := keyBs[:32] // Use full 256-bit key for AES-256

	block, err := aes.NewCipher(encKey)
	if err != nil {
		return "", errors.Wrap(err, "create cipher error")
	}

	// Create GCM mode
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.Wrap(err, "create GCM error")
	}

	// Generate cryptographically secure random nonce
	nonce := make([]byte, aesgcm.NonceSize())
	if _, err := io.ReadFull(cryptorand.Reader, nonce); err != nil {
		return "", errors.Wrap(err, "generate nonce error")
	}

	// Encrypt and authenticate
	encrypted := aesgcm.Seal(nil, nonce, []byte(plaintext), nil)

	// Combine nonce + encrypted (which includes the auth tag)
	result := append(nonce, encrypted...)

	return b64.StdEncoding.EncodeToString(result), nil
}

// FSDecryptGCM decrypts data encrypted by FSEncryptGCM using AES-GCM.
func FSDecryptGCM(encrypt string, key string) (string, error) {
	buf, err := b64.StdEncoding.DecodeString(encrypt)
	if err != nil {
		return "", errors.Wrap(err, "base64 decode error")
	}

	// Derive key
	keyBs := sha256.Sum256([]byte(key))
	encKey := keyBs[:32] // Use full 256-bit key for AES-256

	block, err := aes.NewCipher(encKey)
	if err != nil {
		return "", errors.Wrap(err, "create cipher error")
	}

	// Create GCM mode
	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", errors.Wrap(err, "create GCM error")
	}

	// Minimum size: nonce + ciphertext + auth tag
	if len(buf) < aesgcm.NonceSize()+aesgcm.Overhead() {
		return "", errors.New("cipher too short")
	}

	// Extract nonce and ciphertext
	nonce := buf[:aesgcm.NonceSize()]
	ciphertext := buf[aesgcm.NonceSize():]

	// Decrypt and verify
	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", errors.Wrap(err, "decryption failed")
	}

	return string(plaintext), nil
}

// Deprecated: Use FSEncryptGCM instead. The implementation now uses AES-GCM;
// this wrapper preserves backward compatibility with the old name.
func FSEncryptWithHMAC(plaintext string, key string) (string, error) {
	return FSEncryptGCM(plaintext, key)
}

// Deprecated: Use FSDecryptGCM instead. The implementation now uses AES-GCM;
// this wrapper preserves backward compatibility with the old name.
func FSDecryptWithHMAC(encrypt string, key string) (string, error) {
	return FSDecryptGCM(encrypt, key)
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
