package ioc

import (
	"github.com/CAbrook/golang_learning/config"
	"github.com/CAbrook/golang_learning/internal/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		println(config.Config.DB.DSN)
		panic(err)
	}
	db.Debug()
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}
