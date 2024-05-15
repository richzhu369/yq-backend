package main

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine) {
	user := router.Group("/api/user")
	{
		user.POST("/login", userLogin)
		user.GET("/info", userInfo)
		user.POST("/logout", userLogout)
		user.POST("/new", userCreate)
		user.DELETE("/delete", userDelete)
	}

	k8sIngress := router.Group("/api/ingress")
	{
		k8sIngress.POST("/whitelist", ingressWhitelist)
		k8sIngress.GET("/whitelistLog", fetchAllWhitelist)
	}
}
