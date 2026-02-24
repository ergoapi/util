// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package exhash

import (
	"strings"
	"testing"
)

func TestFSEncryptDecryptGCM(t *testing.T) {
	testCases := []struct {
		name      string
		plaintext string
		key       string
	}{
		{
			name:      "simple text",
			plaintext: "Hello, Feishu!",
			key:       "test-key-123",
		},
		{
			name:      "chinese text",
			plaintext: "飞书加密测试",
			key:       "secure-key-456",
		},
		{
			name:      "long text",
			plaintext: strings.Repeat("Test message for Feishu. ", 50),
			key:       "long-secure-key-789",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			encrypted, err := FSEncryptGCM(tc.plaintext, tc.key)
			if err != nil {
				t.Fatalf("Encryption failed: %v", err)
			}

			// Encrypted should be different from plaintext
			if encrypted == tc.plaintext {
				t.Error("Encrypted data should not match plaintext")
			}

			decrypted, err := FSDecryptGCM(encrypted, tc.key)
			if err != nil {
				t.Fatalf("Decryption failed: %v", err)
			}

			// Decrypted should match original
			if decrypted != tc.plaintext {
				t.Errorf("Decrypted data does not match original.\nGot: %s\nWant: %s", decrypted, tc.plaintext)
			}
		})
	}
}

func TestFSDecryptGCMTamperDetection(t *testing.T) {
	plaintext := "Sensitive Feishu data"
	key := "secure-key-123"

	// Encrypt the data
	encrypted, err := FSEncryptGCM(plaintext, key)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Tamper with encrypted data
	tampered := encrypted[:len(encrypted)-5] + "HACK!"

	// Try to decrypt tampered data
	_, err = FSDecryptGCM(tampered, key)
	if err == nil {
		t.Error("Expected error when decrypting tampered data")
	} else if !strings.Contains(err.Error(), "authentication failed") && !strings.Contains(err.Error(), "base64") && !strings.Contains(err.Error(), "decryption failed") {
		t.Errorf("Expected authentication error, got: %v", err)
	}

	// Try with wrong key
	_, err = FSDecryptGCM(encrypted, "wrong-key")
	if err == nil {
		t.Error("Expected error with wrong key")
	} else if !strings.Contains(err.Error(), "authentication failed") && !strings.Contains(err.Error(), "decryption failed") {
		t.Errorf("Expected authentication error with wrong key, got: %v", err)
	}
}

func TestFSEncryptGCMInvalidInput(t *testing.T) {
	testCases := []struct {
		name      string
		plaintext string
		key       string
		errMsg    string
	}{
		{
			name:      "empty plaintext",
			plaintext: "",
			key:       "test-key",
			errMsg:    "plaintext cannot be empty",
		},
		{
			name:      "empty key",
			plaintext: "test data",
			key:       "",
			errMsg:    "key cannot be empty",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := FSEncryptGCM(tc.plaintext, tc.key)
			if err == nil {
				t.Error("Expected error but got none")
			} else if !strings.Contains(err.Error(), tc.errMsg) {
				t.Errorf("Expected error containing '%s', got: %v", tc.errMsg, err)
			}
		})
	}
}

func TestFSDecryptGCMInvalidInput(t *testing.T) {
	testCases := []struct {
		name    string
		encrypt string
		key     string
		errMsg  string
	}{
		{
			name:    "invalid base64",
			encrypt: "not-valid-base64!@#",
			key:     "test-key",
			errMsg:  "base64",
		},
		{
			name:    "too short data",
			encrypt: "dGVzdA==", // "test" in base64, too short
			key:     "test-key",
			errMsg:  "cipher too short",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := FSDecryptGCM(tc.encrypt, tc.key)
			if err == nil {
				t.Error("Expected error but got none")
			} else if !strings.Contains(strings.ToLower(err.Error()), strings.ToLower(tc.errMsg)) {
				t.Errorf("Expected error containing '%s', got: %v", tc.errMsg, err)
			}
		})
	}
}

// Test that original FSDecrypt still works for backward compatibility
func TestFSDecryptBackwardCompatibility(t *testing.T) {
	// This test ensures the original FSDecrypt function is unchanged
	// We're not testing its functionality extensively since it's for Feishu compatibility
	key := "test-key"

	// Test that it returns error for invalid input
	_, err := FSDecrypt("invalid-base64!", key)
	if err == nil {
		t.Error("Expected error for invalid base64")
	}

	// Test that it returns error for too short cipher
	_, err = FSDecrypt("dGVzdA==", key) // "test" in base64
	if err == nil {
		t.Error("Expected error for cipher too short")
	}
}
