package models

import (
	"fmt"
	"testing"
	"time"
)

func testingUserService() (*UserService, error) {
	const (
		host     = "localhost"
		port     = 5432
		user     = "postgres"
		password = "1234"
		dbname   = "lenslocked_test"
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	us, err := NewUserService(psqlInfo)
	if err != nil {
		return nil, err
	}
	us.db.LogMode(false)

	us.DestructiveReset()
	return us, nil
}

func TestCreateUser(t *testing.T) {
	us, err := testingUserService()
	if err != nil {
		t.Fatal(err)
	}
	u := User{
		Name:  "Mika",
		Email: "mika@email.com",
	}
	err = us.Create(&u)
	if err != nil {
		t.Fatal(err)
	}
	if u.ID == 0 {
		t.Errorf("Expected ID > 0. Received %d", u.ID)
	}
	if time.Since(u.CreatedAt) > time.Duration(5*time.Second) {
		t.Errorf("Expected CreatedAt to be recent. Received %v", u.CreatedAt)
	}
	if time.Since(u.UpdatedAt) > time.Duration(5*time.Second) {
		t.Errorf("Expected UpdatedAt to be recent. Received %v", u.CreatedAt)
	}
}
