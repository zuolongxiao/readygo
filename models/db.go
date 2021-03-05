package models

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/zuolongxiao/readygo/pkg/settings"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// DB database
var DB *gorm.DB

func init() {
	var err error

	charset := "utf8mb4"
	dialector := mysql.Open(
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local",
			settings.DatabaseSetting.User,
			settings.DatabaseSetting.Password,
			settings.DatabaseSetting.Host,
			settings.DatabaseSetting.Name,
			charset,
		))

	consoleLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)

	config := gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   settings.DatabaseSetting.Prefix,
			SingularTable: true,
		},
		Logger: consoleLogger,
	}

	DB, err = gorm.Open(dialector, &config)

	if err != nil {
		log.Println(err)
	}

	sqlDB, err := DB.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
}
