package api

import (
	"fmt"
	"net/http"

	"github.com/aubm/velib-toulouse-carte/backend/bikes"
	"github.com/aubm/velib-toulouse-carte/backend/log"
	"github.com/paulmach/go.geojson"
	"google.golang.org/appengine"
)

type StationsHandlers struct {
	Stations bikes.StationsManager `inject:""`
	Logger   log.Logger            `inject:""`
}

func (h *StationsHandlers) List(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)

	stations, err := h.Stations.Search(ctx, bikes.StationsSearchOptions{
		ContractName: "Toulouse",
	})
	if err != nil {
		httpError(w, serverError, http.StatusInternalServerError)
		h.Logger.Error(ctx, fmt.Sprintf("failed to search stations: %v", err.Error()))
		return
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	writeJSON(w, h.toGeoJsonFeatureCollection(stations), http.StatusOK)
}

func (h *StationsHandlers) toGeoJsonFeatureCollection(stations []bikes.Station) *geojson.FeatureCollection {
	fc := geojson.NewFeatureCollection()
	for _, station := range stations {
		feature := geojson.NewFeature(geojson.NewPointGeometry([]float64{station.Position.Lng, station.Position.Lat}))
		feature.SetProperty("number", station.Number)
		feature.SetProperty("name", station.Name)
		feature.SetProperty("address", station.Address)
		feature.SetProperty("banking", station.Banking)
		feature.SetProperty("bonus", station.Bonus)
		feature.SetProperty("status", station.Status)
		feature.SetProperty("contractName", station.ContractName)
		feature.SetProperty("bikeStands", station.BikeStands)
		feature.SetProperty("availableBikeStands", station.AvailableBikeStands)
		feature.SetProperty("availableBikes", station.AvailableBikes)
		feature.SetProperty("lastUpdateTimestamp", station.LastUpdate())
		fc.AddFeature(feature)
	}
	return fc
}
