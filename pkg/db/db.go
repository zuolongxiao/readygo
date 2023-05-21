package db

import (
	"fmt"
	"log"
	"os"
	"time"

	"readygo/pkg/settings"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// DB database
var DB *gorm.DB

// Redis client
var RDB *redis.Client

func Setup() error {
	if settings.Redis.Enabled {
		if err := setupRedis(); err != nil {
			return err
		}
	}

	var err error
	switch settings.Database.Type {
	case "MySQL":
		err = setupMysql()
	case "SQLite":
		err = setupSqlite()
	default:
		err = fmt.Errorf("database %s has not support yet", settings.Database.Type)
	}

	return err
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

func setupRedis() error {
	RDB = redis.NewClient(&redis.Options{
		Addr:     settings.Redis.Addr,
		Password: settings.Redis.Password,
		DB:       settings.Redis.DB,
	})

	return nil
}
