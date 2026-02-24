// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package expass

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewPwGen(t *testing.T) {
	tests := []struct {
		name    string
		length  int
		chars   string
		wantErr bool
		errMsg  string
	}{
		{
			name:    "正常生成数字密码",
			length:  10,
			chars:   "0123456789",
			wantErr: false,
		},
		{
			name:    "正常生成字母密码",
			length:  20,
			chars:   "abcdefghijklmnopqrstuvwxyz",
			wantErr: false,
		},
		{
			name:    "零长度",
			length:  0,
			chars:   "abc",
			wantErr: true,
			errMsg:  "length must be positive",
		},
		{
			name:    "负长度",
			length:  -5,
			chars:   "abc",
			wantErr: true,
			errMsg:  "length must be positive",
		},
		{
			name:    "空字符集",
			length:  10,
			chars:   "",
			wantErr: true,
			errMsg:  "chars cannot be empty",
		},
		{
			name:    "特殊字符集",
			length:  15,
			chars:   "!@#$%^&*()_+-=[]{}|;:,.<>?",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := NewPwGen(tt.length, tt.chars)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
				assert.Empty(t, result)
			} else {
				require.NoError(t, err)
				assert.Len(t, result, tt.length)

				// 验证所有字符都在指定字符集中
				for _, char := range result {
					assert.Contains(t, tt.chars, string(char))
				}
			}
		})
	}
}

func TestPwGenFunctions(t *testing.T) {
	t.Run("PwGenNum", func(t *testing.T) {
		result, err := PwGenNum(10)
		require.NoError(t, err)
		assert.Len(t, result, 10)

		// 验证只包含数字
		for _, char := range result {
			assert.True(t, char >= '0' && char <= '9', "期望数字，得到: %c", char)
		}
	})

	t.Run("PwGenAlpha", func(t *testing.T) {
		result, err := PwGenAlpha(15)
		require.NoError(t, err)
		assert.Len(t, result, 15)

		// 验证只包含字母
		for _, char := range result {
			isLetter := (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')
			assert.True(t, isLetter, "期望字母，得到: %c", char)
		}
	})

	t.Run("PwGenSymbols", func(t *testing.T) {
		result, err := PwGenSymbols(8)
		require.NoError(t, err)
		assert.Len(t, result, 8)

		// 验证只包含符号
		for _, char := range result {
			isLetter := (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z')
			isDigit := char >= '0' && char <= '9'
			assert.False(t, isLetter || isDigit, "不应包含字母或数字: %c", char)
		}
	})

	t.Run("PwGenAlphaNum", func(t *testing.T) {
		result, err := PwGenAlphaNum(20)
		require.NoError(t, err)
		assert.Len(t, result, 20)
	})

	t.Run("PwGenAlphaNumSymbols", func(t *testing.T) {
		result, err := PwGenAlphaNumSymbols(25)
		require.NoError(t, err)
		assert.Len(t, result, 25)
	})
}

func TestPasswordDistribution(t *testing.T) {
	// 测试分布均匀性
	chars := "ab"
	iterations := 10000
	length := 1

	counts := make(map[rune]int)
	for i := 0; i < iterations; i++ {
		result, err := NewPwGen(length, chars)
		require.NoError(t, err)
		counts[rune(result[0])]++
	}

	// 检查分布是否相对均匀（容差20%）
	expectedCount := iterations / len(chars)
	tolerance := float64(expectedCount) * 0.2

	for char, count := range counts {
		diff := float64(count - expectedCount)
		if diff < 0 {
			diff = -diff
		}
		assert.LessOrEqual(t, diff, tolerance,
			"字符 %c 分布不均: 期望约 %d 次，实际 %d 次", char, expectedCount, count)
	}
}

func TestGenerateHash(t *testing.T) {
	tests := []struct {
		name     string
		password []byte
		wantErr  bool
	}{
		{
			name:     "正常密码",
			password: []byte("mySecurePassword123!"),
			wantErr:  false,
		},
		{
			name:     "空密码",
			password: []byte{},
			wantErr:  true,
		},
		{
			name:     "短密码",
			password: []byte("123"),
			wantErr:  false,
		},
		{
			name:     "长密码",
			password: []byte(strings.Repeat("a", 72)), // bcrypt限制72字节
			wantErr:  false,
		},
		{
			name:     "包含特殊字符",
			password: []byte("密码@#$%^&*()_+中文"),
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := GenerateHash(tt.password)

			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, hash)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, hash)
				assert.NotEmpty(t, hash)

				// 验证相同密码生成的哈希不同（因为包含随机盐）
				hash2, err2 := GenerateHash(tt.password)
				require.NoError(t, err2)
				assert.NotEqual(t, hash, hash2)
			}
		})
	}
}

