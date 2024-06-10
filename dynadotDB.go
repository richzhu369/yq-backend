package main

import (
	"gorm.io/gorm"
)

type MerchantInfo struct {
	gorm.Model
	Status             string `json:"status"`
	Domain             string `json:"Domain"`
	CwDomain           string `json:"CwDomain"`
	AwsCdnDomain       string `json:"AwsCdnDomain"`
	CfDomain           string `json:"CfDomain"`
	MqTopic            string `json:"MqTopic"`
	MerchantCode       string `json:"MerchantCode"`
	MerchantName       string `json:"MerchantName"`
	CloudflareDomainID string `json:"CloudflareDomainID"`
	Process            string `json:"Process"`
	CloudflareNS0      string `json:"cloudflareNS0"`
	CloudflareNS1      string `json:"cloudflareNS1"`
	AwsSSLArn          string `json:"AwsSSLArn"`
	CnameKey           string `json:"CnameKey"`
	CnameValue         string `json:"CnameValue"`
	CloudFrontID       string `json:"CloudFrontID"`
	CloudfrontRecord   string `json:"CloudfrontRecord"`
}

func updateMerchantInfo(merchant MerchantInfo) {
	// 根据传入的 merchant.UID 更新
	DB.Model(&MerchantInfo{}).Where("merchant_name = ?", merchant.MerchantName).Updates(merchant)
}

func  insertMerchantInfo(merchant MerchantInfo) {
	// siteInfo 是您要更新的实例，并且它的 ID 已经设置
	DB.Create(&merchant)
}

func getMerchantByName(merchantName string) (merchant MerchantInfo) {
	DB.Where("merchant_name = ?", merchantName).First(&merchant)
	return merchant
}
