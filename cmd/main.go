package main

import (
	"fmt"
	"userportal/internal/app/config"
	"userportal/internal/app/database"
	"userportal/internal/app/router"

	"github.com/gin-gonic/gin"
)

func main() {
	const configFileRelativePath = "./internal/app/config"
	// Get Server Configuration
	serverConf := config.ReadConfig(configFileRelativePath).Server
	host := serverConf.Host
	port := serverConf.Port

	// Run DB migration
	database.RunDatabaseMigrations()

	// Run gin-engine/app
	engine := gin.Default()
	router.Router(engine)

	engine.Run(fmt.Sprintf("%s:%s", host, port))
}
