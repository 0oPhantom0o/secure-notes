package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "app",
	Short: "my note application",
	Long:  `A longer description that spans multiple lines and likely contains`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello World")
	},
}

func Execute() {
	if len(os.Args) == 1 {
		os.Args = append(os.Args, "serve")
	}
	_ = rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
