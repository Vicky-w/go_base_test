package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	/*
		router.GET("/someGet", getting)
		router.POST("/somePost", posting)
		router.PUT("/somePut", putting)
		router.DELETE("/someDelete", deleting)
		router.PATCH("/somePatch", patching)
		router.HEAD("/someHead", head)
		router.OPTIONS("/someOptions", options)
	*/
	r.Run() // listen and serve on 0.0.0.0:8080
}
