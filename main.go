package main

import (
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/evanw/esbuild/pkg/api"
)

//go:embed frontend/public/index.html
var html string

type Data struct {
	JS template.JS
}

func main() {
	bundled := api.Build(api.BuildOptions{
		EntryPoints: []string{"frontend/App.jsx"},
		Bundle:      true,
	})
	handleMessages(bundled.Errors)

	js := string(bundled.OutputFiles[0].Contents)

	isDevelopment := os.Getenv("ENV") == "development"

	if !isDevelopment {
		minified := api.Transform(string(bundled.OutputFiles[0].Contents), api.TransformOptions{
			MinifyWhitespace:  true,
			MinifyIdentifiers: true,
			MinifySyntax:      true,
		})
		handleMessages(minified.Errors)

		js = string(minified.Code)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	port = fmt.Sprintf(":%s", port)

	tmpl, err := template.New("index").Parse(html)
	if err != nil {
		log.Fatal(err)
	}

	data := &Data{
		JS: template.JS(js),
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	log.Fatal(http.ListenAndServe(port, nil))
}

func handleMessages(messages []api.Message) {
	if len(messages) < 1 {
		return
	}

	for _, message := range messages {
		detail := fmt.Sprintf("\n%s\n    %s:%d:%d\n      %d | %s", message.Text, message.Location.File, message.Location.Line, message.Location.Column, message.Location.Line, message.Location.LineText)

		log.Println(detail)
	}

	os.Exit(1)
}
