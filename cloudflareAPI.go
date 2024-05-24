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
	// 处理创建Cloudflare Zone的逻辑
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
	if getRes.String() == "true" {
		updateSiteInfo(SiteInfo{
			Process:            "在cloudflare中创建站点",
			Status:             gjson.Get(string(body), "result.name_servers.0").String() + "," +gjson.Get(string(body), "result.name_servers.1").String(),
			CloudflareDomainID: gjson.Get(string(body), "result.id").String(),
		})
		c.JSON(http.StatusOK, gin.H{
			"code":    20000,
			"message": "创建站点成功",
			"data": gin.H{
				"cfSiteStatus": gjson.Get(string(body), "result.status").String(),
				"domainId":     site.CloudflareDomainID,
				"domainUDI":    site.UID,
			},
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    40001,
			"message": "创建站点失败",
			"reason":  gjson.Get(string(body), "errors.0.message").String(),
		})
	}
	fmt.Println(string(body))
}
