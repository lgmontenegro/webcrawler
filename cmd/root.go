package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/lgmontenegro/webcrawler/internal/app"
	"github.com/spf13/cobra"
)

var File string
var URL string
var rootCmd = &cobra.Command{
	Use:   "webcrawler",
	Short: "Webcrawler is an assessment test to be part of Parse Digital Consulting",
	Long: `Webcrawler is an assessment test to be part of Parse Digital Consulting, 
	parser construct without any 3rd part libraries`,
	Run: func(cmd *cobra.Command, args []string) {
		a := app.App{
			InputURL: []string{},
		}
		if File != "" {
			urls, err := openURLsFile(File)
			if err != nil {
				log.Fatalln(err)
			}

			a.InputURL = urls		
		}

		if URL != "" {
			a.InputURL = splitURLsArg(URL)
		}

		if len(a.InputURL) == 0 {
			log.Fatalln("No URL supplied")
		}

		if a.Execute() {
			os.Exit(0)
		}

		os.Exit(1)
	},
}

func Execute() {
	rootCmd.PersistentFlags().StringVarP(&File, "file", "f", "", "single file with URLs to be crawled")
	rootCmd.PersistentFlags().StringVarP(&URL, "url", "u", "", "URLs list to be crawled, separeted by comma (,)")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func openURLsFile(filePath string) (urls []string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return []string{}, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return []string{}, err
	}

	return urls, nil
}

func splitURLsArg(urlsString string) (urls []string) {
	urls = strings.Split(urlsString, ",")
	for i, url := range urls {
		urls[i] = strings.TrimSpace(url)
	}
	return urls
}
