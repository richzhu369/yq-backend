package main

import "gorm.io/gorm"

type WhitelistLog struct {
	gorm.Model
	IP     string `json:"ip"`
	OpUser string `json:"opUser"`
}

type WhiteList struct {
	IP string `json:"ip"`
}
