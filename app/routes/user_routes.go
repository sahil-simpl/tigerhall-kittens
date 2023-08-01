package routes

import (
	"github.com/newrelic/go-agent/v3/integrations/nrhttprouter"
	"tigerhallKittens/app/controllers"
	"tigerhallKittens/app/middlewares"
)

func RegisterUserRoutes(router *nrhttprouter.Router) {
	userController := controllers.NewUserController()
	router.POST("/api/v1/users",
		middlewares.ServeV1Endpoint(middlewares.AuthenticateWithHmacDigestMiddleware,
			userController.))
}
