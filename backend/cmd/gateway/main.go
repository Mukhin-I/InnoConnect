// @title Innoconnect API
// @version 1.0
// @description Innoconnect backend API
// @host localhost:8080
// @BasePath /
package main

import (
	gatewayapp "innoconnect/internal/gateway/app"
)

func main() {
	// Initialize the gateway server
	gatewayapp.CreateServer()
}
