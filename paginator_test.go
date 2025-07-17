package paginator

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TestUser struct {
	ID    uint
	Name  string
	Email string
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	db.AutoMigrate(&TestUser{})

	for i := 1; i <= 25; i++ {
		db.Create(&TestUser{
			Name:  "User",
			Email: "user@example.com",
		})
	}

	return db
}

func TestPaginate(t *testing.T) {
	db := setupTestDB(t)

	users, meta, err := Paginate(&TestUser{}, db.Model(&TestUser{}), Params{Page: 2, Limit: 10})
	if err != nil {
		t.Fatal(err)
	}

	userSlice := users.([]*TestUser)
	if len(userSlice) != 10 {
		t.Errorf("esperaba 10 resultados, obtuvo %d", len(userSlice))
	}

	if meta.TotalRecords != 25 {
		t.Errorf("esperaba total 25, obtuvo %d", meta.TotalRecords)
	}
}
