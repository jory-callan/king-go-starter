package identity

// RegisterReq 注册请求
type RegisterReq struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=6,max=20"`
	Nickname string `json:"nickname" validate:"omitempty,max=50"`
	Email    string `json:"email" validate:"omitempty,email"`
	Phone    string `json:"phone" validate:"omitempty"`
}