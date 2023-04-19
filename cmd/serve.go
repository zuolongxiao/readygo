package cmd

import (
	"fmt"
	"log"
	"net/http"

	"readygo/pkg/db"
	"readygo/pkg/jobs"
	"readygo/pkg/settings"
	"readygo/pkg/store"
	"readygo/routing"

	"github.com/spf13/cobra"
)

var host string
var port uint32

// go run main.go serve
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start HTTP server",
	Long:  `Start HTTP server`,
	Run: func(cmd *cobra.Command, args []string) {
		if host == "" {
			host = settings.Server.HTTPHost
		}
		if port == 0 {
			port = settings.Server.HTTPPort
		}
		if port <= 0 || port > 65535 {
			fmt.Println("Invalid HTTP port")
			return
		}
		startHTTP()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVar(&host, "host", "", "HTTP host to listen on")
	serveCmd.Flags().Uint32VarP(&port, "port", "p", 0, "HTTP port to listen on")
}

func startHTTP() {
	if err := db.Setup(); err != nil {
		fmt.Println(err)
		return
	}

	if settings.Captcha.Enabled {
		if err := store.Setup(); err != nil {
			fmt.Println(err)
			return
		}
	}

	go jobs.SetPermissions()

	router := routing.Setup()
	server := &http.Server{
		Addr:           fmt.Sprintf("%s:%d", host, port),
		Handler:        router,
		ReadTimeout:    settings.Server.ReadTimeout,
		WriteTimeout:   settings.Server.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("HTTP server started on: %s:%d\n", host, port)

	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}
