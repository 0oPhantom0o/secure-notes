package cmd

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	appDB "github.com/secure-notes/internal/app"
	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start http server",
	Run: func(cmd *cobra.Command, args []string) {
		app, err := Start()
		if err != nil {
			log.Fatal(err)
		}

		log.Fatal(app.Listen(":3030"))
	},
}

func Start() (*fiber.App, error) {
	app, err := appDB.New(context.Background())
	return app, err
}
