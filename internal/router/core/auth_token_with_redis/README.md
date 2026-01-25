# Auth Token Module

The `auth_token` module serves as a unified "issuing and verification center" for all authentication methods in the system (such as phone login, account login, etc.).

## Features

- **Unified Token Management**: Provides a centralized system for issuing and validating authentication tokens
- **JWT Integration**: Uses `golang-jwt/v5` for secure token generation with `jti` field
- **Redis-backed Sessions**: Stores session information in Redis for enhanced security
- **Single Sign-On (SSO)**: Supports single point of login with automatic logout of previous sessions
- **Session Revocation**: Ability to force logout by revoking user sessions
- **Echo Middleware**: Built-in middleware for Echo framework integration

## Structure

### TokenPayload
Contains essential information stored in JWT tokens:
- `UserID` (int64): Unique identifier for the user
- `Username` (string): User's display name
- `Roles` ([]string): List of user roles
- `Source` (string): Authentication source

### Core Functions

#### Issue(payload TokenPayload) (string, error)
Generates a JWT token and stores the session in Redis with key `user_session:{UserID}`.

#### Revoke(userID int64) error
Removes the user session from Redis, effectively logging out the user.

#### ValidateToken(tokenString string) (*TokenPayload, error)
Validates the token by checking its signature and verifying it hasn't been revoked by comparing the `jti` with the one stored in Redis.

#### JWTMiddleware() echo.MiddlewareFunc
Echo framework middleware that extracts and validates tokens from the Authorization header.

## Security Features

- **JTI Verification**: Each token has a unique identifier that's verified against Redis to prevent unauthorized sessions
- **Session Revocation**: When a new token is issued for a user, the old session is automatically invalidated
- **Token Expiration**: Tokens have configurable expiration times

## Usage

```go
import "king-starter/internal/router/core/auth_token"

// Initialize with Redis client and JWT configuration
auth := auth_token.New(redisClient, jwtSecret, expireSeconds)

// Issue a new token
payload := auth_token.TokenPayload{
    UserID:   123,
    Username: "john_doe",
    Roles:    []string{"user", "admin"},
    Source:   "login",
}
token, err := auth.Issue(payload)

// Use middleware in Echo routes
e.Use(auth.JWTMiddleware())

// Revoke a user's session
auth.Revoke(userID)
```