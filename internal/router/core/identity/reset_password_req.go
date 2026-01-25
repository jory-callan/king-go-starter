package identity

// ResetPasswordReq 重置密码请求
type ResetPasswordReq struct {
	UserID      string `json:"user_id" validate:"required"`
	OldPassword string `json:"old_password" validate:"required,min=6"`
	NewPassword string `json:"new_password" validate:"required,min=6,max=20"`
}