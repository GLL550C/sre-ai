package service

import (
	"bytes"
	"context"
	"core/model"
	"core/repository"
	"encoding/base64"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math/rand"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// JWT密钥 (生产环境应从环境变量或配置中心读取)
var jwtSecret = []byte("sre-ai-platform-jwt-secret-key-change-in-production")

// AuthService 认证服务
type AuthService struct {
	userRepo    *repository.UserRepository
	redisClient *redis.Client
	logger      *zap.Logger
}

// NewAuthService 创建认证服务
func NewAuthService(userRepo *repository.UserRepository, redisClient *redis.Client, logger *zap.Logger) *AuthService {
	return &AuthService{
		userRepo:    userRepo,
		redisClient: redisClient,
		logger:      logger,
	}
}

// GenerateCaptcha 生成验证码
func (s *AuthService) GenerateCaptcha() (*model.Captcha, error) {
	// 生成4位数字验证码
	code := fmt.Sprintf("%04d", rand.Intn(10000))
	captchaID := fmt.Sprintf("%d", rand.Int63())

	// 将验证码存储到Redis，5分钟过期
	ctx := context.Background()
	err := s.redisClient.Set(ctx, "captcha:"+captchaID, code, 5*time.Minute).Err()
	if err != nil {
		s.logger.Error("Failed to store captcha in redis", zap.Error(err))
		return nil, err
	}

	s.logger.Info("Captcha generated", zap.String("id", captchaID), zap.String("code", code))

	// 生成验证码图片（使用自定义验证码）
	var buf bytes.Buffer
	if err := s.writeCustomCaptchaImage(&buf, code, 120, 44); err != nil {
		return nil, err
	}

	// 转换为base64
	imageBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())

	return &model.Captcha{
		ID:        captchaID,
		Image:     "data:image/png;base64," + imageBase64,
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}, nil
}

// VerifyCaptcha 验证验证码
func (s *AuthService) VerifyCaptcha(id, code string) bool {
	ctx := context.Background()
	storedCode, err := s.redisClient.Get(ctx, "captcha:"+id).Result()
	if err != nil {
		s.logger.Warn("Captcha not found or expired", zap.String("id", id), zap.Error(err))
		return false
	}

	// 验证成功后删除验证码
	if storedCode == code {
		s.redisClient.Del(ctx, "captcha:"+id)
		return true
	}

	s.logger.Warn("Captcha verification failed", zap.String("id", id), zap.String("input", code), zap.String("expected", storedCode))
	return false
}

// writeCustomCaptchaImage 生成简单的数字验证码图片
func (s *AuthService) writeCustomCaptchaImage(w *bytes.Buffer, code string, width, height int) error {
	// 创建图片
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	// 填充背景色（浅灰色）
	bgColor := color.RGBA{240, 240, 240, 255}
	draw.Draw(img, img.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	// 添加干扰线
	for i := 0; i < 5; i++ {
		x1 := rand.Intn(width)
		y1 := rand.Intn(height)
		x2 := rand.Intn(width)
		y2 := rand.Intn(height)
		lineColor := color.RGBA{uint8(rand.Intn(200)), uint8(rand.Intn(200)), uint8(rand.Intn(200)), 255}
		s.drawLine(img, x1, y1, x2, y2, lineColor)
	}

	// 添加噪点
	for i := 0; i < 100; i++ {
		x := rand.Intn(width)
		y := rand.Intn(height)
		noiseColor := color.RGBA{uint8(rand.Intn(255)), uint8(rand.Intn(255)), uint8(rand.Intn(255)), 255}
		img.Set(x, y, noiseColor)
	}

	// 绘制数字（使用简单的块字体）
	digitWidth := width / len(code)
	for i, ch := range code {
		x := i*digitWidth + digitWidth/4
		y := height / 4
		digitColor := color.RGBA{uint8(rand.Intn(100)), uint8(rand.Intn(100)), uint8(rand.Intn(100) + 100), 255}
		s.drawDigit(img, x, y, digitWidth*2/3, height/2, int(ch-'0'), digitColor)
	}

	// 编码为PNG
	return png.Encode(w, img)
}

// drawLine 绘制线段
func (s *AuthService) drawLine(img *image.RGBA, x1, y1, x2, y2 int, c color.Color) {
	dx := abs(x2 - x1)
	dy := abs(y2 - y1)
	sx := 1
	if x1 > x2 {
		sx = -1
	}
	sy := 1
	if y1 > y2 {
		sy = -1
	}
	err := dx - dy

	for {
		if x1 >= 0 && x1 < img.Bounds().Dx() && y1 >= 0 && y1 < img.Bounds().Dy() {
			img.Set(x1, y1, c)
		}
		if x1 == x2 && y1 == y2 {
			break
		}
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x1 += sx
		}
		if e2 < dx {
			err += dx
			y1 += sy
		}
	}
}

