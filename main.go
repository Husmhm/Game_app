package main

import (
	"fmt"
	"gameApp/entity"
	"gameApp/repository/mysql"
)

func main() {
	testUserMysqlRepo()
}
func testUserMysqlRepo() {
	mysqlRepo := mysql.New()

	user := entity.User{
		ID:          0,
		PhoneNumber: "0919",
		Name:        "Hossein ",
	}
	isUnique, err := mysqlRepo.IsPhoneNumberUnique(user.PhoneNumber)
	if err != nil {
		fmt.Println("unique err", err)
	}

	fmt.Println("isUnique", isUnique)
	if isUnique == true {
		createdUser, err := mysqlRepo.Register(user)

		if err != nil {
			fmt.Println("register user", err)
		} else {
			fmt.Println("created user", createdUser)
		}
	}

}
