package settings

import (
	"strings"
	"time"

	"github.com/spf13/viper"
)

const Version string = "1.1.0"

// App
var App struct {
	Name            string
	PageSize        uint32
	SuperAdminID    uint64
	RuntimeRootPath string
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
	Type            string
	User            string
	Password        string
	Host            string
	Port            uint32
	Name            string
	Prefix          string
	Charset         string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
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
	viper.SetDefault("Database.Host", "127.0.0.1")
	viper.SetDefault("Database.Port", "3306")
	viper.SetDefault("Database.User", "root")
	viper.SetDefault("Database.Name", "readygo")
	viper.SetDefault("Database.Prefix", "g_")
	viper.SetDefault("Database.Charset", "utf8mb4")
	viper.SetDefault("Database.MaxIdleConns", "10")
	viper.SetDefault("Database.MaxOpenConns", "100")
	viper.SetDefault("Database.ConnMaxLifetime", "3600")
	Database.Type = viper.GetString("Database.Type")
	Database.Host = viper.GetString("Database.Host")
	Database.Port = viper.GetUint32("Database.Port")
	Database.User = viper.GetString("Database.User")
	Database.Password = viper.GetString("Database.Password")
	Database.Name = viper.GetString("Database.Name")
	Database.Prefix = viper.GetString("Database.Prefix")
	Database.Charset = viper.GetString("Database.Charset")
	Database.MaxIdleConns = viper.GetInt("Database.MaxIdleConns")
	Database.MaxOpenConns = viper.GetInt("Database.MaxOpenConns")
	Database.ConnMaxLifetime = time.Duration(viper.GetUint64("Database.ConnMaxLifetime")) * time.Second
}
