package cmd

import (
	"log"

	appDB "github.com/secure-notes/internal/app"
	apihttp "github.com/secure-notes/internal/http"

	"github.com/spf13/cobra"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "start http server",
	Run: func(cmd *cobra.Command, args []string) {
		if err := db(); err != nil {
			log.Fatal(err)
		}

		app := apihttp.NewServer()
		log.Fatal(app.Listen(":3030"))
	},
}

func db() error {
	return appDB.Db()
}
