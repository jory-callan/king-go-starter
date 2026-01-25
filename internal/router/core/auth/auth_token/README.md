# Simple Auth Token Module

The `simple_auth_token` module serves as a lightweight "issuing and verification center" for authentication without external dependencies like Redis.

## Features

- **Lightweight Token Management**: Provides a simple system for issuing and validating authentication tokens
- **JWT Integration**: Uses `golang-jwt/v5` for secure token generation
- **No External Dependencies**: Does not require Redis or other external storage
- **Echo Middleware**: Built-in middleware for Echo framework integration
- **Token Refresh**: Ability to refresh tokens without re-authentication

## Differences from auth_token Module

- **No Redis Dependency**: Tokens are validated based on signature and expiration only
- **Simpler Architecture**: No session storage or revocation mechanism
- **Stateless Operation**: Each token validation is self-contained
- **Reduced Security**: Cannot revoke individual sessions; tokens remain valid until expiration

## Structure

### TokenPayload
Contains essential information stored in JWT tokens:
- `UserID` (int64): Unique identifier for the user
- `Username` (string): User's display name
- `Roles` ([]string): List of user roles
- `Source` (string): Authentication source
- `IssuedAt` (int64): Unix timestamp when token was issued
- `ExpiresAt` (int64): Unix timestamp when token expires

### Core Functions

#### Issue(payload TokenPayload) (string, error)
Generates a JWT token without external storage.

#### ValidateToken(tokenString string) (*TokenPayload, error)
Validates the token by checking its signature and expiration time.

#### RefreshToken(oldToken string) (string, error)
Creates a new token with the same payload but extended validity.

#### GetRemainingTime(tokenString string) (time.Duration, error)
Returns the remaining time before token expires.

#### JWTMiddleware() echo.MiddlewareFunc
Echo framework middleware that extracts and validates tokens from the Authorization header.

## Usage

```go
import "king-starter/internal/router/core/simple_auth_token"

// Initialize with JWT configuration (no Redis needed)
auth := simple_auth_token.New(jwtSecret, expireSeconds)

// Issue a new token
payload := simple_auth_token.TokenPayload{
    UserID:   123,
    Username: "john_doe",
    Roles:    []string{"user", "admin"},
    Source:   "login",
}
token, err := auth.Issue(payload)

// Use middleware in Echo routes
e.Use(auth.JWTMiddleware())

// Refresh a token
newToken, err := auth.RefreshToken(oldToken)
```

## Limitations

- **No Session Revocation**: Once issued, tokens remain valid until expiration
- **No Single Logout**: Cannot force logout of all sessions for a user
- **Reduced Security**: Cannot invalidate compromised tokens immediately