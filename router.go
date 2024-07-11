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
		k8sIngress.GET("/getStatus", checkIsProcessing)
	}

	dynadotApi := router.Group("/api/dynadot")
	{
		dynadotApi.GET("/list", dynadotDomainList)
	}

	awsAPI := router.Group("/api/aws")
	{
		awsAPI.GET("/getCloudFrontDomain", GetCloudFrontDomain)
	}

	publicPropertyAPI := router.Group("/api/publicProperty")
	{
		publicPropertyAPI.GET("/get", getPublicProperty)
		publicPropertyAPI.POST("/edit", editPublicProperty)
	}

	merchantManagementAPI := router.Group("/api/merchant")
	{
		merchantManagementAPI.GET("/get", getAllMerchant)
		merchantManagementAPI.POST("/create", createMerchant)
		merchantManagementAPI.GET("/createProgress", merchantGetProgress)
		merchantManagementAPI.GET("/getFrontendDomain", merchantGetBindDomain)
		merchantManagementAPI.POST("/bindFrontendDomain", bindFrontendDomain)
		merchantManagementAPI.DELETE("/deleteFrontendDomain", deleteFrontendDomain)
	}
}
