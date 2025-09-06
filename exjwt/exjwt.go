// Copyright (c) 2025-2025 All rights reserved.
//
// The original source code is licensed under the DO WHAT THE FUCK YOU WANT TO PUBLIC LICENSE.
//
// You may review the terms of licenses in the LICENSE file.

// Package exjwt provides utilities.
package exjwt

import (
	"os"
	"time"

	"github.com/cockroachdb/errors"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// Exported errors for callers to detect specific failure reasons.
	ErrNoSecret     = errors.New("jwt: secret not set")
	ErrEmptySecret  = errors.New("jwt: empty secret")
	ErrInvalidToken = errors.New("jwt: token invalid")
	ErrAlgMismatch  = errors.New("jwt: signing method mismatch")
	ErrExpired      = errors.New("jwt: token expired")
)

const envJWTSecret = "JWT_SECRET"

// Claims represents the JWT claims used by this package.
type Claims struct {
	Username string `json:"username"`
	UUID     string `json:"uuid"`
	jwt.RegisteredClaims
}

func secretFromEnv() ([]byte, error) {
	s := os.Getenv(envJWTSecret)
	if s == "" {
		return nil, ErrNoSecret
	}
	return []byte(s), nil
}

// Auth signs a JWT using secret from environment variable JWT_SECRET.
func Auth(username string, uuid string) (t string, err error) {
	key, err := secretFromEnv()
	if err != nil {
		return "", err
	}
	return AuthWithSecret(username, uuid, key)
}

// AuthWithSecret signs a JWT using the provided secret.
// TTL policy: admin/root -> 4h, others -> 7d.
func AuthWithSecret(username string, uuid string, key []byte) (t string, err error) {
	if len(key) == 0 {
		return "", ErrEmptySecret
	}

	now := time.Now()
	var ttl time.Duration
	if username == "admin" || username == "root" {
		ttl = 4 * time.Hour
	} else {
		ttl = 7 * 24 * time.Hour
	}

	claims := &Claims{
		Username: username,
		UUID:     uuid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   uuid,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err = token.SignedString(key)
	if err != nil {
		return "", errors.Wrap(err, "JWT generate failure")
	}
	return t, nil
}

func Parse(ts string) (jwt.MapClaims, error) {
	key, err := secretFromEnv()
	if err != nil {
		return nil, err
	}
	return ParseWithSecret(ts, key)
}

// ParseWithSecret parses and validates a token with the provided secret.
func ParseWithSecret(ts string, key []byte) (jwt.MapClaims, error) {
	if len(key) == 0 {
		return nil, ErrEmptySecret
	}

	c := &Claims{}
	token, err := jwt.ParseWithClaims(ts, c, func(token *jwt.Token) (any, error) {
		if m, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || m != jwt.SigningMethodHS256 {
			return nil, ErrAlgMismatch
		}
		return key, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithLeeway(time.Minute))
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, ErrInvalidToken
	}
	// Additional explicit expiry check for clarity.
	if c.ExpiresAt != nil && time.Now().After(c.ExpiresAt.Time) {
		return nil, ErrExpired
	}

	// Return a compatible MapClaims view for callers expecting map.
	mc := jwt.MapClaims{
		"username": c.Username,
		"uuid":     c.UUID,
	}
	if c.ExpiresAt != nil {
		mc["exp"] = c.ExpiresAt.Time.Unix()
	}
	if c.IssuedAt != nil {
		mc["iat"] = c.IssuedAt.Time.Unix()
	}
	if c.NotBefore != nil {
		mc["nbf"] = c.NotBefore.Time.Unix()
	}
	if c.Subject != "" {
		mc["sub"] = c.Subject
	}
	return mc, nil
}
