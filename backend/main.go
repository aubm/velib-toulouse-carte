package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aubm/velib-toulouse-carte/backend/api"
	"github.com/aubm/velib-toulouse-carte/backend/bikes"
	"github.com/aubm/velib-toulouse-carte/backend/shared"
	"github.com/facebookgo/inject"
)

var (
	stationsHandlers = &api.StationsHandlers{}
	stationsManager  = &bikes.DefaultStationsManager{}
	config           *shared.AppConfig
)

func main() {
	var err error
	config, err = shared.ConfigFromEnvVars()
	checkErrorAndExit(err)

	if err := inject.Populate(
		stationsHandlers,
		stationsManager,
		config,
	); err != nil {
		checkErrorAndExit(err)
	}

	http.HandleFunc("/api/v1/stations", stationsHandlers.List)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("Defaulting to port %s", port)
	}

	log.Printf("Listening on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func checkErrorAndExit(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
