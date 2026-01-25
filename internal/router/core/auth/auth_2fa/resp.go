package auth_2fa

// EnableTwoFAReq 启用 2FA 请求参数
type EnableTwoFAReq struct {
	UserID string `json:"user_id" validate:"required"`
	Code   string `json:"code" validate:"required,len=6"`
}

// VerifyTwoFAReq 验证 2FA 请求参数
type VerifyTwoFAReq struct {
	UserID string `json:"user_id" validate:"required"`
	Code   string `json:"code" validate:"required,len=6"`
}

// DisableTwoFAReq 禁用 2FA 请求参数
type DisableTwoFAReq struct {
	UserID string `json:"user_id" validate:"required"`
	Code   string `json:"code" validate:"required,len=6"`
}
