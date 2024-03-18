package main

import (
	"assignment2/controllers"
	"assignment2/database"

	"github.com/gin-gonic/gin"
)

func main() {
	database.Connect()
	r := gin.Default()

	r.POST("/items", controllers.CreateItem)
	r.POST("/orders", controllers.CreateOrder)
	r.PUT("/orders/:id", controllers.UpdateOrder)
	r.DELETE("/orders/:id", controllers.DeleteOrder)
	r.GET("/orders", controllers.GetOrders)

	err := r.Run(":3000")
	if err != nil {
		panic(err)
	}
}
