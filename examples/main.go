package main

import (
	"fmt"
	"log"

	paginator "github.com/hiro-riveros/gorm-paginator"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	ID    uint
	Name  string
	Email string
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&User{})

	for i := 1; i <= 50; i++ {
		db.Create(&User{
			Name:  fmt.Sprintf("User %d", i),
			Email: fmt.Sprintf("user%d@example.com", i),
		})
	}

	params := paginator.Params{Page: 2, Limit: 10, OrderBy: "id desc"}
	users, meta, err := paginator.Paginate(&User{}, db.Model(&User{}), params)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Total usuarios: %d\n", meta.TotalRecords)
	for _, u := range users.([]*User) {
		fmt.Println(u.ID, u.Name)
	}
}
