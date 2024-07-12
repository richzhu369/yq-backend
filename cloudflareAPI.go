package main

import (
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"net/http"
	"strings"
)

func cloudflareCreateZone(merchantName string) bool {

	url := "https://api.cloudflare.com/client/v4/zones"

	site := getMerchantByName(merchantName)

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
		updateMerchantInfo(site)
		upgradeProgress(3, merchantName, "el-icon-check", "primary")
		upgradeProgress(4, merchantName, "el-icon-loading", "primary")
		return true
	} else {
		upgradeProgress(3, merchantName, "el-icon-close", "primary")
		site.Status = "failed"
		updateMerchantInfo(site)
		return false
	}
}

func cloudflareCheckZone(merchantName string) bool {

	site := getMerchantByName(merchantName)

	url := "https://api.cloudflare.com/client/v4/zones/" + site.CloudflareDomainID

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("http.NewRequest error:", err)
		return false
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer Y6U_-NFZ-ww7xeybO3WmZqeJesj7GAkoWx4d9rL_")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("http.DefaultClient.Do error:", err)
		return false
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("io.ReadAll error:", err)
		return false
	}

	fmt.Println(string(body))
	getRes := gjson.Get(string(body), "result.status")
	if getRes.String() == "active" {
		// 更新域名信息
		site.Process = "检车站点是否pending"
		updateMerchantInfo(site)
		upgradeProgress(5, merchantName, "el-icon-check", "primary")
		upgradeProgress(6, merchantName, "el-icon-loading", "primary")
		return true
	}
	return false
}

// 手动再次触发NS检测，暂时无用
func cloudflareForceCheck(merchantName string) {
	site := getMerchantByName(merchantName)

	url := "https://api.cloudflare.com/client/v4/zones/" + site.CloudflareDomainID + "/activation_check"

	req, _ := http.NewRequest("PUT", url, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer Y6U_-NFZ-ww7xeybO3WmZqeJesj7GAkoWx4d9rL_")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()

}

// cloudflareCreateRootRecord 创建CNAME记录* 到 aws的lb
func cloudflareCreateRootRecord(merchantName string) bool {

	site := getMerchantByName(merchantName)

	url := "https://api.cloudflare.com/client/v4/zones/" + site.CloudflareDomainID + "/dns_records"
	log.Println("6 cloudflareCreateRootRecord 访问URL：", url)

	payload := strings.NewReader("{\n  \"content\": \"af400c4b64edf4620810c92cd5dd5d82-499673083.sa-east-1.elb.amazonaws.com\",\n  \"name\": \"*\",\n  \"proxied\": true,\n  \"type\": \"CNAME\",\n  \"comment\": \"由yq-devops平台创建\",\n \"ttl\": 60\n}")
	log.Println("6 cloudflareCreateRootRecord payload：", payload)

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
		updateMerchantInfo(site)
		upgradeProgress(6, merchantName, "el-icon-check", "primary")
		upgradeProgress(7, merchantName, "el-icon-loading", "primary")
		return true
	} else {
		site.Process = "新建CNAME记录失败"
		site.Status = "failed"
		updateMerchantInfo(site)
		upgradeProgress(6, merchantName, "el-icon-close", "primary")
		return false
	}
}

func cloudflareCreateSSLRecord(merchantName string) bool {
	site := getMerchantByName(merchantName)

	url := "https://api.cloudflare.com/client/v4/zones/" + site.CloudflareDomainID + "/dns_records"
	log.Println("9 访问URL：", url)

	payload := strings.NewReader("{\n  \"content\": \"" + site.CnameValue + "\",\n  \"name\": \"" + site.CnameKey + "\",\n  \"proxied\": false,\n  \"type\": \"CNAME\",\n  \"comment\": \"由yq-devops平台创建\",\n \"ttl\": 60\n}")
	log.Println("9 Payload：", payload)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		log.Println("9 http.NewRequest：", err)
		return false
	}
	log.Println("正在创建CNAME记录：", site.CnameKey, site.CnameValue)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer Y6U_-NFZ-ww7xeybO3WmZqeJesj7GAkoWx4d9rL_")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("9 http.DefaultClient：", err)
		return false
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("9 io.ReadAll：", err)
		return false
	}

	log.Println(string(body))
	getRes := gjson.Get(string(body), "success")
	if getRes.String() == "true" {
		site.Process = "新建SSL CNAME记录成功"
		updateMerchantInfo(site)
		upgradeProgress(9, merchantName, "el-icon-check", "primary")
		upgradeProgress(10, merchantName, "el-icon-loading", "primary")
		return true
	} else {
		site.Process = "新建CNAME记录失败"
		site.Status = "failed"
		updateMerchantInfo(site)
		upgradeProgress(9, merchantName, "el-icon-close", "primary")
		return false
	}
}

func cloudflareCreateCloudfrontRecord(merchantName string) bool {

	site := getMerchantByName(merchantName)

	url := "https://api.cloudflare.com/client/v4/zones/" + site.CloudflareDomainID + "/dns_records"
	log.Println("url:", url)

	payload := strings.NewReader("{\n  \"content\": \"" + site.CloudfrontRecord + "\",\n  \"name\": \"" + site.AwsCdnDomain + "\",\n  \"proxied\": false,\n  \"type\": \"CNAME\",\n  \"comment\": \"由yq-devops平台创建\",\n \"ttl\": 60\n}")
	log.Println("payload: ", payload)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		log.Println("http.NewRequest: ", err)
		return false
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer Y6U_-NFZ-ww7xeybO3WmZqeJesj7GAkoWx4d9rL_")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("http.DefaultClient.Do: ", err)
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Println("io.ReadAll: ", err)
		return false
	}

	log.Println(string(body))
	getRes := gjson.Get(string(body), "success")
	if getRes.String() == "true" {
		site.Process = "新建CloudFront CNAME记录"
		updateMerchantInfo(site)
		upgradeProgress(12, merchantName, "el-icon-check", "primary")
		upgradeProgress(13, merchantName, "el-icon-check", "primary")
		return true
	} else {
		site.Process = "新建CloudFront CNAME记录失败"
		site.Status = "failed"
		updateMerchantInfo(site)
		upgradeProgress(12, merchantName, "el-icon-close", "primary")
		log.Println("新建CloudFront CNAME记录失败,gerRest: ", getRes.String())
		return false
	}
}
