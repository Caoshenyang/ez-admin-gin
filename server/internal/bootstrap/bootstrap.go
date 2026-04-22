package bootstrap

import (
	"errors"
	"fmt"

	"ez-admin-gin/server/internal/model"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const (
	defaultAdminUsername = "admin"
	defaultAdminPassword = "EzAdmin@123456"
)

// Run 执行服务启动时必须完成的初始化动作。
func Run(db *gorm.DB, log *zap.Logger) error {
	if err := seedDefaultAdmin(db, log); err != nil {
		return fmt.Errorf("seed default admin: %w", err)
	}

	return nil
}

// seedDefaultAdmin 创建本地起步用的默认管理员。
func seedDefaultAdmin(db *gorm.DB, log *zap.Logger) error {
	var user model.User
	// Unscoped 会把已逻辑删除记录也查出来，避免重复创建同名默认账号。
	err := db.Unscoped().Where("username = ?", defaultAdminUsername).First(&user).Error
	if err == nil {
		return nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(defaultAdminPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash default admin password: %w", err)
	}

	user = model.User{
		Username:     defaultAdminUsername,
		PasswordHash: string(passwordHash),
		Nickname:     "系统管理员",
		Status:       model.UserStatusEnabled,
	}

	if err := db.Create(&user).Error; err != nil {
		return err
	}

	log.Info("default admin user created", zap.String("username", defaultAdminUsername))

	return nil
}
