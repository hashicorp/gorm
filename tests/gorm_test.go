package tests_test

import (
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"testing"
)

func TestConvertNullToZeroValues(t *testing.T) {
	dialect := DB.Dialector.Name()
	switch dialect {
	case "mysql", "sqlserver":
		// these dialects do not support the "returning" clause
		return
	default:
		// This user struct will leverage the existing users table, but override
		// the Name field to default to null.
		type user struct {
			gorm.Model
			Name string `gorm:"default:null"`
		}
		u := user{}
		DB.Config.ConvertNullToZeroValues = true
		c := DB.Callback().Create().Get("gorm:create")
		t.Cleanup(func() {
			DB.Callback().Create().Replace("gorm:create", c)
			DB.Config.ConvertNullToZeroValues = false
		})
		DB.Callback().Create().Replace("gorm:create", callbacks.Create(&callbacks.Config{WithReturning: true}))

		if results := DB.Create(&u); results.Error != nil {
			t.Fatalf("errors happened on create: %v", results.Error)
		} else if results.RowsAffected != 1 {
			t.Fatalf("rows affected expects: %v, got %v", 1, results.RowsAffected)
		} else if u.ID == 0 {
			t.Fatalf("ID expects : not equal 0, got %v", u.ID)
		}

		got := user{}
		results := DB.First(&got, "id = ?", u.ID)
		if results.Error != nil {
			t.Fatalf("errors happened on first: %v", results.Error)
		} else if results.RowsAffected != 1 {
			t.Fatalf("rows affected expects: %v, got %v", 1, results.RowsAffected)
		} else if got.ID != u.ID {
			t.Fatalf("first expects: %v, got %v", u, got)
		}

		results = DB.Select("name").Find(&got)
		if results.Error != nil {
			t.Fatalf("errors happened on first: %v", results.Error)
		} else if results.RowsAffected != 1 {
			t.Fatalf("rows affected expects: %v, got %v", 1, results.RowsAffected)
		} else if got.ID != u.ID {
			t.Fatalf("first expects: %v, got %v", u, got)
		}

		u.Name = "jinzhu"
		if results := DB.Save(&u); results.Error != nil {
			t.Fatalf("errors happened on update: %v", results.Error)
		} else if results.RowsAffected != 1 {
			t.Fatalf("rows affected expects: %v, got %v", 1, results.RowsAffected)
		}
	}
}
