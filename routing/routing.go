package routing

import (
	"e-commerce_api/controller"

	"github.com/julienschmidt/httprouter"
)

func SetUp(router *httprouter.Router) {
	router.POST("/register", controller.Register)
}
