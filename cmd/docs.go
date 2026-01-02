package main

//go:generate sh -c "swag init -g main.go -o ../docs -ot yaml --parseDependency --parseInternal && mv ../docs/swagger.yaml ../docs/openapi.yaml"

// @title Currency Exchange API
// @version 1.0.0
// @description HTTP API for currencies, exchange rates, and conversions.
// @host localhost:8080
// @BasePath /
