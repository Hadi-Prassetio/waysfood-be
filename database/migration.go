package database

import (
	"fmt"
	"waysfood/models"
	"waysfood/pkg/mysql"
)

func RunMigration() {
	err := mysql.DB.AutoMigrate(&models.User{}, &models.Product{})

	if err != nil {
		fmt.Println(err)
		panic("migration error")
	}
	fmt.Println("migration success")
}
