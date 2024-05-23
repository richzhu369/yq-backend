package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
)

func cloudflareCreateZone(c *gin.Context) {
	// 处理创建Cloudflare Zone的逻辑
	url := "https://api.cloudflare.com/client/v4/zones"
	domain := c.PostForm("domain")

	payload := strings.NewReader(fmt.Sprintf("{n  \"account\": {\n    \"id\": \"023e105f4ecef8ad9ca31a8372d0c353\"\n  },\n  \"name\": \"%s\",n  \"type\": \"full\"n}", domain))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer undefined")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fmt.Println(res)
	fmt.Println(string(body))
}