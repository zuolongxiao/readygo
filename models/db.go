package models

import (
	"fmt"
	"log"
	"os"
	"time"

	"readygo/pkg/settings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// DB database
var DB *gorm.DB

func init() {
	var err error

	dialector := mysql.Open(
		fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=True&loc=Local",
			settings.DatabaseSetting.User,
			settings.DatabaseSetting.Password,
			settings.DatabaseSetting.Host,
			settings.DatabaseSetting.Name,
			settings.DatabaseSetting.Charset,
		))

	logLevel := logger.Silent
	if settings.ServerSetting.RunMode == "debug" {
		logLevel = logger.Info
	}
	consoleLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second,
			LogLevel:      logLevel,
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

	sqlDB, _ := DB.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
}
