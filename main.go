package main

import (
	"html/template"
	"net/http"
)

func hello(response http.ResponseWriter, _ *http.Request) {
	response.WriteHeader(http.StatusOK)
	response.Write([]byte("Hello"))
}

type MainData struct {
}

func handleTemplate(response http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	if path != "/" {
		// todo: not found page!
		response.WriteHeader(http.StatusNotFound)
		return
	}

	// Reference: https://gowebexamples.com/templates/
	tmpl := template.Must(template.ParseFiles("./internal/templates/main.html"))
	data := MainData{}
	err := tmpl.Execute(response, data)
	if err != nil {
		// todo
		return
	}
}

func serveSingle(pattern string, filename string, mux *http.ServeMux) {
	// from Deleplace @ https://stackoverflow.com/a/14187941/17836168
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		print("Tests")
		http.ServeFile(w, r, filename)
	})
}

func main() {
	mux := http.NewServeMux()

	// todo(turnip): document how static works
	mux.HandleFunc("/api/hello", hello)

	// root-based resources
	serveSingle("/robots.txt", "./assets/robots.txt", mux)
	// todo: favicon
	// todo: sitemap

	mux.Handle("/assets/", http.FileServer(http.Dir("./")))
	mux.HandleFunc("/", handleTemplate)

	// todo(turnip): make port dynamic
	http.ListenAndServe(":8090", mux)
}
