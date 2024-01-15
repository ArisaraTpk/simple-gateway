package infra

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

func InitConfig() {
	viper.AddConfigPath("./configs")
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("errors viper %+v", err)
	}
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
}
