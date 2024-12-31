package main

import (
	"example/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	router := gin.Default()

	// MongoDB route
	routes.MongodbUserRoute(router)
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	// MySql route
	routes.MysqlUserRoute(router)
	router.Run("localhost:6000")
}
