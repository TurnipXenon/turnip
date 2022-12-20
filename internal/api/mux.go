package api

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/TurnipXenon/Turnip/internal/models"
	turnipserver "github.com/TurnipXenon/Turnip/internal/server"
)

type Mux struct {
	HostMap map[string]models.Host
}

func (m *Mux) hello(response http.ResponseWriter, _ *http.Request) {
	// todo: delete
	response.WriteHeader(http.StatusOK)
	_, _ = response.Write([]byte("Hello"))
}

func (m *Mux) handleIndex(response http.ResponseWriter, request *http.Request) {
	path := request.URL.Path
	if path != "/" {
		// todo: not found page!
		response.WriteHeader(http.StatusNotFound)
		return
	}

	data := models.UserImpl{
		ActualHost: request.Host,
	}

	data.Initialize(request, m.HostMap)

	// Reference: https://gowebexamples.com/templates/
	tmpl := template.Must(template.ParseFiles("./internal/templates/main.html"))
	err := tmpl.Execute(response, data)
	if err != nil {
		print(err)
		// todo: log?
		return
	}
}

func (m *Mux) serveSingle(pattern string, filename string, Mux *http.ServeMux) {
	// from Deleplace @ https://stackoverflow.com/a/14187941/17836168
	Mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		print("Tests")
		http.ServeFile(w, r, filename)
	})
}

func RunServeMux(s *turnipserver.Server, flags models.RunFlags) {
	m := Mux{
		HostMap: s.Storage.GetHostMap(),
	}

	// setup server
	Mux := http.NewServeMux()

	Mux.HandleFunc("/api/hello", m.hello)

	// root-based resources
	m.serveSingle("/robots.txt", "./assets/robots.txt", Mux)
	// todo: favicon
	// todo: sitemap

	Mux.Handle("/assets/", http.FileServer(http.Dir("./")))
	Mux.HandleFunc("/", m.handleIndex)

	fmt.Printf("Serving at http://localhost:%d\n", flags.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", flags.Port), Mux)
	if err != nil {
		log.Fatalln(err)
	}
}
