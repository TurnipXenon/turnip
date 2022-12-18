package main

import (
	"net/http"
)

func hello(response http.ResponseWriter, _ *http.Request) {
	response.WriteHeader(http.StatusOK)
	response.Write([]byte("Hello"))
}

func main() {
	// todo(turnip): document how static works
	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/hello/", hello)
	// todo(turnip): make port dynamic
	http.ListenAndServe(":8090", nil)
}
