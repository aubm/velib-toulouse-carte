package bikes

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	apphttp "github.com/aubm/velib-toulouse-carte/backend/http"

	"github.com/aubm/velib-toulouse-carte/backend/shared"
	"github.com/pkg/errors"
	"time"
)

const (
	apiStationsEndpoint = "https://api.jcdecaux.com/vls/v1/stations"
)

type StationsManager interface {
	Search(ctx context.Context, opts StationsSearchOptions) ([]Station, error)
}

type DefaultStationsManager struct {
	Config     *shared.AppConfig      `inject:""`
	HttpClient apphttp.ClientProvider `inject:""`
}

func (m *DefaultStationsManager) Search(ctx context.Context, opts StationsSearchOptions) ([]Station, error) {
	stations := make([]Station, 0)

	urlQuery := url.Values{}
	urlQuery.Set("apiKey", m.Config.JcdecauxApiKey)
	if opts.ContractName != "" {
		urlQuery.Set("contract", opts.ContractName)
	}

	if err := m.callEndpoint(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s?%s", apiStationsEndpoint, urlQuery.Encode()),
		nil,
		&stations,
	); err != nil {
		return nil, errors.Wrap(err, "failed to perform api call")
	}

	return stations, nil
}

func (m *DefaultStationsManager) callEndpoint(ctx context.Context, httpMethod string, url string, payloadJSON interface{}, targetJSON interface{}) error {
	var b io.ReadWriter
	if payloadJSON != nil {
		b = new(bytes.Buffer)
		if err := json.NewEncoder(b).Encode(payloadJSON); err != nil {
			return errors.Wrap(err, "failed to encode json payload")
		}
	}

	req, err := http.NewRequest(httpMethod, url, b)
	if err != nil {
		return errors.Wrapf(err, "failed to build http request %s %s", httpMethod, url)
	}

	resp, err := m.performRequest(ctx, req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return err
	}

	if targetJSON != nil {
		if err := json.NewDecoder(resp.Body).Decode(targetJSON); err != nil {
			return errors.Wrap(err, "failed to json decode http response")
		}
	}

	return nil
}

func (m *DefaultStationsManager) performRequest(ctx context.Context, req *http.Request) (*http.Response, error) {
	httpClient := m.HttpClient.Provide(ctx, apphttp.ClientOptions{})
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to perform http request")
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return resp, errors.Wrapf(apphttp.ErrUnexpectedStatusCode, "got http status code %v", resp.StatusCode)
	}

	return resp, nil
}

type StationsSearchOptions struct {
	ContractName string
}

// Station is the representation of a station documented here: https://developer.jcdecaux.com/#/opendata/vls?page=dynamic
type Station struct {
	Number              int      `json:"number"`
	Name                string   `json:"name"`
	Address             string   `json:"address"`
	Position            GeoPoint `json:"position"`
	Banking             bool     `json:"banking"`
	Bonus               bool     `json:"bonus"`
	Status              string   `json:"status"`
	ContractName        string   `json:"contract_name"`
	BikeStands          int      `json:"bike_stands"`
	AvailableBikeStands int      `json:"available_bike_stands"`
	AvailableBikes      int      `json:"available_bikes"`
	LastUpdateTimestamp int      `json:"last_update"`
}

func (s Station) LastUpdate() time.Time {
	if s.LastUpdateTimestamp == 0 {
		return time.Time{}
	}
	return time.Unix(int64(s.LastUpdateTimestamp)/1000, 0)
}

type GeoPoint struct {
	Lat float64
	Lng float64
}
