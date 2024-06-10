package main

import "time"

type CreateProgress struct {
	ID           uint      `gorm:"primaryKey"`
	MerchantName string    `json:"MerchantName"`
	Content      string    `json:"Content"`
	TimeStamp    time.Time `json:"TimeStamp"`
	Size         string    `json:"Size"`
	Type         string    `json:"Type"`
	Icon         string    `json:"Icon"`
	Color string `json:"Color"`
	Step  int    `json:"Step"`
}

// todo: 整理所有流程
func createProgress(merchantName string) {
	// 域名是否可被注册
	insertToProgress(merchantName, "el-icon-loading", "域名是否可被注册", "large", "primary", 1)
	// 域名购买
	insertToProgress(merchantName, "el-icon-more", "域名购买", "large", "primary", 2)
	// 添加到cloudflare
	insertToProgress(merchantName, "el-icon-more", "添加到cloudflare", "large", "primary", 3)
	// 在dynadot更改NS
	insertToProgress(merchantName, "el-icon-more", "更改NS服务器到cf", "large", "primary", 4)
	// 验证cloudflare域名是否添加成功
	insertToProgress(merchantName, "el-icon-more", "验证NS", "large", "primary", 5)
	// 解析CNAME * 到aws中的 lb上
	insertToProgress(merchantName, "el-icon-more", "解析*到aws的lb", "large", "primary", 6)
	// 创建aws ssl
	insertToProgress(merchantName, "el-icon-more", "创建aws ssl", "large", "primary", 7)
	// 获取aws ssl 的验证解析
	insertToProgress(merchantName, "el-icon-more", "获取aws ssl 的验证解析", "large", "primary", 8)
	// 在cloudflare中创建 aws ssl需要的 cname
	insertToProgress(merchantName, "el-icon-more", "验证 ssl cname", "large", "primary", 9)
	// 在aws中查看 ssl的状态，是否通过验证
	insertToProgress(merchantName, "el-icon-more", "验证 ssl cname", "large", "primary", 10)


	// 创建cloudfront
	insertToProgress(merchantName, "el-icon-more", "创建cloudfront", "large", "primary", 11)
	// 解析ht到aws cloudfront
	insertToProgress(merchantName, "el-icon-more", "解析ht1到aws cf", "large", "primary", 12)
	// 创建RocketMQ Topic
	insertToProgress(merchantName, "el-icon-more", "创建RocketMQ Topic", "large", "primary", 13)
	// 创建ETCD
	insertToProgress(merchantName, "el-icon-more", "创建ETCD", "large", "primary", 14)


	// todo: 增加ht1 域名解析 到流程中, 看到老站也没解析，应该是被 * 代替了

}

func insertToProgress(merchantName, cpIcon, cpContent, cpSize, cpType string, Step int) {
	var createProgress CreateProgress

	createProgress.TimeStamp = time.Time{}

	createProgress.MerchantName = merchantName
	createProgress.Icon = cpIcon
	createProgress.Content = cpContent
	createProgress.Size = cpSize
	createProgress.Type = cpType
	DB.Create(&createProgress)
}

func upgradeProgress(stepNum int ,merchantName, cpIcon, cpType string) {
	var createProgress CreateProgress
	createProgress.Icon = cpIcon
	createProgress.Type = cpType
	DB.Updates(&createProgress).Where("merchant_name = ? AND step = ?", merchantName, stepNum)
}
