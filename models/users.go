package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	// Import postgres driver
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var (
	// ErrorNotFound is default record not found error
	ErrorNotFound = errors.New("models: resource not found")
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
	err := us.db.Where("id = ?", id).First(&user).Error
	switch err {
	case nil:
		return &user, nil
	case gorm.ErrRecordNotFound:
		return nil, ErrorNotFound
	default:
		return nil, err
	}
}

// Create new user
func (us *UserService) Create(user *User) error {
	return us.db.Create(user).Error
}

// Close database connection
func (us *UserService) Close() error {
	return us.db.Close()
}

// DestructiveReset drop all tables and recreate database
func (us *UserService) DestructiveReset() {
	us.db.DropTableIfExists(&User{})
	us.db.AutoMigrate(&User{})
}

// User model
type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"not null;unique_index"`
}
