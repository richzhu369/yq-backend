package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"strings"
)

func cloudflareCreateZone(c *gin.Context) {

	url := "https://api.cloudflare.com/client/v4/zones"

	domainUID := c.PostForm("domainUID")
	site := getSiteInfoByUID(domainUID)

	payload := strings.NewReader("{\n  \"account\": {\n    \"id\": \"08658db65e224f04a7315d0e4e55ec89\"\n  },\n  \"name\": \"" + site.Domain + "\",\n  \"type\": \"full\"\n}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer Y6U_-NFZ-ww7xeybO3WmZqeJesj7GAkoWx4d9rL_")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	getRes := gjson.Get(string(body), "success")
	fmt.Println(gjson.Get(string(body), "result.name_servers.0").String())
	fmt.Println(gjson.Get(string(body), "result.name_servers.1").String())
	if getRes.String() == "true" {
		// 更新域名信息
		site.CloudflareNS0 = gjson.Get(string(body), "result.name_servers.0").String()
		site.CloudflareNS1 = gjson.Get(string(body), "result.name_servers.1").String()
		site.CloudflareDomainID = gjson.Get(string(body), "result.id").String()
		site.Process = "添加域名到cloudflare"
		site.Status = "addDomainToCloudflare"
		updateSiteInfo(site)

		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"data": gin.H{
				"res":          "success",
				"message":      "创建站点成功",
				"cfSiteStatus": gjson.Get(string(body), "result.status").String(),
				"domainId":     site.CloudflareDomainID,
				"domainUDI":    site.UID,
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"res":     "failed",
			"code":    40001,
			"message": "创建站点失败",
			"reason":  gjson.Get(string(body), "errors.0.message").String(),
		})
	}
}

func cloudflareCheckZone(c *gin.Context) {

	domainUID := c.PostForm("domainUID")
	site := getSiteInfoByUID(domainUID)

	url := "https://api.cloudflare.com/client/v4/zones/" + site.CloudflareDomainID

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer Y6U_-NFZ-ww7xeybO3WmZqeJesj7GAkoWx4d9rL_")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))
	getRes := gjson.Get(string(body), "result.status")
	if getRes.String() == "active" {
		// 更新域名信息
		site.Process = "检车站点是否pending"
		site.Status = "checkSiteStatus"
		updateSiteInfo(site)

		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"data": gin.H{
				"res":          "failed",
				"message":      "站点激活成功",
				"cfSiteStatus": gjson.Get(string(body), "result.status").String(),
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"data": gin.H{
				"res":     "failed",
				"message": "站点pending中，请等待",
				"error1":  gjson.Get(string(body), "errors").String(),
				"error2":  gjson.Get(string(body), "messages").String(),
			},
		})
	}

}

// 手动再次触发NS检测，暂时无用
func cloudflareForceCheck(domainUID string) {
	site := getSiteInfoByUID(domainUID)

	url := "https://api.cloudflare.com/client/v4/zones/" + site.CloudflareDomainID + "/activation_check"

	req, _ := http.NewRequest("PUT", url, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer Y6U_-NFZ-ww7xeybO3WmZqeJesj7GAkoWx4d9rL_")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

}

func cloudflareCreateRootRecord(c *gin.Context) {
	domainUID := c.PostForm("domainUID")
	site := getSiteInfoByUID(domainUID)

	url := "https://api.cloudflare.com/client/v4/zones/" + site.CloudflareDomainID + "/dns_records"

	payload := strings.NewReader("{\n  \"content\": \"af400c4b64edf4620810c92cd5dd5d82-499673083.sa-east-1.elb.amazonaws.com\",\n  \"name\": \"*\",\n  \"proxied\": true,\n  \"type\": \"CNAME\",\n  \"comment\": \"由yq-devops平台创建\",\n \"ttl\": 60\n}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer Y6U_-NFZ-ww7xeybO3WmZqeJesj7GAkoWx4d9rL_")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))
	getRes := gjson.Get(string(body), "success")
	if getRes.String() == "true" {
		site.Process = "新建CNAME记录"
		site.Status = "createCNAMEtoCloudflare"
		updateSiteInfo(site)

		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"data": gin.H{
				"res":     "success",
				"message": "新建CNAME记录成功",
			},
		})
	} else {
		site.Process = "新建CNAME记录失败"
		site.Status = "failedCreateCNAMEtoCloudflare"
		updateSiteInfo(site)
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"data": gin.H{
				"message": "新建CNAME记录失败",
				"res":     "failed",
				"error":   gjson.Get(string(body), "errors").String(),
			},
		})
	}
}


func cloudflareCreateSSLRecord(c *gin.Context) {
	domainUID := c.PostForm("domainUID")
	site := getSiteInfoByUID(domainUID)

	url := "https://api.cloudflare.com/client/v4/zones/" + site.CloudflareDomainID + "/dns_records"

	payload := strings.NewReader("{\n  \"content\": \""+site.CnameValue+"\",\n  \"name\": \""+site.CnameKey+"\",\n  \"proxied\": false,\n  \"type\": \"CNAME\",\n  \"comment\": \"由yq-devops平台创建\",\n \"ttl\": 60\n}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer Y6U_-NFZ-ww7xeybO3WmZqeJesj7GAkoWx4d9rL_")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))
	getRes := gjson.Get(string(body), "success")
	if getRes.String() == "true" {
		site.Process = "新建SSL CNAME记录"
		site.Status = "create SSL CNAME to Cloudflare"
		updateSiteInfo(site)

		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"data": gin.H{
				"res":     "success",
				"message": "新建CNAME记录成功",
			},
		})
	} else {
		site.Process = "新建CNAME记录失败"
		site.Status = "failedCreateCNAMEtoCloudflare"
		updateSiteInfo(site)
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"data": gin.H{
				"message": "新建CNAME记录失败",
				"res":     "failed",
				"error":   gjson.Get(string(body), "errors").String(),
			},
		})
	}
}

func cloudflareCreateCloudfrontRecord(c *gin.Context) {
	domainUID := c.PostForm("domainUID")
	site := getSiteInfoByUID(domainUID)

	url := "https://api.cloudflare.com/client/v4/zones/" + site.CloudflareDomainID + "/dns_records"

	payload := strings.NewReader("{\n  \"content\": \""+site.CloudfrontRecord+"\",\n  \"name\": \""+site.AwsCdnDomain+"\",\n  \"proxied\": false,\n  \"type\": \"CNAME\",\n  \"comment\": \"由yq-devops平台创建\",\n \"ttl\": 60\n}")

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer Y6U_-NFZ-ww7xeybO3WmZqeJesj7GAkoWx4d9rL_")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(string(body))
	getRes := gjson.Get(string(body), "success")
	if getRes.String() == "true" {
		site.Process = "新建CloudFront CNAME记录"
		site.Status = "create CloudFront CNAME to Cloudflare"
		updateSiteInfo(site)

		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"data": gin.H{
				"res":     "success",
				"message": "新建CloudFront CNAME记录成功",
			},
		})
	} else {
		site.Process = "新建CloudFront CNAME记录失败"
		site.Status = "failedCreateCloudFront CNAMEtoCloudflare"
		updateSiteInfo(site)
		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"data": gin.H{
				"message": "新建CloudFront CNAME记录失败",
				"res":     "failed",
				"error":   gjson.Get(string(body), "errors").String(),
			},
		})
	}
}