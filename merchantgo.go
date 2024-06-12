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

	// 创建前端timeline步骤
	createProgress(merchantName)

	var res bool
	// 1. 查域名
	res = dynadotSearchDomain(merchantName, merchantCode)
	if !res {
		createMerchantFailed(c)
	}
	// 2. 买域名
	res = dynadotBuyDomain(merchantName)
	if !res {
		createMerchantFailed(c)
	}
	// 3. 域名添加到 cloudflare
	res = cloudflareCreateZone(merchantName)
	if !res {
		createMerchantFailed(c)
	}
	// 4. 更改NS到 cloudflare
	res = dynadotChangeNS(merchantName)
	if !res {
		createMerchantFailed(c)
	}
	// 5. 验证cloudflare NS
	for !cloudflareCheckZone(merchantName) {
		log.Println("cloudflare NS验证未通过，正在重试...")
		time.Sleep(5 * time.Second)
	}
	// 6. 创建cname * ，到cf的lb
	res = cloudflareCreateRootRecord(merchantName)
	if !res {
		createMerchantFailed(c)
	}
	// 7. aws中创建SSL
	res = createSSL(merchantName)
	if !res {
		createMerchantFailed(c)
	}
	// 8. 获得SSL验证信息
	res = GetSSLVerifyInfo(merchantName)
	if !res {
		createMerchantFailed(c)
	}
	// 9. 在cloudflare中创建 aws ssl需要的 cname
	res = cloudflareCreateSSLRecord(merchantName)
	if !res {
		createMerchantFailed(c)
	}
	// 10. 在aws中检测ssl的状态，是否通过验证
	for !GetSSLStatus(merchantName) {
		log.Println("aws SSL状态检查未通过，正在重试...")
		time.Sleep(5 * time.Second) // 等待30秒后重试
	}

	// 11. 在aws中创建 cloudfront
	res = createCloudFront(merchantName)
	if !res {
		createMerchantFailed(c)
	}
	// 12. 在cloudflare中创建 cloudfront的cname ht
	res = cloudflareCreateCloudfrontRecord(merchantName)
	if !res {
		createMerchantFailed(c)
	}

	// 13. 创建RocketMQ Topic
	res = createTopic(merchantName)
	if !res {
		createMerchantFailed(c)
	}
	// 14. 创建ETCD配置
	res = createETCD(merchantName)
	if !res {
		createMerchantFailed(c)
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"message": "创建商户成功",
		"res":     "success",
	})
}

func createMerchantFailed(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    40001,
		"message": "创建商户失败",
	})
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
