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
				"res":     "success",
				"domains": res.String(),
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

func dynadotSearchDomain(merchantName, merchantCode string) bool {

	domain := randomString(6) + merchantName + randomString(2) + ".xyz"
	fmt.Println(domain)

	url := "https://api.dynadot.com/api3.json?key=pE8G6Q608b8a6l8x7C6u7oR6fU6V7t8t6Y746g656S7i&command=search&domain0=" + domain + "&show_price=1&currency=USD"
	log.Println(url)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	log.Println(string(body))

	res := gjson.Get(string(body), "SearchResponse.SearchResults.0.Available")
	log.Println(res)
	if res.String() == "yes" {
		insertMerchantInfo(MerchantInfo{
			Domain:       domain,
			MerchantCode: merchantCode,
			MerchantName: merchantName,
			CwDomain:     merchantName + "cw." + domain,
			AwsCdnDomain: merchantName + "ht." + domain,
			CfDomain:     merchantName + "ht1" + domain,
			MqTopic:      merchantName,
			Process:      "搜索是否可被注册",
			Status:       "process",
		})
		upgradeProgress(1, merchantName, "el-icon-check", "primary")
		upgradeProgress(2, merchantName, "el-icon-loading", "primary")
		return true
	} else {
		upgradeProgress(1, merchantName, "el-icon-close", "danger")
		return false

	}
}

func dynadotBuyDomain(merchantName string) bool {
	merchant := getMerchantByName(merchantName)

	url := "https://api.dynadot.com/api3.json?key=pE8G6Q608b8a6l8x7C6u7oR6fU6V7t8t6Y746g656S7i&command=register&domain=" + merchant.Domain + "&duration=1&currency=USD"
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
		merchant.Process = "购买完成"
		updateMerchantInfo(merchant)
		upgradeProgress(2, merchantName, "el-icon-check", "primary")
		upgradeProgress(3, merchantName, "el-icon-loading", "primary")
		return true
	} else {
		updateMerchantInfo(MerchantInfo{
			Process: "购买失败",
			Status: "failed",
		})
		upgradeProgress(2, merchantName, "el-icon-close", "primary")
		return false
	}
}

func dynadotChangeNS(merchantName string) bool {

	site := getMerchantByName(merchantName)

	url := "https://api.dynadot.com/api3.json?key=pE8G6Q608b8a6l8x7C6u7oR6fU6V7t8t6Y746g656S7i&command=set_ns&domain=" + site.Domain + "&ns0=" + site.CloudflareNS0 + "&ns1=" + site.CloudflareNS1

	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("The HTTP request failed with error %s\n", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	log.Println(string(body))

	res := gjson.Get(string(body), "SetNsResponse.Status").String()
	if res == "success" {
		site.Process = "NS服务器更改完成"
		updateMerchantInfo(site)
		upgradeProgress(4, merchantName, "el-icon-check", "primary")
		upgradeProgress(5, merchantName, "el-icon-loading", "primary")
		return true
	} else {
		site.Process = "NS服务器更改失败"
		site.Status = "failed"
		updateMerchantInfo(site)
		upgradeProgress(4, merchantName, "el-icon-close", "primary")
		return false
	}
}