// abs 绝对值
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// drawDigit 绘制数字
func (s *AuthService) drawDigit(img *image.RGBA, x, y, w, h, digit int, c color.Color) {
	// 简单的7段数码管显示
	segmentWidth := w / 5
	segmentHeight := h / 8

	// 数字对应的段（0-9）
	segments := [][]int{
		{1, 1, 1, 1, 1, 1, 0}, // 0
		{0, 1, 1, 0, 0, 0, 0}, // 1
		{1, 1, 0, 1, 1, 0, 1}, // 2
		{1, 1, 1, 1, 0, 0, 1}, // 3
		{0, 1, 1, 0, 0, 1, 1}, // 4
		{1, 0, 1, 1, 0, 1, 1}, // 5
		{1, 0, 1, 1, 1, 1, 1}, // 6
		{1, 1, 1, 0, 0, 0, 0}, // 7
		{1, 1, 1, 1, 1, 1, 1}, // 8
		{1, 1, 1, 1, 0, 1, 1}, // 9
	}

	seg := segments[digit]
	ox, oy := x, y

	// 绘制水平段
	if seg[0] == 1 {
		s.fillRect(img, ox+segmentWidth, oy, ox+w-segmentWidth, oy+segmentHeight, c)
	}
	if seg[6] == 1 {
		s.fillRect(img, ox+segmentWidth, oy+h/2-segmentHeight/2, ox+w-segmentWidth, oy+h/2+segmentHeight/2, c)
	}
	if seg[3] == 1 {
		s.fillRect(img, ox+segmentWidth, oy+h-segmentHeight, ox+w-segmentWidth, oy+h, c)
	}

	// 绘制垂直段
	if seg[1] == 1 {
		s.fillRect(img, ox+w-segmentWidth, oy+segmentHeight, ox+w, oy+h/2-segmentHeight/2, c)
	}
	if seg[2] == 1 {
		s.fillRect(img, ox+w-segmentWidth, oy+h/2+segmentHeight/2, ox+w, oy+h-segmentHeight, c)
	}
	if seg[5] == 1 {
		s.fillRect(img, ox, oy+segmentHeight, ox+segmentWidth, oy+h/2-segmentHeight/2, c)
	}
	if seg[4] == 1 {
		s.fillRect(img, ox, oy+h/2+segmentHeight/2, ox+segmentWidth, oy+h-segmentHeight, c)
	}
}

// fillRect 填充矩形
func (s *AuthService) fillRect(img *image.RGBA, x1, y1, x2, y2 int, c color.Color) {
	for x := x1; x < x2 && x < img.Bounds().Dx(); x++ {
		for y := y1; y < y2 && y < img.Bounds().Dy(); y++ {
			if x >= 0 && y >= 0 {
				img.Set(x, y, c)
			}
		}
	}
}

// Login 用户登录
func (s *AuthService) Login(req *model.LoginRequest, clientIP string) (*model.LoginResponse, error) {
	// 验证验证码
	if !s.VerifyCaptcha(req.CaptchaID, req.Captcha) {
		return nil, fmt.Errorf("验证码错误或已过期")
	}

	// 验证用户
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		s.logger.Error("Failed to get user", zap.Error(err))
		return nil, fmt.Errorf("用户名或密码错误")
	}
	if user == nil {
		return nil, fmt.Errorf("用户名或密码错误")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		s.logger.Warn("Invalid password attempt", zap.String("username", req.Username))
		return nil, fmt.Errorf("用户名或密码错误")
	}

	// 更新登录信息
	if err := s.userRepo.UpdateLoginInfo(user.ID, clientIP); err != nil {
		s.logger.Error("Failed to update login info", zap.Error(err))
	}

	// 生成JWT token
	token, expiresAt, err := s.GenerateToken(user)
	if err != nil {
		s.logger.Error("Failed to generate token", zap.Error(err))
		return nil, fmt.Errorf("生成令牌失败")
	}

	// 清除密码
	user.Password = ""

	return &model.LoginResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		User:      *user,
	}, nil
}

// GenerateToken 生成JWT token
func (s *AuthService) GenerateToken(user *model.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(24 * time.Hour)

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      expiresAt.Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// ParseToken 解析JWT token
func (s *AuthService) ParseToken(tokenString string) (*model.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecret, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID, _ := claims["user_id"].(float64)
		username, _ := claims["username"].(string)
		role, _ := claims["role"].(string)

		return &model.User{
			ID:       int64(userID),
			Username: username,
			Role:     role,
		}, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// HashPassword 密码加密
func (s *AuthService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword 验证密码
func (s *AuthService) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// ChangePassword 修改密码
func (s *AuthService) ChangePassword(userID int64, req *model.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("用户不存在")
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.OldPassword)); err != nil {
		return fmt.Errorf("原密码错误")
	}

	// 加密新密码
	hashedPassword, err := s.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}

	return s.userRepo.UpdatePassword(userID, hashedPassword)
}