func TestCompareHash(t *testing.T) {
	password := []byte("testPassword123")
	wrongPassword := []byte("wrongPassword456")

	hash, err := GenerateHash(password)
	require.NoError(t, err)

	tests := []struct {
		name     string
		hash     []byte
		password []byte
		want     bool
	}{
		{
			name:     "正确密码",
			hash:     hash,
			password: password,
			want:     true,
		},
		{
			name:     "错误密码",
			hash:     hash,
			password: wrongPassword,
			want:     false,
		},
		{
			name:     "空密码",
			hash:     hash,
			password: []byte{},
			want:     false,
		},
		{
			name:     "无效哈希",
			hash:     []byte("invalid hash"),
			password: password,
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CompareHash(tt.hash, tt.password)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestGenerateHashWithCost(t *testing.T) {
	password := []byte("testPassword")

	t.Run("默认成本", func(t *testing.T) {
		hash, err := GenerateHashWithCost(password, 10)
		require.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.True(t, CompareHash(hash, password))
	})

	t.Run("最小成本", func(t *testing.T) {
		hash, err := GenerateHashWithCost(password, 4)
		require.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.True(t, CompareHash(hash, password))
	})

	t.Run("较高成本", func(t *testing.T) {
		hash, err := GenerateHashWithCost(password, 12)
		require.NoError(t, err)
		assert.NotEmpty(t, hash)
		assert.True(t, CompareHash(hash, password))
	})
}

func TestSaltPbkdf2Pass(t *testing.T) {
	tests := []struct {
		name       string
		salt       string
		password   string
		iterations int
	}{
		{
			name:       "默认迭代次数",
			salt:       "mysalt",
			password:   "mypassword",
			iterations: 0,
		},
		{
			name:       "自定义迭代次数",
			salt:       "anothersalt",
			password:   "anotherpassword",
			iterations: 100000,
		},
		{
			name:       "高迭代次数",
			salt:       "strongsalt",
			password:   "strongpassword",
			iterations: 1000000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SaltPbkdf2Pass(tt.salt, tt.password, tt.iterations)

			// 验证格式
			assert.True(t, strings.HasPrefix(result, "pbkdf2_sha256$"))
			parts := strings.Split(result, "$")
			assert.Len(t, parts, 4)
			assert.Equal(t, "pbkdf2_sha256", parts[0])
			assert.Equal(t, tt.salt, parts[2])

			// 验证相同输入产生相同输出
			result2 := SaltPbkdf2Pass(tt.salt, tt.password, tt.iterations)
			assert.Equal(t, result, result2)

			// 验证不同密码产生不同输出
			result3 := SaltPbkdf2Pass(tt.salt, tt.password+"diff", tt.iterations)
			assert.NotEqual(t, result, result3)
		})
	}
}

func TestAesEncryptDecrypt(t *testing.T) {
	tests := []struct {
		name     string
		data     string
		password string
		wantErr  bool
	}{
		{
			name:     "正常加解密",
			data:     "Hello, World!",
			password: "mySecretPassword",
			wantErr:  false,
		},
		{
			name:     "长文本",
			data:     strings.Repeat("Lorem ipsum dolor sit amet. ", 100),
			password: "anotherPassword123",
			wantErr:  false,
		},
		{
			name:     "包含特殊字符",
			data:     "数据加密测试@#$%^&*()_+-=[]{}|;:,.<>?",
			password: "特殊密码!@#",
			wantErr:  false,
		},
		{
			name:     "空数据",
			data:     "",
			password: "password",
			wantErr:  true,
		},
		{
			name:     "空密码",
			data:     "data",
			password: "",
			wantErr:  true,
		},
		{
			name:     "单字符数据",
			data:     "a",
			password: "password",
			wantErr:  false,
		},
		{
			name:     "单字符密码",
			data:     "test data",
			password: "p",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encrypted, err := AesEncryptGCM(tt.data, tt.password)

			if tt.wantErr {
				require.Error(t, err)
				assert.Empty(t, encrypted)
				return
			}

			require.NoError(t, err)
			assert.NotEmpty(t, encrypted)
			assert.NotEqual(t, tt.data, encrypted)

			// 验证相同输入产生不同输出（因为随机IV和盐）
			encrypted2, err2 := AesEncryptGCM(tt.data, tt.password)
			require.NoError(t, err2)
			assert.NotEqual(t, encrypted, encrypted2)

			// 解密测试
			decrypted, err := AesDecryptGCM(encrypted, tt.password)
			require.NoError(t, err)
			assert.Equal(t, tt.data, decrypted)

			// 错误密码解密失败
			_, err = AesDecryptGCM(encrypted, tt.password+"wrong")
			assert.Error(t, err)
		})
	}
}

func TestAesDecryptGCM_ErrorCases(t *testing.T) {
	validData := "test data"
	validPassword := "password"

	encrypted, err := AesEncryptGCM(validData, validPassword)
	require.NoError(t, err)

	tests := []struct {
		name     string
		data     string
		password string
		errMsg   string
	}{
		{
			name:     "空加密数据",
			data:     "",
			password: validPassword,
			errMsg:   "encrypted data cannot be empty",
		},
		{
			name:     "空密码",
			data:     encrypted,
			password: "",
			errMsg:   "password cannot be empty",
		},
		{
			name:     "无效base64",
			data:     "not-valid-base64!@#$",
			password: validPassword,
			errMsg:   "failed to decode base64",
		},
		{
			name:     "数据太短",
			data:     "dGVzdA==", // "test" in base64
			password: validPassword,
			errMsg:   "invalid encrypted data format",
		},
		{
			name:     "损坏的加密数据",
			data:     encrypted[:len(encrypted)-10] + "corrupted",
			password: validPassword,
			errMsg:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := AesDecryptGCM(tt.data, tt.password)
			require.Error(t, err)
			assert.Empty(t, result)
			if tt.errMsg != "" {
				assert.Contains(t, err.Error(), tt.errMsg)
			}
		})
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	assert.NotNil(t, config)
	assert.Equal(t, 600000, config.PBKDF2Iterations)
	assert.Equal(t, 32, config.SaltLength)
	assert.Equal(t, 10, config.BCryptCost) // bcrypt.DefaultCost
}

func BenchmarkNewPwGen(b *testing.B) {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	b.Run("Length10", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = NewPwGen(10, chars)
		}
	})

	b.Run("Length32", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = NewPwGen(32, chars)
		}
	})

	b.Run("Length128", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = NewPwGen(128, chars)
		}
	})
}

func BenchmarkGenerateHash(b *testing.B) {
	password := []byte("benchmarkPassword123!")

	for i := 0; i < b.N; i++ {
		_, _ = GenerateHash(password)
	}
}

func BenchmarkCompareHash(b *testing.B) {
	password := []byte("benchmarkPassword123!")
	hash, _ := GenerateHash(password)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CompareHash(hash, password)
	}
}

func BenchmarkAesEncryptGCM(b *testing.B) {
	data := "benchmark test data for encryption"
	password := "benchmarkPassword"

	for i := 0; i < b.N; i++ {
		_, _ = AesEncryptGCM(data, password)
	}
}

func BenchmarkAesDecryptGCM(b *testing.B) {
	data := "benchmark test data for decryption"
	password := "benchmarkPassword"
	encrypted, _ := AesEncryptGCM(data, password)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = AesDecryptGCM(encrypted, password)
	}
}
