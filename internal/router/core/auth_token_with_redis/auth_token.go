package auth_token

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// TokenPayload contains the essential information stored in a JWT token
type TokenPayload struct {
	UserID   int64    `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	Source   string   `json:"source"`
}

// AuthToken represents the authentication token manager
type AuthToken struct {
	redisClient *redis.Client
	jwtSecret   []byte
	tokenExpire time.Duration
}

// New creates a new AuthToken instance
func New(redisClient *redis.Client, jwtSecret string, expireSeconds int) *AuthToken {
	return &AuthToken{
		redisClient: redisClient,
		jwtSecret:   []byte(jwtSecret),
		tokenExpire: time.Duration(expireSeconds) * time.Second,
	}
}

// Issue generates a JWT token and stores the session in Redis
func (at *AuthToken) Issue(payload TokenPayload) (string, error) {
	// Generate a unique identifier for the token
	jti := uuid.New().String()

	// Create claims with the payload and jti
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id":  payload.UserID,
		"username": payload.Username,
		"roles":    payload.Roles,
		"source":   payload.Source,
		"jti":      jti,
		"iat":      now.Unix(),
		"exp":      now.Add(at.tokenExpire).Unix(),
		"nbf":      now.Unix(),
	}

	// Create the token with HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret
	signedToken, err := token.SignedString(at.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	// Store the jti in Redis with the key user_session:{UserID}
	key := fmt.Sprintf("user_session:%d", payload.UserID)
	err = at.redisClient.Set(context.Background(), key, jti, at.tokenExpire).Err()
	if err != nil {
		return "", fmt.Errorf("failed to store session in Redis: %w", err)
	}

	return signedToken, nil
}

// Revoke removes the user session from Redis, effectively logging out the user
func (at *AuthToken) Revoke(userID int64) error {
	key := fmt.Sprintf("user_session:%d", userID)
	err := at.redisClient.Del(context.Background(), key).Err()
	if err != nil {
		return fmt.Errorf("failed to revoke session for user %d: %w", userID, err)
	}
	return nil
}

// ValidateToken checks if the token is valid and not revoked
func (at *AuthToken) ValidateToken(tokenString string) (*TokenPayload, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return at.jwtSecret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	// Verify the token hasn't been revoked by comparing jti with Redis
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid user_id in token")
	}
	userID := int64(userIDFloat)

	jti, ok := claims["jti"].(string)
	if !ok {
		return nil, errors.New("invalid jti in token")
	}

	// Get the stored jti from Redis
	ctx := context.Background()
	storedJTI, err := at.redisClient.Get(ctx, fmt.Sprintf("user_session:%d", userID)).Result()
	if err != nil {
		if err == redis.Nil {
			// No session found in Redis, token is invalid (revoked)
			return nil, errors.New("session not found in Redis, token may be revoked")
		}
		return nil, fmt.Errorf("failed to get session from Redis: %w", err)
	}

	// Compare the jti from token with the one stored in Redis
	if jti != storedJTI {
		return nil, errors.New("token has been revoked (jti mismatch)")
	}

	// Build and return the TokenPayload
	payload := &TokenPayload{
		UserID:   userID,
		Username: claims["username"].(string),
		Source:   claims["source"].(string),
	}

	// Handle roles array
	if rolesClaim, ok := claims["roles"].([]interface{}); ok {
		roles := make([]string, len(rolesClaim))
		for i, v := range rolesClaim {
			roles[i] = v.(string)
		}
		payload.Roles = roles
	} else if rolesStr, ok := claims["roles"].(string); ok {
		// Handle roles as a single string (fallback)
		payload.Roles = []string{rolesStr}
	}

	return payload, nil
}
