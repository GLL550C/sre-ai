package service

import (
	"core/model"
	"core/repository"
	"fmt"

	"go.uber.org/zap"
)

// UserService 用户服务
type UserService struct {
	userRepo    *repository.UserRepository
	authService *AuthService
	logger      *zap.Logger
}

// NewUserService 创建用户服务
func NewUserService(userRepo *repository.UserRepository, authService *AuthService, logger *zap.Logger) *UserService {
	return &UserService{
		userRepo:    userRepo,
		authService: authService,
		logger:      logger,
	}
}

// ListUsers 获取用户列表
func (s *UserService) ListUsers(page, pageSize int) ([]model.User, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize
	return s.userRepo.List(offset, pageSize)
}

// GetUser 获取用户详情
func (s *UserService) GetUser(id int64) (*model.User, error) {
	return s.userRepo.GetByID(id)
}

// CreateUser 创建用户
func (s *UserService) CreateUser(req *model.CreateUserRequest, createdBy string) (*model.User, error) {
	// 检查用户名是否已存在
	exists, err := s.userRepo.CheckUsernameExists(req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("username already exists")
	}

	// 加密密码
	hashedPassword, err := s.authService.HashPassword(req.Password)
	if err != nil {
		s.logger.Error("Failed to hash password", zap.Error(err))
		return nil, fmt.Errorf("failed to process password")
	}

	// 验证角色
	validRoles := map[string]bool{"admin": true, "operator": true, "viewer": true}
	if !validRoles[req.Role] {
		req.Role = "viewer"
	}

	user := &model.User{
		Username: req.Username,
		Password: hashedPassword,
		Email:    req.Email,
		Phone:    req.Phone,
		Role:     req.Role,
		Status:   1,
	}

	if err := s.userRepo.Create(user); err != nil {
		s.logger.Error("Failed to create user", zap.Error(err))
		return nil, fmt.Errorf("failed to create user")
	}

	// 清除密码
	user.Password = ""
	return user, nil
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(id int64, email, phone, role string, status int) (*model.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// 验证角色
	validRoles := map[string]bool{"admin": true, "operator": true, "viewer": true}
	if validRoles[role] {
		user.Role = role
	}

	user.Email = email
	user.Phone = phone
	user.Status = status

	if err := s.userRepo.Update(user); err != nil {
		s.logger.Error("Failed to update user", zap.Error(err))
		return nil, fmt.Errorf("failed to update user")
	}

	user.Password = ""
	return user, nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id int64) error {
	// 检查是否是最后一个管理员
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	if user.Role == "admin" {
		adminCount, err := s.userRepo.CountAdmins()
		if err != nil {
			return err
		}
		if adminCount <= 1 {
			return fmt.Errorf("cannot delete the last admin")
		}
	}

	return s.userRepo.Delete(id)
}

// ResetPassword 重置用户密码
func (s *UserService) ResetPassword(id int64, newPassword string) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user not found")
	}

	hashedPassword, err := s.authService.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(id, hashedPassword)
}

// GetCurrentUser 获取当前用户信息
func (s *UserService) GetCurrentUser(userID int64) (*model.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	user.Password = ""
	return user, nil
}

// ValidateRole 验证角色权限
func (s *UserService) ValidateRole(role string) bool {
	validRoles := map[string]bool{"admin": true, "operator": true, "viewer": true}
	return validRoles[role]
}

// HasPermission 检查用户是否有权限
func HasPermission(userRole, requiredRole string) bool {
	roleLevels := map[string]int{
		"viewer":   1,
		"operator": 2,
		"admin":    3,
	}

	userLevel := roleLevels[userRole]
	requiredLevel := roleLevels[requiredRole]

	return userLevel >= requiredLevel
}
