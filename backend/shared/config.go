package shared

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

const envConfigPrefix = "VELIB_TOULOUSE_CARTE"

type AppConfig struct {
	JcdecauxApiKey string `split_words:"true",required:"true"`
}

func ConfigFromEnvVars() (*AppConfig, error) {
	config := &AppConfig{}
	if err := envconfig.Process(envConfigPrefix, config); err != nil {
		return nil, errors.Wrap(err, "failed to parse env vars")
	}

	return config, nil
}
