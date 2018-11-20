package main

import (
	"fmt"
	"lenslocked/models"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "1234"
	dbname   = "lenslocked_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	defer us.Close()

	// us.DestructiveReset()

	// u := models.User{
	// 	Name:  "Mika",
	// 	Email: "mika@email.com",
	// }

	// if err := us.Create(&u); err != nil {
	// 	panic(err)
	// }

	user, err := us.ByID(1)
	if err != nil {
		panic(err)
	}
	fmt.Println(user)

}
