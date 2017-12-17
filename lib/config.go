package dementor

import (
	"github.com/kelseyhightower/envconfig"
)

type CommonConf struct {
	Url      string `short:"u" long:"url" description:"Azkaban URL" env:"DEM_URL" default:"http://localhost:80/" envconfig:"DEM_URL"`
	UserName string `short:"U" long:"username" description:"Username" env:"DEM_USERNAME" default:"azkaban" envconfig:"DEM_USERNAME"`
	Password string `short:"P" long:"password" description:"Password" env:"DEM_PASSWORD" default:"azkaban" envconfig:"DEM_PASSWORD"`
}

var Config CommonConf

// Initialize CommonConf mainly for tests
func InitConf() error {
	return envconfig.Process("DEM", &Config)
}
