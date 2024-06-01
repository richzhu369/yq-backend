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
		user.POST("/reset", userReset)
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
		dynadotApi.PATCH("/changeNS", dynadotChangeNS)
	}

	cloudflareAPI := router.Group("/api/cloudflare")
	{
		cloudflareAPI.POST("/createZone", cloudflareCreateZone)
		cloudflareAPI.GET("/checkZone", cloudflareCheckZone)
		cloudflareAPI.PUT("/createRootRecord", cloudflareCreateRootRecord)
		cloudflareAPI.PUT("/createSSLRecord", cloudflareCreateSSLRecord)
		cloudflareAPI.PUT("/cloudflareCreateCloudfrontRecord", cloudflareCreateCloudfrontRecord)
	}

	awsAPI := router.Group("/api/aws")
	{
		awsAPI.POST("/createCloudfront", createCloudFront)
		awsAPI.POST("/createSSL", createSSL)
		awsAPI.GET("/getSSLVerifyInfo", GetSSLVerifyInfo)
		awsAPI.GET("/getSSLStatus", GetSSLStatus)
		awsAPI.GET("/getCloudFrontDomain", GetCloudFrontDomain)
	}

	rocketmqAPI := router.Group("/api/rocketmq")
	{
		rocketmqAPI.POST("/createTopic", createTopic)
	}

	publicPropertyAPI := router.Group("/api/publicProperty")
	{
		publicPropertyAPI.GET("/get", getPublicProperty)
		publicPropertyAPI.POST("/edit", editPublicProperty)
	}

	etcdAPI := router.Group("/api/etcd")
	{
		etcdAPI.PUT("/create", createETCD)
	}
}
