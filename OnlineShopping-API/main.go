package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"github.com/melardev/GoGonicEcommerceApi/controllers"
	"github.com/melardev/GoGonicEcommerceApi/infrastructure"
	"github.com/melardev/GoGonicEcommerceApi/middlewares"
	"github.com/melardev/GoGonicEcommerceApi/models"
	"github.com/melardev/GoGonicEcommerceApi/seeds"
	"os"
)


func main() {

	e := godotenv.Load() //Load .env file
	if e != nil {
		fmt.Print(e)
	}
	println(os.Getenv("DB_DIALECT"))

	database := infrastructure.OpenDbConnection()

	defer database.Close()
	args := os.Args
	if len(args) > 1 {
		first := args[1]
		second := ""
		if len(args) > 2 {
			second = args[2]
		}

		if first == "create" {
			create(database)
		} else if first == "seed" {
			seeds.Seed()
			os.Exit(0)
		} else if first == "migrate" {
			migrate(database)
		}

		if second == "seed" {
			seeds.Seed()
			os.Exit(0)
		} else if first == "migrate" {
			migrate(database)
		}

		if first != "" && second == "" {
			os.Exit(0)
		}
	}

	migrate(database)

	// gin.New() - new gin Instance with no middlewares
	// goGonicEngine.Use(gin.Logger())
	// goGonicEngine.Use(gin.Recovery())
	goGonicEngine := gin.Default() // gin with the Logger and Recovery Middlewares attached
	// Allow all Origins
	goGonicEngine.Use(cors.Default())

	goGonicEngine.Use(middlewares.Benchmark())

	// goGonicEngine.Use(middlewares.Cors())

	goGonicEngine.Use(middlewares.UserLoaderMiddleware())
	goGonicEngine.Static("/static", "./static")
	apiRouteGroup := goGonicEngine.Group("/api")

	controllers.RegisterUserRoutes(apiRouteGroup.Group("/users"))
	controllers.RegisterProductRoutes(apiRouteGroup.Group("/items"))
	controllers.RegisterCartRoutes(apiRouteGroup.Group("/carts"))
	controllers.RegisterOrderRoutes(apiRouteGroup.Group("/orders"))

	goGonicEngine.Run(":8080") // listen and serve on 0.0.0.0:8080
}
