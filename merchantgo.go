package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

func getAllMerchant(c *gin.Context) {
	var merchants []MerchantInfo
	var total int64

	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	sort := c.DefaultQuery("sort", "+id")

	// 计算分页的偏移量
	offset := (page - 1) * limit

	// 排序处理，根据前端传来的参数决定升序还是降序
	sortOrder := ""
	if sort[0] == '%' { // 假设传来的是 URL 编码后的 '+'
		if sort[1] == '2' {
			sortOrder = "id desc" // 降序
		}
	} else {
		sortOrder = "id" // 默认升序
	}

	// 查询数据库
	DB.Order(sortOrder).Offset(offset).Limit(limit).Find(&merchants)
	DB.Model(&MerchantInfo{}).Count(&total) // 计算总数

	//DB.Find(&merchants)

	c.JSON(http.StatusOK, gin.H{
		"code": 20000,
		"data": gin.H{
			"items": merchants,
			"total": total,
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

// createPipeline 开始创建商户
func createPipeline(merchantName, merchantCode string, merchant MerchantInfo) {
	// 创建前端timeline步骤
	createProgress(merchantName)

	//1. 查域名
	for !dynadotSearchDomain(merchantName, merchantCode) {
		log.Println("dynadot查询域名出错，正在重试...")
		time.Sleep(5 * time.Second)
	}

	// 2. 买域名
	for !dynadotBuyDomain(merchantName) {
		log.Println("dynadot购买域名出错，正在重试...")
		time.Sleep(5 * time.Second)
	}

	// 3. 域名添加到 cloudflare
	for !cloudflareCreateZone(merchantName) {
		log.Println("cloudflare创建zone出错，正在重试...")
		time.Sleep(5 * time.Second)
	}

	// 4. 更改NS到 cloudflare
	for !dynadotChangeNS(merchantName) {
		log.Println("更改NS到cloudflare出错，正在重试...")
		time.Sleep(5 * time.Second)
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
	for !cloudflareCreateRootRecord(merchantName) {
		log.Println("cloudflare创建cname出错，正在重试...")
		time.Sleep(5 * time.Second)
	}

	// 7. aws中创建SSL
	for !createSSL(merchantName) {
		log.Println("aws创建SSL出错，正在重试...")
		time.Sleep(5 * time.Second)
	}

	// 8. 获得SSL验证信息
	for !GetSSLVerifyInfo(merchantName) {
		log.Println("获取aws的ssl验证信息出错，正在重试...")
		time.Sleep(5 * time.Second)
	}
	// 9. 在cloudflare中创建 aws ssl需要的 cname
	for cloudflareCreateSSLRecord(merchantName) {
		log.Println("cloudflare创建SSL cname出错，正在重试...")
		time.Sleep(5 * time.Second)
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
	for !cloudflareCreateCloudfrontRecord(merchantName) {
		log.Println("cloudflare创建cloudfront cname出错，正在重试...")
		time.Sleep(5 * time.Second)
	}

	// 13. 创建RocketMQ Topic
	for !createTopic(merchantName) {
		log.Println("创建RocketMQ Topic出错，正在重试...")
		time.Sleep(5 * time.Second)
	}

	// 14. 创建ETCD配置
	for !createETCD(merchantName) {
		log.Println("创建ETCD配置出错，正在重试...")
		time.Sleep(5 * time.Second)
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
