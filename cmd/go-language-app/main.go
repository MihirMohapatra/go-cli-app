package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/MihirMohapatra/go-cli-app/internal/todos"
)

func main() {
	port := flag.Int("port", 8080, "HTTP server port")
	flag.Parse()

	store := todos.NewStore()
	handler := todos.NewHandler(store)
	addr := fmt.Sprintf(":%d", *port)

	log.Printf("listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, handler); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
