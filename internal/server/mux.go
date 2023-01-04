package server

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/twitchtv/twirp"

	"github.com/TurnipXenon/turnip_api/rpc/turnip"

	"github.com/TurnipXenon/turnip/internal/models"
)

type Mux struct {
	HostMap map[string]models.Host
}

func (m *Mux) handleIndex(response http.ResponseWriter, request *http.Request) {
	// todo(turnip): delete
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

func (m *Mux) serveSingle(pattern string, filename string, Mux *mux.Router) {
	// from Deleplace @ https://stackoverflow.com/a/14187941/17836168
	Mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		print("Tests")
		http.ServeFile(w, r, filename)
	})
}

func RunServeMux(s *Server, flags *models.RunFlags) {
	//m := Mux{
	//	HostMap: s.Storage.GetHostMap(),
	//}

	// setup turnip
	ti := NewTurnipHandler(s)
	twirpHandler := turnip.NewTurnipServer(ti, twirp.WithServerPathPrefix("/api"))

	// grab header details
	authWrapper := NewAuthMiddleware(twirpHandler, s)

	// todo: we might remove mux later
	//router := mux.NewRouter()
	//router.Handle(twirpHandler.PathPrefix(), twirpHandler)

	// root-based resources
	//m.serveSingle("/robots.txt", "./assets/robots.txt", router)
	// todo: favicon
	// todo: sitemap

	// from dodgy_coder @ https://stackoverflow.com/a/21251658/17836168
	//router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	// todo: take a look at CORS more for safety stuff
	flags.CorsAllowList = append(flags.CorsAllowList, []string{
		"http://localhost:3000",
		"http://127.0.0.1:3000",
	}...)
	// todo: remove in the future
	fmt.Println("Allowing CORS on the following hosts:")
	fmt.Println(flags.CorsAllowList)
	c := cors.New(cors.Options{
		AllowedOrigins:   flags.CorsAllowList,
		AllowCredentials: false,
		Debug:            true,
		AllowedHeaders: []string{
			"Content-Type",
			"Authorization",
		},
	})
	c.Log = nil // remove logging for the cors middleware
	corsHandler := c.Handler(authWrapper)

	// todo: enforce timeouts
	srv := &http.Server{
		Handler: corsHandler, // todo: fix
		//Handler: http.TimeoutHandler(corsHandler, 6*time.Second, "Timeout"), // todo: fix
		Addr: fmt.Sprintf(":%d", flags.Port),
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 6 * time.Second,
		ReadTimeout:  6 * time.Second, // todo: when local, extend timeout for debugging
	}

	fmt.Printf("Serving at http://localhost:%d\n", flags.Port)
	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
