package routes

import (
	"github.com/maykealisson/fin-news/config"
	"github.com/maykealisson/fin-news/controllers"
)

func HandlerRequests() {

	server := config.SetupGin()

	server.GET("/fin-news/v1/noticias", controllers.BuscarNoticias)

	server.Run(":3001")
}
