package app

import (
	"fmt"
	"net/http"
	"os"

	"github.com/aubm/velib-toulouse-carte/backend/api"
	"github.com/aubm/velib-toulouse-carte/backend/bikes"
	apphttp "github.com/aubm/velib-toulouse-carte/backend/http"
	"github.com/aubm/velib-toulouse-carte/backend/log"
	"github.com/aubm/velib-toulouse-carte/backend/shared"
	"github.com/facebookgo/inject"
)

var (
	stationsHandlers   = &api.StationsHandlers{}
	stationsManager    = &bikes.DefaultStationsManager{}
	httpClientProvider = &apphttp.AppEngineClientProvider{}
	logger             = &log.AppEngineLogger{}
	config             *shared.AppConfig
)

func init() {
	var err error
	config, err = shared.ConfigFromEnvVars()
	checkErrorAndExit(err)

	if err := inject.Populate(
		stationsHandlers,
		stationsManager,
		httpClientProvider,
		logger,
		config,
	); err != nil {
		checkErrorAndExit(err)
	}

	http.HandleFunc("/api/v1/stations", stationsHandlers.List)
}

func checkErrorAndExit(err error) {
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
