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

func TestCheckPasswordStrength(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected PasswordStrength
	}{
		{
			name:     "空密码",
			password: "",
			expected: StrengthVeryWeak,
		},
		{
			name:     "短密码",
			password: "abc",
			expected: StrengthVeryWeak,
		},
		{
			name:     "纯数字短密码",
			password: "12345",
			expected: StrengthVeryWeak,
		},
		{
			name:     "8位简单密码",
			password: "password",
			expected: StrengthVeryWeak,
		},
		{
			name:     "8位混合密码",
			password: "Pass123!",
			expected: StrengthFair, // 8位长度，4种字符类型
		},
		{
			name:     "12位强密码",
			password: "MyP@ssw0rd12",
			expected: StrengthStrong, // 12位长度，4种字符类型
		},
		{
			name:     "16位超强密码",
			password: "MyV3ry$tr0ngP@ss",
			expected: StrengthVeryStrong, // 16位长度，4种字符类型
		},
		{
			name:     "20位复杂密码",
			password: "MyV3ry$tr0ngP@ssw0rd!",
			expected: StrengthVeryStrong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			strength := CheckPasswordStrength(tt.password)
			assert.Equal(t, tt.expected, strength)
			assert.NotEmpty(t, strength.String())
		})
	}
}

func TestPasswordStrengthString(t *testing.T) {
	tests := []struct {
		strength PasswordStrength
		expected string
	}{
		{StrengthVeryWeak, "Very Weak"},
		{StrengthWeak, "Weak"},
		{StrengthFair, "Fair"},
		{StrengthStrong, "Strong"},
		{StrengthVeryStrong, "Very Strong"},
		{PasswordStrength(999), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.strength.String())
		})
	}
}

func TestGenerateSecurePassword(t *testing.T) {
	tests := []struct {
		name           string
		length         int
		requireUpper   bool
		requireLower   bool
		requireDigit   bool
		requireSpecial bool
		wantErr        bool
	}{
		{
			name:           "所有要求",
			length:         16,
			requireUpper:   true,
			requireLower:   true,
			requireDigit:   true,
			requireSpecial: true,
			wantErr:        false,
		},
		{
			name:           "仅字母",
			length:         10,
			requireUpper:   true,
			requireLower:   true,
			requireDigit:   false,
			requireSpecial: false,
			wantErr:        false,
		},
		{
			name:           "长度太短",
			length:         3,
			requireUpper:   true,
			requireLower:   true,
			requireDigit:   true,
			requireSpecial: true,
			wantErr:        true,
		},
		{
			name:           "无要求",
			length:         12,
			requireUpper:   false,
			requireLower:   false,
			requireDigit:   false,
			requireSpecial: false,
			wantErr:        false,
		},
		{
			name:           "仅数字",
			length:         8,
			requireUpper:   false,
			requireLower:   false,
			requireDigit:   true,
			requireSpecial: false,
			wantErr:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			password, err := GenerateSecurePassword(
				tt.length,
				tt.requireUpper,
				tt.requireLower,
				tt.requireDigit,
				tt.requireSpecial,
			)

			if tt.wantErr {
				require.Error(t, err)
				assert.Empty(t, password)
			} else {
				require.NoError(t, err)
				assert.Len(t, password, tt.length)

				// 验证要求
				if tt.requireUpper {
					assert.True(t, containsUpperCase(password), "应包含大写字母")
				}
				if tt.requireLower {
					assert.True(t, containsLowerCase(password), "应包含小写字母")
				}
				if tt.requireDigit {
					assert.True(t, containsDigit(password), "应包含数字")
				}
				if tt.requireSpecial {
					assert.True(t, containsSpecial(password), "应包含特殊字符")
				}
			}
		})
	}
}

func TestContainsRune(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		r        rune
		expected bool
	}{
		{
			name:     "包含字符",
			s:        "hello",
			r:        'e',
			expected: true,
		},
		{
			name:     "不包含字符",
			s:        "hello",
			r:        'x',
			expected: false,
		},
		{
			name:     "空字符串",
			s:        "",
			r:        'a',
			expected: false,
		},
		{
			name:     "中文字符",
			s:        "你好世界",
			r:        '好',
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := containsRune(tt.s, tt.r)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Helper functions for testing
func containsUpperCase(s string) bool {
	return strings.ContainsAny(s, "ABCDEFGHIJKLMNOPQRSTUVWXYZ")
}

func containsLowerCase(s string) bool {
	return strings.ContainsAny(s, "abcdefghijklmnopqrstuvwxyz")
}

func containsDigit(s string) bool {
	return strings.ContainsAny(s, "0123456789")
}

func containsSpecial(s string) bool {
	for _, r := range s {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9')) {
			return true
		}
	}
	return false
}

func BenchmarkCheckPasswordStrength(b *testing.B) {
	password := "MyV3ry$tr0ngP@ssw0rd!"

	for i := 0; i < b.N; i++ {
		_ = CheckPasswordStrength(password)
	}
}

func BenchmarkGenerateSecurePassword(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GenerateSecurePassword(16, true, true, true, true)
	}
}
