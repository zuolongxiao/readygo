package db

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

func Setup() error {
	// log.Println("Setup DB")
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		settings.Database.User,
		settings.Database.Password,
		settings.Database.Host,
		settings.Database.Port,
		settings.Database.Name,
		settings.Database.Charset,
	)
	dialector := mysql.Open(dsn)

	logLevel := logger.Silent
	if settings.Gin.Mode == "debug" {
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
			TablePrefix:   settings.Database.Prefix,
			SingularTable: true,
		},
		Logger: consoleLogger,
	}

	DB, err = gorm.Open(dialector, &config)
	if err != nil {
		log.Println(err)
		return err
	}

	// Connection Pool setting
	sqlDB, _ := DB.DB()
	sqlDB.SetMaxIdleConns(settings.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(settings.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(settings.Database.ConnMaxLifetime)

	return nil
}
