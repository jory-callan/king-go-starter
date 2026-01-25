package simple_auth_token

// Example usage of the simple_auth_token module
func ExampleUsage() {
	// Create simple auth token instance with JWT secret and expiration time
	auth := New("your-jwt-secret-key", 3600) // 1 hour expiration

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

	// The token can be validated using ValidateToken
	validPayload, err := auth.ValidateToken(token)
	if err != nil {
		// token is invalid or has been expired
		return
	}

	// validPayload now contains the user information
	_ = validPayload

	// Refresh a token
	_, err = auth.RefreshToken(token)
	if err != nil {
		// handle error
		return
	}

	// Check remaining time before expiration
	// remaining, err := auth.GetRemainingTime(token)
	// if err != nil {
	// 	// token is expired or invalid
	// 	return
	// }
}
