package exhash

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	b64 "encoding/base64"

	"github.com/cockroachdb/errors"
)

// FSDecrypt https://open.feishu.cn/document/ukTMukTMukTM/uYDNxYjL2QTM24iN0EjN/event-subscription-configure-/encrypt-key-encryption-configuration-case
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
