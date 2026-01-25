package identity

import (
	"fmt"
	"king-starter/internal/app"
	"king-starter/internal/response"
	"king-starter/internal/router/core/user"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type RegisterHandler struct {
	app          *app.App
	registerRepo *RegisterRepo
	userRepo     *user.Repository
}

func NewRegisterHandler(app *app.App) *RegisterHandler {
	return &RegisterHandler{
		app:          app,
		registerRepo: NewRegisterRepo(app.Db.DB),
		userRepo:     user.NewRepository(app.Db.DB),
	}
}

// Register 用户注册
func (h *RegisterHandler) Register(c echo.Context) error {
	var req RegisterReq
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	if err := c.Validate(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	// 检查用户名是否已存在
	existingUser, _ := h.userRepo.GetByUsername(c.Request().Context(), req.Username)
	if existingUser != nil {
		return response.Error(c, http.StatusBadRequest, "用户名已存在")
	}

	// 密码加密
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, fmt.Sprintf("密码加密失败: %v", err))
	}

	// 生成用户ID
	userID := uuid.NewString()

	// 创建用户
	newUser := &user.CoreUser{
		ID:        userID,
		Username:  req.Username,
		Password:  string(passwordHash),
		Nickname:  req.Nickname,
		Email:     req.Email,
		Phone:     req.Phone,
		Status:    1,        // 默认启用状态
		CreatedBy: "system", // 注册时默认创建者
		UpdatedBy: "system",
	}

	if err := h.userRepo.Create(c.Request().Context(), newUser); err != nil {
		return response.Error(c, http.StatusInternalServerError, fmt.Sprintf("用户创建失败: %v", err))
	}

	// 为新用户分配默认角色（例如：普通用户角色）
	defaultRoleID := "role-user" // 可以根据实际需求调整
	if err := h.registerRepo.AssignRoleToUser(c.Request().Context(), newUser.ID, defaultRoleID, "system"); err != nil {
		// 如果分配角色失败，记录日志但不阻止注册成功
		fmt.Printf("分配角色失败: %v\n", err)
	}

	resp := &RegisterResp{
		ID:       newUser.ID,
		Username: newUser.Username,
		Nickname: newUser.Nickname,
		Email:    newUser.Email,
		Phone:    newUser.Phone,
	}

	return response.Success(c, resp)
}

// ResetPassword 重置密码
func (h *RegisterHandler) ResetPassword(c echo.Context) error {
	var req ResetPasswordReq
	if err := c.Bind(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, "请求参数错误")
	}

	if err := c.Validate(&req); err != nil {
		return response.Error(c, http.StatusBadRequest, err.Error())
	}

	// 获取当前用户信息
	currentUser, err := h.userRepo.GetByID(c.Request().Context(), req.UserID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return response.Error(c, http.StatusNotFound, "用户不存在")
		}
		return response.Error(c, http.StatusInternalServerError, fmt.Sprintf("查询用户失败: %v", err))
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(currentUser.Password), []byte(req.OldPassword)); err != nil {
		return response.Error(c, http.StatusBadRequest, "原密码错误")
	}

	// 新密码不能与旧密码相同
	if req.OldPassword == req.NewPassword {
		return response.Error(c, http.StatusBadRequest, "新密码不能与原密码相同")
	}

	// 加密新密码
	newPasswordHash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return response.Error(c, http.StatusInternalServerError, fmt.Sprintf("新密码加密失败: %v", err))
	}

	// 更新密码
	currentUser.Password = string(newPasswordHash)
	if err := h.userRepo.Update(c.Request().Context(), currentUser); err != nil {
		return response.Error(c, http.StatusInternalServerError, fmt.Sprintf("更新密码失败: %v", err))
	}

	return response.SuccessWithMsg[any](c, "密码重置成功", nil)
}
