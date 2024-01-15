package main

import (
	"simple-gateway/https"
	"simple-gateway/infra"
)

func init() {
	infra.InitConfig()
	infra.InitDb()
}
func main() {
	https.InitRoutes()
}
