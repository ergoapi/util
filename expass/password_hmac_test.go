// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package expass

import (
	"strings"
	"testing"
)

func TestAesEncryptDecryptCBCWithHMAC(t *testing.T) {
	testCases := []struct {
		name     string
		data     string
		password string
	}{
		{
			name:     "simple text",
			data:     "Hello, World!",
			password: "test-password-123",
		},
		{
			name:     "empty password edge case",
			data:     "test data",
			password: "",
		},
		{
			name:     "long text",
			data:     strings.Repeat("This is a test message. ", 100),
			password: "strong-password-with-symbols!@#$%",
		},
		{
			name:     "unicode text",
			data:     "测试中文字符串 🎉 Test Unicode",
			password: "unicode-password-测试",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test encryption with empty password
			if tc.password == "" {
				_, err := AesEncryptCBC(tc.data, tc.password)
				if err == nil || !strings.Contains(err.Error(), "password cannot be empty") {
					t.Errorf("Expected error for empty password, got: %v", err)
				}
				return
			}

			// Encrypt data
			encrypted, err := AesEncryptCBC(tc.data, tc.password)
			if err != nil {
				t.Fatalf("Encryption failed: %v", err)
			}

			// Encrypted data should be different from original
			if encrypted == tc.data {
				t.Error("Encrypted data should not match original data")
			}

			// Decrypt data
			decrypted, err := AesDecryptCBC(encrypted, tc.password)
			if err != nil {
				t.Fatalf("Decryption failed: %v", err)
			}

			// Decrypted data should match original
			if decrypted != tc.data {
				t.Errorf("Decrypted data does not match original.\nGot: %s\nWant: %s", decrypted, tc.data)
			}
		})
	}
}

func TestAesDecryptCBCWithHMACTamperDetection(t *testing.T) {
	data := "Secret message that should not be tampered with"
	password := "secure-password-123"

	// Encrypt the data
	encrypted, err := AesEncryptCBC(data, password)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Try to tamper with the encrypted data
	// Convert from base64, modify, and convert back
	tampered := encrypted[:len(encrypted)-10] + "TAMPERED!!"

	// Try to decrypt tampered data
	_, err = AesDecryptCBC(tampered, password)
	if err == nil {
		t.Error("Expected error when decrypting tampered data, but got none")
	} else if !strings.Contains(err.Error(), "base64") && !strings.Contains(err.Error(), "HMAC") && !strings.Contains(err.Error(), "invalid") {
		t.Errorf("Expected HMAC verification or format error, got: %v", err)
	}

	// Try with wrong password (should fail HMAC verification)
	_, err = AesDecryptCBC(encrypted, "wrong-password")
	if err == nil {
		t.Error("Expected error with wrong password, but got none")
	}
}

func TestAesEncryptCBCDeterministic(t *testing.T) {
	data := "Test deterministic encryption"
	password := "test-password"

	// Encrypt the same data twice
	encrypted1, err := AesEncryptCBC(data, password)
	if err != nil {
		t.Fatalf("First encryption failed: %v", err)
	}

	encrypted2, err := AesEncryptCBC(data, password)
	if err != nil {
		t.Fatalf("Second encryption failed: %v", err)
	}

	// Due to random salt and IV, encrypted results should be different
	if encrypted1 == encrypted2 {
		t.Error("Encrypted data should be different due to random salt and IV")
	}

	// But both should decrypt to the same original data
	decrypted1, err := AesDecryptCBC(encrypted1, password)
	if err != nil {
		t.Fatalf("First decryption failed: %v", err)
	}

	decrypted2, err := AesDecryptCBC(encrypted2, password)
	if err != nil {
		t.Fatalf("Second decryption failed: %v", err)
	}

	if decrypted1 != data || decrypted2 != data {
		t.Error("Both encrypted values should decrypt to the same original data")
	}
}

func TestAesDecryptCBCInvalidInput(t *testing.T) {
	testCases := []struct {
		name      string
		encrypted string
		password  string
		errMsg    string
	}{
		{
			name:      "empty encrypted data",
			encrypted: "",
			password:  "test",
			errMsg:    "encrypted data cannot be empty",
		},
		{
			name:      "invalid base64",
			encrypted: "not-valid-base64!@#$%",
			password:  "test",
			errMsg:    "base64",
		},
		{
			name:      "too short encrypted data",
			encrypted: "dGVzdA==", // "test" in base64, too short for our format
			password:  "test",
			errMsg:    "invalid encrypted data format",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := AesDecryptCBC(tc.encrypted, tc.password)
			if err == nil {
				t.Error("Expected error but got none")
			} else if !strings.Contains(strings.ToLower(err.Error()), strings.ToLower(tc.errMsg)) {
				t.Errorf("Expected error containing '%s', got: %v", tc.errMsg, err)
			}
		})
	}
}
