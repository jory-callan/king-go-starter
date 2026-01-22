package identity

// LoginReq 登录请求参数
type LoginReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
	Captcha  string `json:"captcha,omitempty"`
	Remember bool   `json:"remember,omitempty"`
}

// LogoutReq 登出请求参数
type LogoutReq struct {
	RefreshToken string `json:"refresh_token,omitempty"`
}

// RefreshTokenReq 刷新令牌请求参数
type RefreshTokenReq struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}
