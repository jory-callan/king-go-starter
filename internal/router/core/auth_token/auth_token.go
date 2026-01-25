package simple_auth_token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// TokenPayload contains the essential information stored in a JWT token
type TokenPayload struct {
	UserID    int64    `json:"user_id"`
	Username  string   `json:"username"`
	Roles     []string `json:"roles"`
	Source    string   `json:"source"`
	IssuedAt  int64    `json:"iat"`
	ExpiresAt int64    `json:"exp"`
}

// SimpleAuthToken represents the authentication token manager without Redis
type SimpleAuthToken struct {
	jwtSecret   []byte
	tokenExpire time.Duration
}

// New creates a new SimpleAuthToken instance without Redis dependency
func New(jwtSecret string, expireSeconds int) *SimpleAuthToken {
	return &SimpleAuthToken{
		jwtSecret:   []byte(jwtSecret),
		tokenExpire: time.Duration(expireSeconds) * time.Second,
	}
}

// Issue generates a JWT token without storing anything in Redis
func (sat *SimpleAuthToken) Issue(payload TokenPayload) (string, error) {
	// Create claims with the payload
	now := time.Now()
	claims := jwt.MapClaims{
		"user_id":  payload.UserID,
		"username": payload.Username,
		"roles":    payload.Roles,
		"source":   payload.Source,
		"iat":      now.Unix(),
		"exp":      now.Add(sat.tokenExpire).Unix(),
		"nbf":      now.Unix(),
	}

	// Create the token with HS256 signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret
	signedToken, err := token.SignedString(sat.jwtSecret)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

// ValidateToken checks if the token is valid based only on signature and expiration
func (sat *SimpleAuthToken) ValidateToken(tokenString string) (*TokenPayload, error) {
	// Parse the token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return sat.jwtSecret, nil
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

	// Extract basic information from claims
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid user_id in token")
	}
	userID := int64(userIDFloat)

	username, ok := claims["username"].(string)
	if !ok {
		return nil, errors.New("invalid username in token")
	}

	source, ok := claims["source"].(string)
	if !ok {
		return nil, errors.New("invalid source in token")
	}

	iat, ok := claims["iat"].(float64)
	if !ok {
		return nil, errors.New("invalid issued at time in token")
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, errors.New("invalid expiration time in token")
	}

	// Build and return the TokenPayload
	payload := &TokenPayload{
		UserID:    userID,
		Username:  username,
		Source:    source,
		IssuedAt:  int64(iat),
		ExpiresAt: int64(exp),
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

// RefreshToken creates a new token with the same payload but extended validity
func (sat *SimpleAuthToken) RefreshToken(oldToken string) (string, error) {
	// First validate the old token
	oldPayload, err := sat.ValidateToken(oldToken)
	if err != nil {
		return "", fmt.Errorf("cannot refresh invalid token: %w", err)
	}

	// Create a new token with the same payload
	newToken, err := sat.Issue(*oldPayload)
	if err != nil {
		return "", fmt.Errorf("failed to issue refreshed token: %w", err)
	}

	return newToken, nil
}

// GetRemainingTime returns the remaining time before token expires
func (sat *SimpleAuthToken) GetRemainingTime(tokenString string) (time.Duration, error) {
	payload, err := sat.ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}

	expirationTime := time.Unix(payload.ExpiresAt, 0)
	remaining := time.Until(expirationTime)

	if remaining < 0 {
		return 0, errors.New("token has already expired")
	}

	return remaining, nil
}
