package user

import (
	"king-starter/internal/response"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type RegisterRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=6"`
	Email    string `json:"email" validate:"required,email"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type ResetCodeRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ResetPasswordRequest struct {
	Email       string `json:"email" validate:"required,email"`
	Code        string `json:"code" validate:"required,len=6"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

func (h *Handler) Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, 400, "invalid request")
	}

	user, err := h.service.Register(c.Request().Context(), req.Username, req.Password, req.Email)
	if err != nil {
		return response.Error(c, 500, err.Error())
	}

	return response.SuccessWithMsg(c, "register success", map[string]interface{}{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
	})
}

func (h *Handler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, 400, "invalid request")
	}

	token, err := h.service.Login(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		return response.Error(c, 401, err.Error())
	}

	return response.SuccessWithMsg(c, "login success", map[string]interface{}{
		"token": token,
	})
}

func (h *Handler) GenerateResetCode(c echo.Context) error {
	var req ResetCodeRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, 400, "invalid request")
	}

	code, err := h.service.GenerateResetCode(c.Request().Context(), req.Email)
	if err != nil {
		return response.Error(c, 500, err.Error())
	}

	return response.SuccessWithMsg(c, "reset code sent", map[string]interface{}{
		"code": code,
	})
}

func (h *Handler) ResetPassword(c echo.Context) error {
	var req ResetPasswordRequest
	if err := c.Bind(&req); err != nil {
		return response.Error(c, 400, "invalid request")
	}

	if err := h.service.ResetPassword(c.Request().Context(), req.Email, req.Code, req.NewPassword); err != nil {
		return response.Error(c, 500, err.Error())
	}

	return response.SuccessWithMsg(c, "password reset success", nil)
}

func (h *Handler) GetProfile(c echo.Context) error {
	//userID := c.Get("user_id").(uint)
	userID := c.QueryParam("user_id")

	var user User
	if err := h.service.GetDB().First(&user, userID).Error; err != nil {
		return response.Error(c, 404, "user not found")
	}

	return response.Success(c, map[string]interface{}{
		"user_id":  user.ID,
		"username": user.Username,
		"email":    user.Email,
		"phone":    user.Phone,
	})
}
