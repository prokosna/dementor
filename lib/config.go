package dementor

import (
	"github.com/kelseyhightower/envconfig"
)

type CommonConf struct {
	HTTP struct {
		Url      string `default:"http://localhost:80/" envconfig:"DEM_URL"`
		Insecure bool   `default:"false" envconfig:"DEM_INSECURE"`
	}
	Session struct {
		UserName string `default:"azkaban" envconfig:"DEM_USERNAME"`
		Password string `default:"azkaban" envconfig:"DEM_PASSWORD"`
	}
}

var Config CommonConf

func InitConf() error {
	return envconfig.Process("DEM", &Config)
}
