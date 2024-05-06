package main

import (
	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/networking/v1"
	"log"
	"net/http"
	"strings"
	"sync"
)

func ingressWhitelist(c *gin.Context) {
	ips := c.PostForm("ips")
	act := c.PostForm("act")

	ChangeWhitelist(ips, act)

	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"message": "Ingress whitelist update completed.",
	})
}

func ChangeWhitelist(ips, act string) {
	ipsToAdd := strings.Split(ips, ",")

	ingressList := GetAllIngress(ClientSet)

	// 使用并发控制
	var wg sync.WaitGroup
	updateChan := make(chan error)

	for _, ingress := range ingressList.Items {
		wg.Add(1)
		go func(ingress v1.Ingress) {
			defer wg.Done()
			if act == "add" {
				err := AddIPsToWhitelist(ClientSet, ingress.Namespace, ingress.Name, ipsToAdd)
				if err != nil {
					updateChan <- err
					return
				}
			} else if act == "del" {
				err := RemoveIPsFromWhitelist(ClientSet, ingress.Namespace, ingress.Name, ipsToAdd)
				if err != nil {
					updateChan <- err
					return
				}
			}

		}(ingress)
	}

	go func() {
		wg.Wait()
		close(updateChan)
	}()

	// 处理并发操作结果
	for err := range updateChan {
		if err != nil {
			log.Fatalf("Failed to update Ingress whitelist: %v", err)
		}
	}

}

func fetchAllWhitelistLog(c *gin.Context) {
	var allLog []WhitelistLog
	DB.Find(&allLog)

	c.JSON(http.StatusOK, gin.H{
		"code": 20000,
		"data": allLog,
	})

}
