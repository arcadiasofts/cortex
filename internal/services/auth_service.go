package services

import (
	"context"
	"crypto/ed25519"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

type AuthService struct {
	redis     *redis.Client
	jwtSecret []byte
}

func NewAuthService(r *redis.Client, secret string) *AuthService {
	return &AuthService{
		redis:     r,
		jwtSecret: []byte(secret),
	}git reset HEAD^
}

// Create challenge (Generate nonce and save to redis)
func (s *AuthService) CreateChallenge(ctx context.Context, did string) (string, error) {
	nonce := fmt.Sprintf("%d-challenge", time.Now().UnixNano())
	key := fmt.Sprintf("auth:nonce:%s", did)

	err := s.redis.Set(ctx, key, nonce, 120*time.Second).Err()
	if err != nil {
		return "", err
	}
	return nonce, nil
}

func (s *AuthService) GenerateJwt(did string) (string, error) {
	claims := jwt.MapClaims{
		"sub": did,
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}

func (s *AuthService) VerifyAndLogin(ctx context.Context, did string, nonce string, signature []byte) (string, error) {
	key := fmt.Sprintf("auth:nonce:%s", did)
	storedNonce, err := s.redis.Get(ctx, key).Result()

	if err == redis.Nil {
		return "", errors.New("invalid token")
	} else if err != nil {
		return "", err
	}

	if storedNonce != nonce {
		return "", errors.New("nonce mismatch")
	}

	parts := strings.Split(did, ":")
	if len(parts) != 3 {
		return "", errors.New("invalid did format")
	}

	pubKeyHex := parts[2]
	pubKeyBytes, err := hex.DecodeString(pubKeyHex)
	if err != nil || len(pubKeyBytes) != 32 {
		return "", errors.New("invalid public key")
	}

	isValid := ed25519.Verify(pubKeyBytes, []byte(nonce), signature)
	if !isValid {
		return "", errors.New("invalid signature")
	}

	s.redis.Del(ctx, key)

	return s.GenerateJwt(did)
}
