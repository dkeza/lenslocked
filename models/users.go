package models

import (
	"errors"
	"lenslocked/hash"
	"lenslocked/rand"

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
	// ErrorInvalidPassword error
	ErrorInvalidPassword = errors.New("models: incorrect password provided")
)

const usPasswordPepper = "some-secret"
const hmacSecretKey = "hmac-secret"

type UserDB interface {
	// Methods for querying
	ByID(id uint) (*User, error)
	ByEmail(email string) (*User, error)
	ByRemember(token string) (*User, error)

	// Methods for altering users
	Create(user *User) error
	Update(user *User) error
	Delete(id uint) error

	// Clode DB connection
	Close() error

	// Migration helpers
	AutoMigrate() error
	DestructiveReset() error
}

// NewUserService creates new user service
func NewUserService(connectionInfo string) (*UserService, error) {
	ug, err := newUserGorm(connectionInfo)
	if err != nil {
		return nil, err
	}
	return &UserService{
		UserDB: &userValidator{
			UserDB: ug,
		},
	}, nil
}

// UserService to access users
type UserService struct {
	UserDB
}

type userValidator struct {
	UserDB
}

// newUserGorm creates new user service
func newUserGorm(connectionInfo string) (*userGorm, error) {
	db, err := gorm.Open("postgres", connectionInfo)
	if err != nil {
		return nil, err
	}
	db.LogMode(true)
	hmac := hash.NewHMAC(hmacSecretKey)
	return &userGorm{
		db:   db,
		hmac: hmac,
	}, nil
}

var _ UserDB = &userGorm{}

type userGorm struct {
	db   *gorm.DB
	hmac hash.HMAC
}

// ByID finds user
func (ug *userGorm) ByID(id uint) (*User, error) {
	var user User
	db := ug.db.Where("id = ?", id)
	err := first(db, &user)
	return &user, err
}

// ByEmail finds user by email
func (ug *userGorm) ByEmail(email string) (*User, error) {
	var user User
	db := ug.db.Where("email = ?", email)
	err := first(db, &user)
	return &user, err
}

// Authenticate authenticates user
func (us *UserService) Authenticate(email, password string) (*User, error) {
	userFound, err := us.ByEmail(email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(userFound.PasswordHash), []byte(password+usPasswordPepper))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return nil, ErrorInvalidPassword
	} else if err != nil {
		return nil, err
	}
	return userFound, nil
}

// ByRemember finds user by token
func (ug *userGorm) ByRemember(token string) (*User, error) {
	var user User
	rememberHash := ug.hmac.Hash(token)
	err := first(ug.db.Where("remember_hash = ?", rememberHash), &user)
	if err != nil {
		return nil, err
	}
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
func (ug *userGorm) Create(user *User) error {
	hashedPassword := []byte(user.Password + usPasswordPepper)
	hashedBytes, err := bcrypt.GenerateFromPassword(hashedPassword, bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedBytes)
	user.Password = ""

	if user.Remember == "" {
		token, err := rand.RememberToken()
		if err != nil {
			return err
		}
		user.Remember = token
	}
	user.RememberHash = ug.hmac.Hash(user.Remember)
	return ug.db.Create(user).Error
}

// Update user
func (ug *userGorm) Update(user *User) error {
	if user.Remember != "" {
		user.RememberHash = ug.hmac.Hash(user.Remember)
	}
	return ug.db.Save(user).Error
}

// Delete user
func (ug *userGorm) Delete(id uint) error {
	if id == 0 {
		return ErrorInvalidID
	}
	user := User{Model: gorm.Model{ID: id}}
	return ug.db.Delete(&user).Error
}

// Close database connection
func (ug *userGorm) Close() error {
	return ug.db.Close()
}

// DestructiveReset drop all tables and recreate database
func (ug *userGorm) DestructiveReset() error {
	if err := ug.db.DropTableIfExists(&User{}).Error; err != nil {
		return err
	}
	return ug.AutoMigrate()
}

// AutoMigrate atempt to migrate users table
func (ug *userGorm) AutoMigrate() error {
	if err := ug.db.AutoMigrate(&User{}).Error; err != nil {
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
	Remember     string `gorm:"-"`
	RememberHash string `gorm:"not null;unique_index"`
}
