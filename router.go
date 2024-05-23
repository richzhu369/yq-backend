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
		user.GET("/list", userList)
		user.POST("/reset",userReset)
	}

	k8sIngress := router.Group("/api/ingress")
	{
		k8sIngress.POST("/whitelist", ingressWhitelist)
		k8sIngress.GET("/fetchWhitelist", fetchAllWhitelist)
		k8sIngress.GET("/fetchWhitelistLogs", fetchWhitelistLogs)
	}

	dynadotApi := router.Group("/api/dynadot")
	{
		dynadotApi.GET("/list", dynadotDomainList)
		dynadotApi.GET("/search", dynadotSearchDomain)
		dynadotApi.GET("/buy", dynadotBuyDomain)
	}

	cloudflareAPI := router.Group("/api/cloudflare")
	{
		cloudflareAPI.POST("/createZone", cloudflareCreateZone)
	}
}
