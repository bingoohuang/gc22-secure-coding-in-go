package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/caarlos0/httperr"
	_ "modernc.org/sqlite"
)

type API struct {
	db       *sql.DB
	sessions map[int]string
}

func main() {
	if !strings.Contains(os.Args[0], "go") {
		fmt.Print(grumpy)
		fmt.Println(`DO NOT BUILD OR INSTALL THIS!`)
		os.Exit(1)
	}

	api := &API{}

	var err error
	if api.db, err = InitDB(); err != nil {
		log.Fatal(err)
	}

	router := http.NewServeMux()

	// Auth
	router.Handle("/login", httperr.NewF(api.Login))

	// Users
	router.Handle("/user", httperr.NewF(api.User))
	router.Handle("/users", httperr.NewF(api.Users))
	router.Handle("/friend", httperr.NewF(api.Friend))
	router.Handle("/friends", httperr.NewF(api.Friends))

	// Images
	router.Handle("/images", httperr.NewF(api.Pictures))
	router.Handle("/imgs/", httperr.NewF(api.Image))
	router.Handle("/upload", httperr.NewF(api.Upload))

	log.Printf("Listening on port 8081")
	log.Fatal(http.ListenAndServe("127.0.0.1:8081", router))
}
