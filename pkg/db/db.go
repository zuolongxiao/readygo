package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"readygo/pkg/settings"

	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// DB database
var DB *gorm.DB

func Setup() error {
	if settings.Database.Type == "MySQL" {
		return setupMysql()
	}

	if settings.Database.Type == "SQLite" {
		return setupSqlite()
	}

	return fmt.Errorf("database %s has not support yet", settings.Database.Type)
}

func setupMysql() error {
	var err error

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		settings.MySQL.User,
		settings.MySQL.Password,
		settings.MySQL.Host,
		settings.MySQL.Port,
		settings.MySQL.Name,
		settings.MySQL.Charset,
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
	sqlDB.SetMaxIdleConns(settings.MySQL.MaxIdleConns)
	sqlDB.SetMaxOpenConns(settings.MySQL.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(settings.MySQL.ConnMaxLifetime)

	return nil
}

func setupSqlite() error {
	var err error

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
	DB, err = gorm.Open(sqlite.Open(settings.SQLite.Name), &config)

	return err
}
