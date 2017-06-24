package histories

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetLastLendingTime_Empty(t *testing.T) {
	db, err := NewDatabase(":memory:", nil)
	if err != nil {
		t.Fatalf("NewDatabase: %v", err)
	}
	defer db.Close()

	actual := db.GetLastLendingTime()
	assert.Equal(t, time.Unix(0, 0), actual)
}

func TestGetLastLendingTime_NotEmpty(t *testing.T) {
	db, err := NewDatabase(":memory:", nil)
	if err != nil {
		t.Fatalf("NewDatabase: %v", err)
	}
	defer db.Close()

	t1 := time.Date(2017, time.January, 1, 0, 0, 0, 0, time.Local)
	t2 := time.Date(2017, time.January, 2, 0, 0, 0, 0, time.Local)
	row1 := PoloniexLendingRow{LendingID: 1, Open: t1}
	row2 := PoloniexLendingRow{LendingID: 2, Open: t2}
	db.db.Create(&row1)
	db.db.Create(&row2)

	actual := db.GetLastLendingTime()
	assert.Equal(t, t2.Unix(), actual.Unix())
}
