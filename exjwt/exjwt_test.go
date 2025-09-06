// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

package exjwt

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func setEnv(key, val string) func() {
	old, ok := os.LookupEnv(key)
	if val == "" {
		_ = os.Unsetenv(key)
	} else {
		_ = os.Setenv(key, val)
	}
	return func() {
		if ok {
			_ = os.Setenv(key, old)
		} else {
			_ = os.Unsetenv(key)
		}
	}
}

func TestAuthRequiresSecret(t *testing.T) {
	defer setEnv(envJWTSecret, "")()
	if _, err := Auth("alice", "uuid-1"); !errors.Is(err, ErrNoSecret) {
		t.Fatalf("expected ErrNoSecret, got %v", err)
	}
}

func TestAuthAndParseNormal(t *testing.T) {
	defer setEnv(envJWTSecret, "test-secret-1234567890")()

	user := "alice"
	uid := "uuid-123"
	token, err := Auth(user, uid)
	if err != nil {
		t.Fatalf("Auth failed: %v", err)
	}

	claims, err := Parse(token)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if got := claims["username"]; got != user {
		t.Fatalf("username mismatch: got %v want %v", got, user)
	}
	if got := claims["uuid"]; got != uid {
		t.Fatalf("uuid mismatch: got %v want %v", got, uid)
	}
	// Check presence of standard fields
	if _, ok := claims["exp"]; !ok {
		t.Fatalf("exp not set")
	}
	if _, ok := claims["iat"]; !ok {
		t.Fatalf("iat not set")
	}
	if _, ok := claims["nbf"]; !ok {
		t.Fatalf("nbf not set")
	}
	if got, ok := claims["sub"]; !ok || got != uid {
		t.Fatalf("sub mismatch: got %v want %v (ok=%v)", got, uid, ok)
	}

	// exp - iat should be approximately 7 days for normal users
	exp, _ := claims["exp"].(int64)
	iat, _ := claims["iat"].(int64)
	if exp == 0 || iat == 0 {
		t.Fatalf("bad exp/iat types: exp=%T iat=%T", claims["exp"], claims["iat"])
	}
	diff := exp - iat
	// 7 days in seconds
	want := int64(7 * 24 * 60 * 60)
	if diff < want-100 || diff > want+100 { // allow ~100s skew in CI
		t.Fatalf("unexpected TTL: got %d, want ~%d", diff, want)
	}
}

func TestAuthAdminTTL(t *testing.T) {
	defer setEnv(envJWTSecret, "test-secret-admin")()
	token, err := Auth("admin", "uuid-admin")
	if err != nil {
		t.Fatalf("Auth failed: %v", err)
	}
	claims, err := Parse(token)
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	exp, _ := claims["exp"].(int64)
	iat, _ := claims["iat"].(int64)
	got := exp - iat
	// 4 hours in seconds
	want := int64(4 * 60 * 60)
	// allow some tolerance
	if got < want-50 || got > want+50 {
		t.Fatalf("admin TTL mismatch: got %d want ~%d", got, want)
	}
}

func TestParseAlgorithmMismatch(t *testing.T) {
	key := []byte("sign-key-512")
	defer setEnv(envJWTSecret, string(key))()

	now := time.Now()
	c := &Claims{
		Username: "bob",
		UUID:     "uuid-bob",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(10 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   "uuid-bob",
		},
	}
	// Sign with HS512 while parser only accepts HS256
	tok := jwt.NewWithClaims(jwt.SigningMethodHS512, c)
	s, err := tok.SignedString(key)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	if _, err := Parse(s); err == nil {
		t.Fatalf("expected error for alg mismatch, got nil")
	} else if !strings.Contains(err.Error(), "signing method HS512 is invalid") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestParseExpired(t *testing.T) {
	key := []byte("sign-key-expired")
	defer setEnv(envJWTSecret, string(key))()
	past := time.Now().Add(-2 * time.Minute) // beyond leeway
	c := &Claims{
		Username: "eve",
		UUID:     "uuid-eve",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(past),
			IssuedAt:  jwt.NewNumericDate(past.Add(-time.Minute)),
			NotBefore: jwt.NewNumericDate(past.Add(-time.Minute)),
			Subject:   "uuid-eve",
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	s, err := tok.SignedString(key)
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	if _, err := Parse(s); err == nil {
		t.Fatalf("expected error for expired token, got nil")
	} else if !strings.Contains(err.Error(), "expired") {
		t.Fatalf("unexpected error: %v", err)
	}
}
func Example() {
	defer setEnv(envJWTSecret, "my-very-secret-key")()
	token, err := Auth("alice", "uuid-123")
	if err != nil {
		fmt.Println("auth error:", err)
		return
	}
	claims, err := Parse(token)
	if err != nil {
		fmt.Println("parse error:", err)
		return
	}
	fmt.Printf("username=%v uuid=%v sub=%v\n", claims["username"], claims["uuid"], claims["sub"])
	// Output: username=alice uuid=uuid-123 sub=uuid-123
}
