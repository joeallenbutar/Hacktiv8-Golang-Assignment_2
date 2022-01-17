//Joe Allen Butarbutar (GLNG020ONL003)

package main

import (
	"Assignment-2/controller"
	"Assignment-2/db"
	"fmt"

	"github.com/gin-gonic/gin"
)

func init() {
	fmt.Println("initialize function.")
	db.InitializeDB()
}

func main() {
	route := gin.Default()
	orderRoute := route.Group("/orders")
	{
		orderRoute.POST("/", controller.CreateOrder)
		orderRoute.GET("/", controller.GetOrder)
		orderRoute.PUT("/:id", controller.UpdateOrder)
		orderRoute.DELETE("/:id", controller.DeleteOrder)
	}
	fmt.Println("Server run at port 8080 (http://localhost:8080).")
	route.Run(":8080")
}
