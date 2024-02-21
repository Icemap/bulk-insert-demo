package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strconv"
	"time"
)

var tableName = "user"

type User struct {
	ID    int    `gorm:"primaryKey;column:id"`
	Names string `gorm:"column:names"`
}

func (*User) TableName() string {
	return tableName
}

func bulkInsert(dsn string, amount int, batchSize int) time.Duration {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{})

	startTime := time.Now()

	userList := make([]*User, 0, amount)
	for i := 1; i <= amount; i++ {
		userList = append(userList, &User{ID: i, Names: "user" + strconv.Itoa(i)})
	}
	db.CreateInBatches(userList, batchSize)

	endTime := time.Now()

	return endTime.Sub(startTime)
}

func loopInsert(dsn string, amount int) time.Duration {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&User{})

	startTime := time.Now()

	for i := 1; i <= amount; i++ {
		db.Create(&User{ID: i, Names: "user" + strconv.Itoa(i)})
	}

	endTime := time.Now()
	return endTime.Sub(startTime)
}

func main() {
	dsn := "root:@tcp(127.0.0.1:4000)/test"
	insertAmount := []int{10, 100, 1000, 10000, 100000, 1000000, 10000000}

	for _, amount := range insertAmount {
		baseTableName := "user_" + strconv.Itoa(amount)

		// Bulk inserts
		tableName = baseTableName + "_bulk"
		bulkDuring := bulkInsert(dsn, amount, 1000)
		fmt.Printf("Bulk\t%d\t%f\n", amount, bulkDuring.Seconds())

		// Bulk inserts
		tableName = baseTableName + "_loop"
		loopDuring := loopInsert(dsn, amount)
		fmt.Printf("Loop\t%d\t%f\n", amount, loopDuring.Seconds())
	}
}
