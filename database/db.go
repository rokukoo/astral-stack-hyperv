package database

import (
	"astralstack-hyperv/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() error {
	//dbName := "data.db"
	dbName := `D:\Projects\astralstack-hyperv\data.db`
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	DB = db
	err = DB.AutoMigrate(
		&model.VirtualHardDisk{},
		&model.SystemImage{},
	)
	if err != nil {
		return err
	}
	return nil
}
