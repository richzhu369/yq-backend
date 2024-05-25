package main

import (
	"gorm.io/gorm"
)

type SiteInfo struct {
	gorm.Model
	UID                string `json:"UID" gorm:"size:64,unique"`
	Status             string `json:"status"`
	Domain             string `json:"Domain"`
	CwDomain           string `json:"CwDomain"`
	AwsCdnDomain       string `json:"AwsCdnDomain"`
	CfDomain           string `json:"CfDomain"`
	MqTopic            string `json:"MqTopic"`
	MerchantCode       string `json:"MerchantCode"`
	SiteName           string `json:"SiteName"`
	CloudflareDomainID string `json:"CloudflareDomainID"`
	Process            string `json:"Process"`
	CloudflareNS0      string `json:"cloudflareNS0"`
	CloudflareNS1      string `json:"cloudflareNS1"`
}

// xxxxd22.xyz
// xxxxd23.xyz
func updateSiteInfo(siteInfo SiteInfo) {
	// 根据传入的 siteInfo.UID 更新
	DB.Model(&SiteInfo{}).Where("UID = ?", siteInfo.UID).Updates(siteInfo)
}

func insertSiteInfo(siteInfo SiteInfo) {
	// siteInfo 是您要更新的实例，并且它的 ID 已经设置
	DB.Create(&siteInfo)
}

func getSiteInfoByUID(UID string) (siteInfo SiteInfo) {
	DB.Where("UID = ?", UID).First(&siteInfo)
	return siteInfo
}
