package cmd

import (
	"fmt"
	"os"

	"github.com/lgmontenegro/webcrawler/internal/app"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "webcrawler",
	Short: "Hugo is a very fast static site generator",
	Long: `A Fast and Flexible Static Site Generator built with
				  love by spf13 and friends in Go.
				  Complete documentation is available at http://hugo.spf13.com`,
	Run: func(cmd *cobra.Command, args []string) {
		a := app.App{
			InputURL: append([]string{
				"https://gobyexample.com6",
				"https://gobyexample.com5",
				"https://gobyexample.com4",
				"https://gobyexample.com3",
				"https://gobyexample.com2",
				"https://gobyexample.com1",
			}),
		}

		a.Execute()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
