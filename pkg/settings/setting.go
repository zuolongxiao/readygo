package settings

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

// App App
type App struct {
	Name            string
	Version         string
	PageSize        int
	JwtSecret       string
	JwtExpires      time.Duration
	JwtIssuer       string
	RuntimeRootPath string
	SuperAdminID    uint64
}

// AppSetting AppSetting
var AppSetting = &App{}

// Server Server
type Server struct {
	RunMode      string
	HTTPHost     string
	HTTPPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// ServerSetting ServerSetting
var ServerSetting = &Server{}

// Database Database
type Database struct {
	Type     string
	User     string
	Password string
	Host     string
	Name     string
	Prefix   string
}

// DatabaseSetting DatabaseSetting
var DatabaseSetting = &Database{}

var cfg *ini.File

// Setup initialize the configuration instance
func init() {
	var err error
	cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("setting.init, fail to parse 'conf/app.ini': %v", err)
	}

	mapTo("app", AppSetting)
	mapTo("server", ServerSetting)
	mapTo("database", DatabaseSetting)

	ServerSetting.ReadTimeout = ServerSetting.ReadTimeout * time.Second
	ServerSetting.WriteTimeout = ServerSetting.WriteTimeout * time.Second
}

// mapTo map section
func mapTo(section string, v interface{}) {
	err := cfg.Section(section).MapTo(v)
	if err != nil {
		log.Fatalf("setting.MapTo %s err: %v", section, err)
	}
}
