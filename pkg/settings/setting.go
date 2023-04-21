package settings

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

// App
var App struct {
	Name            string
	PageSize        uint32
	SuperAdminID    uint64
	RuntimeRootPath string
}

// CORS
var CORS struct {
	AllowOrigin  string
	AllowMethods string
}

// Gin
var Gin struct {
	Mode string
}

// JWT
var JWT struct {
	Secret  string
	Expires time.Duration
	Issuer  string
}

// Server
var Server struct {
	HTTPHost     string
	HTTPPort     uint32
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// Database
var Database struct {
	Type   string
	Prefix string
}

// MySQL
var MySQL struct {
	User            string
	Password        string
	Host            string
	Port            uint32
	Name            string
	Charset         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

// SQLite
var SQLite struct {
	Name string
}

// Redis
var Redis struct {
	Enabled  bool
	Addr     string
	Password string
	DB       int
}

// Captcha
var Captcha struct {
	Enabled bool
	Store   string
	Height  int
	Width   int
	Length  int
	Prefix  string
	Expires time.Duration
}

// Load initialize the configuration instance
func Load() {
	// log.Println("settings.Load")

	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)

	viper.SetDefault("App.Name", "readygo")
	viper.SetDefault("App.PageSize", "20")
	viper.SetDefault("App.RuntimeRootPath", "runtime/")
	viper.SetDefault("App.SuperAdminID", "1")
	App.Name = viper.GetString("App.Name")
	App.PageSize = viper.GetUint32("App.PageSize")
	App.SuperAdminID = viper.GetUint64("App.SuperAdminID")
	App.RuntimeRootPath = viper.GetString("App.RuntimeRootPath")

	viper.SetDefault("CORS.AllowOrigin", "*")
	viper.SetDefault("CORS.AllowMethods", "OPTIONS, GET, POST, PUT, DELETE")
	CORS.AllowOrigin = viper.GetString("CORS.AllowOrigin")
	CORS.AllowMethods = viper.GetString("CORS.AllowMethods")

	viper.SetDefault("Gin.Mode", "debug")
	Gin.Mode = viper.GetString("Gin.Mode")

	viper.SetDefault("JWT.Expires", "8")
	viper.SetDefault("JWT.Issuer", "readygo")
	JWT.Expires = time.Duration(viper.GetUint64("JWT.Expires")) * time.Hour
	JWT.Secret = viper.GetString("JWT.Secret")
	JWT.Issuer = viper.GetString("JWT.Issuer")

	viper.SetDefault("Server.HTTPHost", "127.0.0.1")
	viper.SetDefault("Server.HTTPPort", "9331")
	viper.SetDefault("Server.ReadTimeout", "60")
	viper.SetDefault("Server.WriteTimeOut", "60")
	Server.HTTPHost = viper.GetString("Server.HTTPHost")
	Server.HTTPPort = viper.GetUint32("Server.HTTPPort")
	Server.ReadTimeout = time.Duration(viper.GetUint64("Server.ReadTimeout")) * time.Second
	Server.WriteTimeout = time.Duration(viper.GetUint64("Server.WriteTimeout")) * time.Second

	viper.SetDefault("Database.Type", "MySQL")
	viper.SetDefault("Database.Prefix", "g_")
	Database.Type = viper.GetString("Database.Type")
	Database.Prefix = viper.GetString("Database.Prefix")

	viper.SetDefault("MySQL.Host", "127.0.0.1")
	viper.SetDefault("MySQL.Port", "3306")
	viper.SetDefault("MySQL.User", "root")
	viper.SetDefault("MySQL.Name", "readygo")
	viper.SetDefault("MySQL.Charset", "utf8mb4")
	viper.SetDefault("MySQL.MaxIdleConns", "10")
	viper.SetDefault("MySQL.MaxOpenConns", "100")
	viper.SetDefault("MySQL.ConnMaxLifetime", "3600")
	MySQL.Host = viper.GetString("Database.Host")
	MySQL.Port = viper.GetUint32("Database.Port")
	MySQL.User = viper.GetString("Database.User")
	MySQL.Password = viper.GetString("Database.Password")
	MySQL.Name = viper.GetString("Database.Name")
	MySQL.Charset = viper.GetString("Database.Charset")
	MySQL.MaxIdleConns = viper.GetInt("Database.MaxIdleConns")
	MySQL.MaxOpenConns = viper.GetInt("Database.MaxOpenConns")
	MySQL.ConnMaxLifetime = time.Duration(viper.GetUint64("Database.ConnMaxLifetime")) * time.Second

	viper.SetDefault("SQLite.Name", "db.sqlite")
	SQLite.Name = viper.GetString("SQLite.Name")

	viper.SetDefault("Redis.Enabled", "0")
	viper.SetDefault("Redis.Addr", "127.0.0.1:6379")
	viper.SetDefault("Redis.DB", "0")
	Redis.Enabled = viper.GetBool("Redis.Enabled")
	Redis.Addr = viper.GetString("Redis.Addr")
	Redis.Password = viper.GetString("Redis.Password")
	Redis.DB = viper.GetInt("Redis.DB")

	viper.SetDefault("Captcha.Enabled", "0")
	viper.SetDefault("Captcha.Store", "Memory")
	viper.SetDefault("Captcha.Height", "40")
	viper.SetDefault("Captcha.Width", "100")
	viper.SetDefault("Captcha.Length", "6")
	viper.SetDefault("Captcha.Prefix", "captcha_")
	viper.SetDefault("Captcha.Expires", "600")
	Captcha.Enabled = viper.GetBool("Captcha.Enabled")
	Captcha.Store = viper.GetString("Captcha.Store")
	Captcha.Height = viper.GetInt("Captcha.Height")
	Captcha.Width = viper.GetInt("Captcha.Width")
	Captcha.Length = viper.GetInt("Captcha.Length")
	Captcha.Prefix = viper.GetString("Captcha.Prefix")
	Captcha.Expires = time.Duration(viper.GetUint64("Captcha.Expires")) * time.Second
}
