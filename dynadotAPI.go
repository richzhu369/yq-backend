package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"math/rand"
	"net/http"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyz0123456789"
)

func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func dynadotDomainList(c *gin.Context) {
	url := "https://api.dynadot.com/api3.json?key=pE8G6Q608b8a6l8x7C6u7oR6fU6V7t8t6Y746g656S7i&command=list_domain"

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	res := gjson.Get(string(body), "ListDomainInfoResponse.Status")
	if res.String() == "success" {
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"data": gin.H{
				"res": "success",
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 40001,
			"data": gin.H{
				"res": "failed",
			},
		})
	}

}

func dynadotSearchDomain(c *gin.Context) {

	merchantCode := c.PostForm("MerchantCode")
	siteName := c.PostForm("SiteName")

	fmt.Println(merchantCode)
	fmt.Println(siteName)

	domain := randomString(6) + siteName + randomString(2) + ".xyz"
	fmt.Println(domain)

	url := "https://api.dynadot.com/api3.json?key=pE8G6Q608b8a6l8x7C6u7oR6fU6V7t8t6Y746g656S7i&command=search&domain0=" + domain + "&show_price=1&currency=USD"
	println(url)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	fmt.Println(string(body))

	res := gjson.Get(string(body), "SearchResponse.SearchResults.0.Available")
	fmt.Println(res)
	if res.String() == "yes" {
		domainUID := randomString(24)
		insertSiteInfo(SiteInfo{
			Domain:       domain,
			MerchantCode: merchantCode,
			SiteName:     siteName,
			CwDomain:     siteName + "cw." + domain,
			AwsCdnDomain: siteName + "ht." + domain,
			CfDomain:     siteName + "ht1" + domain,
			MqTopic:      siteName,
			UID:          domainUID,
			Process:      "搜索是否可被注册",
			Status:       "searching",
		})
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"data": gin.H{
				"message":   "域名可被注册",
				"res":       "success",
				"domainUID": domainUID,
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 40001,
			"data": gin.H{
				"res":     "failed",
				"message": "域名不可被注册",
			},
		})
	}

}

func dynadotBuyDomain(c *gin.Context) {
	domainUID := c.PostForm("domainUID")
	site := getSiteInfoByUID(domainUID)

	url := "https://api.dynadot.com/api3.json?key=pE8G6Q608b8a6l8x7C6u7oR6fU6V7t8t6Y746g656S7i&command=register&domain=" + site.Domain + "&duration=1&currency=USD"
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	log.Println(string(body))

	res := gjson.Get(string(body), "RegisterResponse.Status")
	log.Println(res)

	if res.String() == "success" {
		//todo: 这个地方要更新 process字段
		updateSiteInfo(SiteInfo{
			UID:     domainUID,
			Process: "购买完成",
			Status:  "purchased",
		})
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"data": gin.H{
				"message":   "域名购买成功",
				"res":       "success",
				"domainUDI": site.UID,
			},
		})
	} else {
		updateSiteInfo(SiteInfo{
			UID:     domainUID,
			Process: "购买失败",
			Status:  "purchaseFailed",
		})
		c.JSON(http.StatusOK, gin.H{
			"code": 40001,
			"data": gin.H{
				"message": "域名购买失败",
				"res":     "failed",
				"reason":  res.String(),
			},
		})
	}
}


// todo: 完善这个功能，更改NS服务器
func dynadotChangeNS(c *gin.Context) {
	domainUID := c.PostForm("domainUID")
	site := getSiteInfoByUID(domainUID)

	url := "https://api.dynadot.com/api3.json?key=pE8G6Q608b8a6l8x7C6u7oR6fU6V7t8t6Y746g656S7i&command=register&domain=" + site.Domain + "&duration=1&currency=USD"
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	log.Println(string(body))
}
