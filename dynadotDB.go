package main

import (
	"gorm.io/gorm"
	"time"
)

type SiteInfo struct {
	gorm.Model
	UID          string    `json:"UID" gorm:"size:64,unique"`
	Status       string    `json:"status"`
	CwDomain     string    `json:"CwDomain"`
	AwsCdnDomain time.Time `json:"AwsCdnDomain"`
	CfDomain     string    `json:"CfDomain"`
	MqTopic      string    `json:"MqTopic"`
	MerchantCode string    `json:"MerchantCode"`
	PlatformName string    `json:"PlatformName"`
}
