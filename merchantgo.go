package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func getAllMerchant(c *gin.Context) {
	var merchants []MerchantInfo

	DB.Find(&merchants)

	c.JSON(http.StatusOK, gin.H{
		"code": 20000,
		"data": gin.H{
			"items": merchants,
		},
	})
}

func createMerchant(c *gin.Context) {
	merchantName := c.PostForm("merchantName")
	merchantCode := c.PostForm("merchantCode")

	var merchant MerchantInfo
	result := DB.Where("merchant_name = ?", merchantName).Find(&merchant)
	if result.RowsAffected > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code":    40001,
			"message": "商户名" + merchantName + "已存在",
			"res":     "failed",
		})
		return
	}

	// 先返回成功，创建任务在后台执行
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"message": "创建商户成功",
		"res":     "success",
	})

	go createPipeline(merchantName, merchantCode, merchant)
}

func createPipeline(merchantName, merchantCode string, merchant MerchantInfo) {
	// 创建前端timeline步骤
	createProgress(merchantName)

	var res bool
	//1. 查域名
	if res = dynadotSearchDomain(merchantName, merchantCode); !res {
		return
	}

	// 2. 买域名
	if res = dynadotBuyDomain(merchantName); !res{
		return
	}

	// 3. 域名添加到 cloudflare
	if res = cloudflareCreateZone(merchantName);!res{
		return
	}

	// 4. 更改NS到 cloudflare
	if res = dynadotChangeNS(merchantName);!res{
		return
	}

	// 5. 验证cloudflare NS
	time.Sleep(10 * time.Second)
	cloudflareForceCheck(merchantName)
	log.Println("执行NS server检测")
	time.Sleep(10 * time.Second)

	for !cloudflareCheckZone(merchantName) {
		log.Println("cloudflare NS验证未通过，正在重试...")
		time.Sleep(5 * time.Second)
	}
	// 6. 创建cname * ，到cf的lb
	if res = cloudflareCreateRootRecord(merchantName);!res{
		return
	}

	// 7. aws中创建SSL
	if res = createSSL(merchantName);!res{
		return
	}

	// 8. 获得SSL验证信息
	for !GetSSLVerifyInfo(merchantName) {
		log.Println("获取aws的ssl验证信息出错，正在重试...")
		time.Sleep(5 * time.Second)
	}
	// 9. 在cloudflare中创建 aws ssl需要的 cname
	if res = cloudflareCreateSSLRecord(merchantName);!res{
		return
	}

	// 10. 在aws中检测ssl的状态，是否通过验证
	for !GetSSLStatus(merchantName) {
		log.Println("aws SSL状态检查未通过，正在重试...")
		time.Sleep(5 * time.Second) // 等待30秒后重试
	}

	// 11. 在aws中创建 cloudfront
	for !createCloudFront(merchantName) {
		log.Println("在aws中创建 cloudfront出错，正在重试...")
		time.Sleep(5 * time.Second)
	}

	// 12. 在cloudflare中创建 cloudfront的cname ht
	if res = cloudflareCreateCloudfrontRecord(merchantName);!res{
		return
	}

	// 13. 创建RocketMQ Topic
	if res = createTopic(merchantName); !res{
		return
	}

	// 14. 创建ETCD配置
	if res = createETCD(merchantName);!res{
		return
	}

	// 更新商户表的 status
	merchant = getMerchantByName(merchantName)
	merchant.Status = "done"
	insertMerchantInfo(merchant)
}

func merchantGetProgress(c *gin.Context) {
	var cp []CreateProgress
	merchantName := c.Query("merchantName")
	fmt.Println(merchantName)

	result := DB.Where("merchant_name = ?", merchantName).Find(&cp)
	if result.RowsAffected == 0 {

		c.JSON(http.StatusOK, gin.H{
			"code":    40001,
			"message": "商户名不存在",
			"data": gin.H{
				"items": cp,
			},
		})
	} else {
		// 如果查询到数据，返回查询结果
		c.JSON(http.StatusOK, gin.H{
			"code":    20000,
			"message": "查询成功",
			"data": gin.H{
				"items": cp,
			},
		})
	}
}
