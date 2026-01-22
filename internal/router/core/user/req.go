package user

// CreateUserReq 创建用户请求
type CreateUserReq struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=6"`
	Nickname string `json:"nickname"`
	Email    string `json:"email" validate:"omitempty,email"`
	Phone    string `json:"phone"`
}

// UpdateUserReq 更新用户请求
type UpdateUserReq struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email" validate:"omitempty,email"`
	Phone    string `json:"phone"`
}

// ChangePasswordReq 修改密码请求
type ChangePasswordReq struct {
	OldPassword string `json:"old_password" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}
