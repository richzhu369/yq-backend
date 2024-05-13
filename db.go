package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	ID            int       `gorm:"primaryKey"`
	Username      string    `json:"username" gorm:"size:64;index:idx_username,unique" binding:"required"`
	Password      string    `json:"password" gorm:"-" binding:"required"`
	PasswordHash  string    `gorm:"size:128"`
	LastLoginTime time.Time `gorm:"default:null"`
	Role          string    `json:"role" gorm:"size:20;not null"`
}

type WhitelistLog struct {
	gorm.Model
	IP     string `json:"ip"`
	OpUser string `json:"opUser"`
}

func (u *User) SetPassword(password string) string {

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err.Error()
	}

	return string(hashedPassword)

}

func (u *User) CheckPassword(password, hashedPassword string) bool {

	// 比较密码哈希
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil

}

func (u *User) CreateUser(db *gorm.DB, user User) error {
	var count int64
	db.Model(&User{}).Where("username = ?", user.Username).Count(&count)
	if count > 0 {
		return fmt.Errorf("user already exists")
	}

	if err := db.Create(&user).Error; err != nil {
		return err
	}

	return nil
}

func (u *User) DeleteUser(db *gorm.DB, user User) error {
	var count int64

	db.Model(&User{}).Where("username = ?", user.Username).Count(&count)
	if count < 1 {
		return fmt.Errorf("user 不存在")
	}

	if err := db.Delete(&user).Error; err != nil {
		return err
	}

	return nil
}
