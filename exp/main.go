package main

import (
	"fmt"
	"simplegallery/models"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1234"
	dbname   = "simplegallery_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()

	us.DestructiveReset()

	u := models.User{
		Name:  "Mika",
		Email: "mika@email.com",
	}

	if err := us.Create(&u); err != nil {
		panic(err)
	}

	u.Email = "mika.mikic@google.com"
	if err := us.Update(&u); err != nil {
		panic(err)
	}

	user, err := us.ByID(u.ID)
	if err != nil {
		panic(err)
	}
	fmt.Println(user)

	err = us.Delete(u.ID)
	if err != nil {
		panic(err)
	}

	user, err = us.ByEmail("mika.mikic@google.com")
	if err != nil {
		panic(err)
	}
	fmt.Println(user)

}
