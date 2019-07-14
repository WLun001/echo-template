package routes

import (
	"api/controller"
	apiMiddleware "api/controller/middleware"
	"api/repository"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

func APIRoutes(repo repository.Repository, client *mongo.Client) *echo.Echo {
	routes := echo.New()

	ctrl := controller.New(repo, client)
	mw := apiMiddleware.New(client)

	routes.Use(middleware.Logger())
	routes.Use(middleware.Recover())
	routes.Use(middleware.CORS())

	routes.GET("/", ctrl.HelloWorld)
	routes.POST("/", ctrl.Receiver)

	protected := routes.Group("", mw.Authentication)
	r := protected.Group("/protected")
	r.GET("/", ctrl.HelloWorld)

	return routes
}
