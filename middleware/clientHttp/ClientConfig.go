package clientHttp

import (
	"fmt"
	"github.com/spf13/viper"
)

type ClientConfig struct {
	Name    string         `mapstructure:"name"`
	BaseUrl string         `mapstructure:"baseurl"`
	Apis    map[string]API `mapstructure:"apis"`
}

type API struct {
	Name   string `mapstructure:"name"`
	Uri    string `mapstructure:"uri"`
	Method string `mapstructure:"method"`
}

func NewClientConfig(endpoint string) *ClientConfig {
	var cfg ClientConfig
	err := viper.UnmarshalKey(endpoint, &cfg)
	if err != nil {
		fmt.Println("errors when call configs")
		panic(err)
	}
	return &cfg
}
