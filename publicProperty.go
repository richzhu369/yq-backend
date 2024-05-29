package main

import "github.com/gin-gonic/gin"

func getPublicProperty(c *gin.Context) {
	c.JSON(200, gin.H{
		"code": 20000,
		"data": gin.H{
			"items": publicProperty,
		},
	})
}
