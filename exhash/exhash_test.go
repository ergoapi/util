// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package exhash

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestB64Encode(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", ""},
		{"hello", "hello", "aGVsbG8="},
		{"world", "world", "d29ybGQ="},
		{"special", "Hello, 世界!", "SGVsbG8sIOS4lueVjCE="},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := B64Encode(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestB64Decode(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"empty", "", "", true},
		{"hello", "aGVsbG8=", "hello", false},
		{"world", "d29ybGQ=", "world", false},
		{"special", "SGVsbG8sIOS4lueVjCE=", "Hello, 世界!", false},
		{"invalid", "invalid!@#", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := B64Decode(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestB64EncodeBytes(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  string
	}{
		{"empty", []byte{}, ""},
		{"hello", []byte("hello"), "aGVsbG8="},
		{"binary", []byte{0x01, 0x02, 0x03}, "AQID"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := B64EncodeBytes(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestB64DecodeBytes(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []byte
		wantErr bool
	}{
		{"empty", "", nil, true},
		{"hello", "aGVsbG8=", []byte("hello"), false},
		{"binary", "AQID", []byte{0x01, 0x02, 0x03}, false},
		{"invalid", "!!!!", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := B64DecodeBytes(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestB58Encode(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  string
	}{
		{"empty", []byte{}, ""},
		{"hello", []byte("hello"), "Cn8eVZg"},
		{"leading zeros", []byte{0x00, 0x00, 0x01}, "112"},
		{"single zero", []byte{0x00}, "1"},
		{"multiple zeros", []byte{0x00, 0x00, 0x00}, "111"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := B58Encode(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestB58Decode(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []byte
		wantErr bool
	}{
		{"empty", "", nil, true},
		{"hello", "Cn8eVZg", []byte("hello"), false},
		{"leading ones", "112", []byte{0x00, 0x00, 0x01}, false},
		{"single one", "1", []byte{0x00}, false},
		{"multiple ones", "111", []byte{0x00, 0x00, 0x00}, false},
		{"invalid char", "0OIl", nil, true}, // 0, O, I, l are not in base58
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := B58Decode(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestB58DecodeString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"empty", "", "", true},
		{"hello", "Cn8eVZg", "hello", false},
		{"invalid", "0OIl", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := B58DecodeString(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestB58RoundTrip(t *testing.T) {
	// 测试编解码往返
	inputs := [][]byte{
		[]byte("hello"),
		[]byte("world"),
		[]byte("The quick brown fox jumps over the lazy dog"),
		{0x00, 0x01, 0x02, 0x03},
		{0x00, 0x00, 0x00},
		{0xff, 0xfe, 0xfd},
	}

	for _, input := range inputs {
		encoded := B58Encode(input)
		decoded, err := B58Decode(encoded)
		assert.NoError(t, err)
		assert.Equal(t, input, decoded, "Round trip failed for %v", input)
	}
}

func TestB32Encode(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", ""},
		{"hello", "hello", "NBSWY3DP"},
		{"world", "world", "O5XXE3DE"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := B32Encode(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestB32Decode(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"empty", "", "", true},
		{"hello", "NBSWY3DP", "hello", false},
		{"world", "O5XXE3DE", "world", false},
		{"invalid", "!!!!", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := B32Decode(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestReverseBytes(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  []byte
	}{
		{"empty", []byte{}, []byte{}},
		{"single", []byte{1}, []byte{1}},
		{"even", []byte{1, 2, 3, 4}, []byte{4, 3, 2, 1}},
		{"odd", []byte{1, 2, 3}, []byte{3, 2, 1}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := reverseBytes(tt.input)
			assert.Equal(t, tt.want, got)
			// 确保原始切片没有被修改（空切片除外）
			if len(tt.input) > 0 {
				assert.NotSame(t, &tt.input[0], &got[0])
			}
		})
	}
}

func TestGenSha256(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{"hello", "hello", "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824"},
		{"world", "world", "486ea46224d1bb4fb680f34f7c9ad96a8f24ec88be73ea8e5a6c65260e9cb8a7"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenSha256(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGenSha512(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", "cf83e1357eefb8bdf1542850d66d8007d620e4050b5715dc83f4a921d36ce9ce47d0d13c5d85f2b0ff8318d2877eec2f63b931bd47417a81a538327af927da3e"},
		{"hello", "hello", "9b71d224bd62f3785d96d46ad3ea3d73319bfbc2890caadae2dff72519673ca72323c3d99ba5c11d7c7acc6e14b8c5da0c4663475c2e5c3adef46f73bcdec043"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenSha512(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGenSha1(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", "da39a3ee5e6b4b0d3255bfef95601890afd80709"},
		{"hello", "hello", "aaf4c61ddcc5e8a2dabede0f3b482cd9aea9434d"},
		{"world", "world", "7c211433f02071597741e6ff5a8ea34789abbf43"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenSha1(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestStringToNumber(t *testing.T) {
	// 测试相同的字符串产生相同的哈希值
	str := "test string"
	hash1 := StringToNumber(str)
	hash2 := StringToNumber(str)
	assert.Equal(t, hash1, hash2, "Same string should produce same hash")

	// 测试不同的字符串产生不同的哈希值
	hash3 := StringToNumber("different string")
	assert.NotEqual(t, hash1, hash3, "Different strings should produce different hashes")
}

func TestString(t *testing.T) {
	// 测试相同的字符串产生相同的哈希值
	str := "test string"
	hash1 := String(str)
	hash2 := String(str)
	assert.Equal(t, hash1, hash2, "Same string should produce same hash")
	assert.Len(t, hash1, 64, "SHA256 hash should be 64 hex characters")

	// 测试不同的字符串产生不同的哈希值
	hash3 := String("different string")
	assert.NotEqual(t, hash1, hash3, "Different strings should produce different hashes")
}

func TestHex(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		length int
		want   int // 期望的长度
	}{
		{"short", "test", 8, 8},
		{"exact", "test", 16, 16},
		{"long", "test", 20, 20},
		{"zero length", "test", 0, 32},
		{"negative length", "test", -1, 32},
		{"exceed md5 length", "test", 50, 32},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Hex(tt.input, tt.length)
			assert.Len(t, got, tt.want)
			// 验证是十六进制字符串
			for _, c := range got {
				assert.True(t, (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f'))
			}
		})
	}
}

func TestMD5(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", "d41d8cd98f00b204e9800998ecf8427e"},
		{"hello", "hello", "5d41402abc4b2a76b9719d911017c592"},
		{"world", "world", "7d793037a0760186574b0282f2f435e7"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MD5(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFSDecrypt(t *testing.T) {
	// 这个测试需要一个有效的加密字符串
	// 由于FSDecrypt是用于解密飞书的数据，我们需要模拟一个加密过程
	t.Run("invalid base64", func(t *testing.T) {
		_, err := FSDecrypt("invalid!@#", "key")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "base64")
	})

	t.Run("cipher too short", func(t *testing.T) {
		shortCipher := B64Encode("short")
		_, err := FSDecrypt(shortCipher, "key")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cipher too short")
	})

	t.Run("invalid block size", func(t *testing.T) {
		// 创建一个长度不是块大小整数倍的密文
		invalidCipher := make([]byte, 17) // IV(16) + 1字节密文
		encoded := B64EncodeBytes(invalidCipher)
		_, err := FSDecrypt(encoded, "key")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not a multiple")
	})
}

func TestRemovePKCS7Padding(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		want    []byte
		wantErr bool
	}{
		{
			name:    "empty data",
			input:   []byte{},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "valid padding 1",
			input:   []byte("hello\x01"),
			want:    []byte("hello"),
			wantErr: false,
		},
		{
			name:    "valid padding 3",
			input:   []byte("hello\x03\x03\x03"),
			want:    []byte("hello"),
			wantErr: false,
		},
		{
			name:    "full block padding",
			input:   bytes.Repeat([]byte{16}, 16),
			want:    []byte{},
			wantErr: false,
		},
		{
			name:    "invalid padding size",
			input:   []byte("hello\x20"), // padding > block size
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid padding bytes",
			input:   []byte("hello\x03\x02\x03"), // inconsistent padding
			want:    nil,
			wantErr: true,
		},
		{
			name:    "padding larger than data",
			input:   []byte("\x10"),
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := removePKCS7Padding(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func BenchmarkB64Encode(b *testing.B) {
	data := strings.Repeat("hello world", 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = B64Encode(data)
	}
}

func BenchmarkB58Encode(b *testing.B) {
	data := []byte(strings.Repeat("hello world", 100))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = B58Encode(data)
	}
}

func BenchmarkGenSha256(b *testing.B) {
	data := strings.Repeat("hello world", 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = GenSha256(data)
	}
}

func BenchmarkMD5(b *testing.B) {
	data := strings.Repeat("hello world", 100)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = MD5(data)
	}
}

// 测试并发安全性
func TestConcurrentSafety(t *testing.T) {
	data := "test data"
	iterations := 100

	t.Run("B64Encode concurrent", func(t *testing.T) {
		for i := 0; i < iterations; i++ {
			go func() {
				result := B64Encode(data)
				assert.NotEmpty(t, result)
			}()
		}
	})

	t.Run("B58Encode concurrent", func(t *testing.T) {
		for i := 0; i < iterations; i++ {
			go func() {
				result := B58Encode([]byte(data))
				assert.NotEmpty(t, result)
			}()
		}
	})

	t.Run("GenSha256 concurrent", func(t *testing.T) {
		for i := 0; i < iterations; i++ {
			go func() {
				result := GenSha256(data)
				assert.NotEmpty(t, result)
			}()
		}
	})
}

// 测试边界条件
func TestBoundaryConditions(t *testing.T) {
	t.Run("B58 max byte value", func(t *testing.T) {
		data := []byte{0xff, 0xff, 0xff, 0xff}
		encoded := B58Encode(data)
		decoded, err := B58Decode(encoded)
		require.NoError(t, err)
		assert.Equal(t, data, decoded)
	})

	t.Run("B58 large data", func(t *testing.T) {
		data := make([]byte, 1024)
		for i := range data {
			data[i] = byte(i % 256)
		}
		encoded := B58Encode(data)
		decoded, err := B58Decode(encoded)
		require.NoError(t, err)
		assert.Equal(t, data, decoded)
	})

	t.Run("Base64 large data", func(t *testing.T) {
		data := strings.Repeat("x", 10000)
		encoded := B64Encode(data)
		decoded, err := B64Decode(encoded)
		require.NoError(t, err)
		assert.Equal(t, data, decoded)
	})
}
