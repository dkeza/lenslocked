package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
	// Import postgres driver
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	// ErrorNotFound is default record not found error
	ErrorNotFound = errors.New("models: resource not found")
	// ErrorInvalidID is invalid id
	ErrorInvalidID = errors.New("models: invalid ID")
)

// NewUserService creates new user service
func NewUserService(connectionInfo string) (*UserService, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)

	return &UserService{
		db: db,
	}, nil
}

// UserService to access users
type UserService struct {
	db *gorm.DB
}

// ByID finds user
func (us *UserService) ByID(id uint) (*User, error) {
	var user User
	db := us.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

// ByEmail finds user by email
func (us *UserService) ByEmail(email string) (*User, error) {
	var user User
	db := us.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

func first(db *gorm.DB, dst interface{}) error {
	err := db.First(dst).Error
	if err == gorm.ErrRecordNotFound {
		return ErrorNotFound
	}
	return err
}

// Create new user
func (us *UserService) Create(user *User) error {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""
	return us.db.Create(user).Error
}

// Update user
func (us *UserService) Update(user *User) error {
	return us.db.Save(user).Error
}

// Delete user
func (us *UserService) Delete(id uint) error {
	if id == 0 {
		return ErrorInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return us.db.Delete(&user).Error
}

// Close database connection
func (us *UserService) Close() error {
	return us.db.Close()
}

// DestructiveReset drop all tables and recreate database
func (us *UserService) DestructiveReset() error {
	if err := us.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return us.AutoMigrate()
}

// AutoMigrate atempt to migrate users table
func (us *UserService) AutoMigrate() error {
	if err := us.db.AutoMigrate(&User{}).Error; err != nil {
		return err
	}
	return nil
}

// User model
type User struct {
	gorm.Model
	Name         string
	Email        string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
}
