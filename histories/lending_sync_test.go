package histories

import (
	"testing"
	"time"

	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/assert"
)

func Test_LendingSync_GetLastTime_Empty(t *testing.T) {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("gorm.Open: %v", err)
	}
	db.AutoMigrate(&LendingRow{})
	defer db.Close()

	sync := NewLendingSync(db, nil)

	actual := sync.GetLastTime()
	assert.Equal(t, time.Unix(0, 0), actual)
}

func Test_LendingSync_GetLastTime_NotEmpty(t *testing.T) {
	db, err := gorm.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("gorm.Open: %v", err)
	}
	db.AutoMigrate(&LendingRow{})
	defer db.Close()

	sync := NewLendingSync(db, nil)

	t1 := time.Date(2017, time.January, 1, 0, 0, 0, 0, time.Local)
	t2 := time.Date(2017, time.January, 2, 0, 0, 0, 0, time.Local)
	row1 := LendingRow{LendingID: 1, Open: t1}
	row2 := LendingRow{LendingID: 2, Open: t2}
	db.Create(&row1)
	db.Create(&row2)

	actual := sync.GetLastTime()
	assert.Equal(t, t2.Unix(), actual.Unix())
}
