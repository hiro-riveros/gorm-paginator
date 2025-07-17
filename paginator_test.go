package paginator

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TestUser struct {
	ID      uint
	Name    string
	Email   string
	Wallets []TestWallet `gorm:"foreignKey:UserID"`
}

type TestWallet struct {
	ID     uint
	UserID uint
	Amount int
}

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	err = db.AutoMigrate(&TestUser{}, &TestWallet{})
	if err != nil {
		t.Fatal(err)
	}

	for i := 1; i <= 25; i++ {
		user := TestUser{
			Name:  "User",
			Email: "user@example.com",
		}
		db.Create(&user)

		// Add wallets only to some users
		if i%5 == 0 {
			db.Create(&TestWallet{
				UserID: user.ID,
				Amount: i * 100,
			})
		}
	}

	return db
}

func TestPaginateBasic(t *testing.T) {
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

	if meta.Page != 2 {
		t.Errorf("esperaba página 2, obtuvo %d", meta.Page)
	}
}

func TestPaginateWithPreload(t *testing.T) {
	db := setupTestDB(t)

	users, meta, err := Paginate(&TestUser{}, db.Model(&TestUser{}), Params{Page: 1, Limit: 10}, "Wallets")
	if err != nil {
		t.Fatal(err)
	}

	userSlice := users.([]*TestUser)
	if len(userSlice) != 10 {
		t.Errorf("esperaba 10 resultados, obtuvo %d", len(userSlice))
	}

	// Verify that Wallets are preloaded on users who have wallets
	for _, u := range userSlice {
		if len(u.Wallets) > 0 {
			if u.Wallets[0].Amount == 0 {
				t.Errorf("wallet cargado pero con amount 0 inesperado")
			}
		}
	}

	if meta.TotalRecords != 25 {
		t.Errorf("esperaba total 25, obtuvo %d", meta.TotalRecords)
	}
}

func TestPaginateWithWhere(t *testing.T) {
	db := setupTestDB(t)

	query := db.Model(&TestUser{}).Where("email LIKE ?", "%@example.com")
	users, meta, err := Paginate(&TestUser{}, query, Params{Page: 1, Limit: 5})
	if err != nil {
		t.Fatal(err)
	}

	userSlice := users.([]*TestUser)
	if len(userSlice) != 5 {
		t.Errorf("esperaba 5 resultados, obtuvo %d", len(userSlice))
	}

	if meta.TotalRecords != 25 {
		t.Errorf("esperaba total 25, obtuvo %d", meta.TotalRecords)
	}
}

func TestPaginateOrderBy(t *testing.T) {
	db := setupTestDB(t)

	params := Params{
		Page:    1,
		Limit:   10,
		OrderBy: "id desc",
	}
	users, meta, err := Paginate(&TestUser{}, db.Model(&TestUser{}), params)
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

	// Check descending order by ID
	prevID := userSlice[0].ID
	for i := 1; i < len(userSlice); i++ {
		if userSlice[i].ID > prevID {
			t.Errorf("resultado no ordenado desc por ID")
		}
		prevID = userSlice[i].ID
	}
}
