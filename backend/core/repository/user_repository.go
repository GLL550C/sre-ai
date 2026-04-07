package repository

import (
	"core/model"
	"database/sql"
	"fmt"

	"go.uber.org/zap"
)

// UserRepository 用户仓库
type UserRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewUserRepository 创建用户仓库
func NewUserRepository(db *sql.DB, logger *zap.Logger) *UserRepository {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

// GetByUsername 根据用户名获取用户
func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	query := `SELECT id, username, password, email, phone, role, status,
		last_login_at, last_login_ip, created_at, updated_at
		FROM users WHERE username = ? AND status = 1`

	var user model.User
	var lastLoginAt sql.NullTime
	var lastLoginIP sql.NullString
	err := r.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Password, &user.Email, &user.Phone,
		&user.Role, &user.Status, &lastLoginAt, &lastLoginIP,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}
	if lastLoginIP.Valid {
		user.LastLoginIP = lastLoginIP.String
	}
	return &user, nil
}

// GetByID 根据ID获取用户
func (r *UserRepository) GetByID(id int64) (*model.User, error) {
	query := `SELECT id, username, password, email, phone, role, status,
		last_login_at, last_login_ip, created_at, updated_at
		FROM users WHERE id = ?`

	var user model.User
	var lastLoginAt sql.NullTime
	var lastLoginIP sql.NullString
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Password, &user.Email, &user.Phone,
		&user.Role, &user.Status, &lastLoginAt, &lastLoginIP,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	if lastLoginAt.Valid {
		user.LastLoginAt = &lastLoginAt.Time
	}
	if lastLoginIP.Valid {
		user.LastLoginIP = lastLoginIP.String
	}
	return &user, nil
}

// List 获取用户列表
func (r *UserRepository) List(offset, limit int) ([]model.User, int64, error) {
	// 获取总数
	var total int64
	err := r.db.QueryRow("SELECT COUNT(*) FROM users WHERE status = 1").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, username, email, phone, role, status,
		last_login_at, last_login_ip, created_at, updated_at
		FROM users WHERE status = 1
		ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		var lastLoginAt sql.NullTime
		var lastLoginIP sql.NullString
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.Phone,
			&user.Role, &user.Status, &lastLoginAt, &lastLoginIP,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan user", zap.Error(err))
			continue
		}
		if lastLoginAt.Valid {
			user.LastLoginAt = &lastLoginAt.Time
		}
		if lastLoginIP.Valid {
			user.LastLoginIP = lastLoginIP.String
		}
		users = append(users, user)
	}

	return users, total, nil
}

// Create 创建用户
func (r *UserRepository) Create(user *model.User) error {
	query := `INSERT INTO users (username, password, email, phone, role, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())`

	result, err := r.db.Exec(query, user.Username, user.Password, user.Email, user.Phone, user.Role, user.Status)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = id
	return nil
}

// Update 更新用户
func (r *UserRepository) Update(user *model.User) error {
	query := `UPDATE users SET email = ?, phone = ?, role = ?, status = ?, updated_at = NOW() WHERE id = ?`
	_, err := r.db.Exec(query, user.Email, user.Phone, user.Role, user.Status, user.ID)
	return err
}

// UpdatePassword 更新密码
func (r *UserRepository) UpdatePassword(userID int64, password string) error {
	query := `UPDATE users SET password = ?, updated_at = NOW() WHERE id = ?`
	_, err := r.db.Exec(query, password, userID)
	return err
}

// UpdateLoginInfo 更新登录信息
func (r *UserRepository) UpdateLoginInfo(userID int64, ip string) error {
	query := `UPDATE users SET last_login_at = NOW(), last_login_ip = ?, updated_at = NOW() WHERE id = ?`
	_, err := r.db.Exec(query, ip, userID)
	return err
}

// Delete 删除用户(软删除)
func (r *UserRepository) Delete(id int64) error {
	query := `UPDATE users SET status = 0, updated_at = NOW() WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

// CheckUsernameExists 检查用户名是否存在
func (r *UserRepository) CheckUsernameExists(username string) (bool, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM users WHERE username = ? AND status = 1", username).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// CountAdmins 统计管理员数量
func (r *UserRepository) CountAdmins() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM users WHERE role = 'admin' AND status = 1").Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

// ValidateUserCredentials 验证用户凭据
func (r *UserRepository) ValidateUserCredentials(username, password string) (*model.User, error) {
	user, err := r.GetByUsername(username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}
	return user, nil
}
