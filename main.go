package main

import (
	"flag"
	"fmt"
	"net/http"
	"runtime"

	"github.com/zuolongxiao/readygo/commands"
	"github.com/zuolongxiao/readygo/pkg/jobs"
	"github.com/zuolongxiao/readygo/pkg/settings"
	"github.com/zuolongxiao/readygo/pkg/utils"
	"github.com/zuolongxiao/readygo/routing"
)

func main() {
	flag.Parse()

	var cmd string
	args := flag.Args()
	if len(args) > 0 {
		cmd = args[0]
	}

	maps := make(map[string]string)
	maps["info"] = "Show information"
	maps["http:start"] = "Start HTTP server"
	maps["admin:create"] = "Create an admin"
	maps["permission:load"] = "Load permissions to database"
	maps["migration:run"] = "Run migration"

	if _, ok := maps[cmd]; !ok {
		fmt.Println("Available CMDs:")
		for key, val := range maps {
			fmt.Printf("%s: %s\n", key, val)
		}

		return
	}

	flagSet := flag.NewFlagSet(cmd, flag.ExitOnError)

	switch cmd {
	case "info":
		fmt.Printf("Go version: %s\n", runtime.Version())
		fmt.Printf("App version: %s\n", settings.AppSetting.Version)
		fmt.Printf("Run mode: %s\n", settings.ServerSetting.RunMode)
		fmt.Printf("Server host: %s\n", settings.ServerSetting.HTTPHost)
		fmt.Printf("Server port: %d\n", settings.ServerSetting.HTTPPort)

	case "http:start":
		fmt.Printf("HTTP server running on: %s:%d\n", settings.ServerSetting.HTTPHost, settings.ServerSetting.HTTPPort)
		startHTTP()

	case "admin:create":
		var username string
		var password string
		flagSet.StringVar(&username, "u", "", "Username, required")
		flagSet.StringVar(&password, "p", "", "Password, optional")

		flagSet.Parse(args[1:])

		if password == "" {
			password = utils.RandomString(12)
		}

		if err := commands.CreateAdmin(username, password); err != nil {
			fmt.Printf("Admin created failed with error:\n")
			fmt.Println(err.Error())
			return
		}

		fmt.Printf("Admin created successfully\n")
		fmt.Printf("Username: %s\n", username)
		fmt.Printf("Password: %s\n", password)

	case "permission:load":
		commands.LoadPermissions()

	case "migration:run":
		commands.RunMigration()

	default:
		fmt.Printf("Not implemented command: %s\n", cmd)
	}
}

func startHTTP() {
	go jobs.SetPermissions()

	router := routing.Setup()

	s := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", settings.ServerSetting.HTTPHost, settings.ServerSetting.HTTPPort),
		Handler:        router,
		ReadTimeout:    settings.ServerSetting.ReadTimeout,
		WriteTimeout:   settings.ServerSetting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	err := s.ListenAndServe()
	fmt.Println(err.Error())
}
