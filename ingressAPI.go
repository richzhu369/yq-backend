package main

import (
	"fmt"
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

	fmt.Println(ips)
	fmt.Println(act)

	// todo: 完善去重功能
	//ChangeWhitelist(ips, act)
	prepareWhitelist(ips, act)

	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"message": "Ingress whitelist update completed.",
	})
}

// 获得数据库中所有的，跟新增或删除的ip组成新是数组，return出来，给ChangeWhitelist 函数
func prepareWhitelist(ips, act string) (resIps []string) {

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

func fetchAllWhitelist(c *gin.Context) {
	var allList []WhiteList
	DB.Find(&allList)

	c.JSON(http.StatusOK, gin.H{
		"code": 20000,
		"data": gin.H{
			"total": len(allList),
			"items": allList,
		},
	})
}
