package auth_token

import (
	"github.com/redis/go-redis/v9"
)

// Example usage of the auth_token module
func ExampleUsage() {
	// Initialize Redis client
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		DB:   0,
	})

	// Create auth token instance with JWT secret and expiration time
	auth := New(rdb, "your-jwt-secret-key", 3600) // 1 hour expiration

	// Create token payload
	payload := TokenPayload{
		UserID:   12345,
		Username: "test_user",
		Roles:    []string{"user", "admin"},
		Source:   "account_login",
	}

	// Issue a new token
	token, err := auth.Issue(payload)
	if err != nil {
		// handle error
		return
	}

	// The token is now stored in Redis with key "user_session:12345"
	// and can be validated using ValidateToken

	// To validate a token
	validPayload, err := auth.ValidateToken(token)
	if err != nil {
		// token is invalid or has been revoked
		return
	}

	// validPayload now contains the user information
	_ = validPayload

	// To revoke a user's session (force logout)
	err = auth.Revoke(12345)
	if err != nil {
		// handle error
		return
	}
}
