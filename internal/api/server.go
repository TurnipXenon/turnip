package api

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/TurnipXenon/Turnip/internal/models"
)

func hello(response http.ResponseWriter, _ *http.Request) {
	// todo: delete
	response.WriteHeader(http.StatusOK)
	response.Write([]byte("Hello"))
}

func handleIndex(response http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	if path != "/" {
		// todo: not found page!
		response.WriteHeader(http.StatusNotFound)
		return
	}

	data := models.UserData{
		Host: request.Host,
	}

	// Reference: https://gowebexamples.com/templates/
	tmpl := template.Must(template.ParseFiles("./internal/templates/main.html"))
	err := tmpl.Execute(response, data)
	if err != nil {
		print(err)
		// todo: log?
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

func InitializeServer(flags models.RunFlags) {
	// setup domain information
	models.InitializeDomainMap()

	// setup server
	mux := http.NewServeMux()

	// todo(turnip): document how static works
	mux.HandleFunc("/api/hello", hello)

	// root-based resources
	serveSingle("/robots.txt", "./assets/robots.txt", mux)
	// todo: favicon
	// todo: sitemap

	mux.Handle("/assets/", http.FileServer(http.Dir("./")))
	mux.HandleFunc("/", handleIndex)

	fmt.Printf("Serving at http://localhost:%d\n", flags.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", flags.Port), mux)
	if err != nil {
		log.Fatalln(err)
	}
}
