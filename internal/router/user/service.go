package user

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"king-starter/pkg/jwt"
	"math/rand"
	"strconv"
)

type Service struct {
	db        *gorm.DB
	jwt       *jwt.JWT
	jwtExpire int
}

func NewService(db *gorm.DB, jwt *jwt.JWT, jwtExpire int) *Service {
	return &Service{
		db:        db,
		jwt:       jwt,
		jwtExpire: jwtExpire,
	}
}

// GetDB returns the database instance
func (s *Service) GetDB() *gorm.DB {
	return s.db
}

func (s *Service) hashPassword(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}

func (s *Service) Register(ctx context.Context, username, password, email string) (*User, error) {
	var existingUser User
	if err := s.db.Where("username = ? OR email = ?", username, email).First(&existingUser).Error; err == nil {
		return nil, errors.New("username or email already exists")
	}

	user := User{
		Username: username,
		Password: s.hashPassword(password),
		Email:    email,
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *Service) Login(ctx context.Context, username, password string) (string, error) {
	var user User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("user not found")
		}
		return "", err
	}

	if user.Password != s.hashPassword(password) {
		return "", errors.New("invalid password")
	}

	token, err := s.jwt.GenerateToken(strconv.Itoa(int(user.ID)), user.Username, strconv.Itoa(s.jwtExpire))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) GenerateResetCode(ctx context.Context, email string) (string, error) {
	var user User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", errors.New("email not found")
		}
		return "", err
	}

	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	return code, nil
}

func (s *Service) ResetPassword(ctx context.Context, email, code, newPassword string) error {
	var user User
	if err := s.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("email not found")
		}
		return err
	}

	user.Password = s.hashPassword(newPassword)
	if err := s.db.Save(&user).Error; err != nil {
		return err
	}

	return nil
}
