package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/networking/v1"
	"log"
	"net"
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
	_, res := prepareWhitelist(ips, act)
	if !res {
		c.JSON(http.StatusOK, gin.H{
			"code":    40003,
			"message": "IP地址不正确.",
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"message": "操作完成.",
	})
}

// 验证 IP 地址是否为有效的 IPv4 地址
func isValidIPv4(ip string) bool {
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil && parsedIP.To4() != nil
}

// 获得数据库中所有的，跟新增或删除的ip组成新是数组，return出来，给ChangeWhitelist 函数
func prepareWhitelist(ips, act string) (resIps []string, runRes bool) {
	var whiteList []WhiteList

	// 创建一个 map 来去重和验证 IP 地址
	ipMap := make(map[string]bool)
	// 分割传入的 ips 参数
	inputIPs := strings.Split(ips, ",")
	for _, ip := range inputIPs {
		ip = strings.TrimSpace(ip) // 去除空格
		if isValidIPv4(ip) {
			ipMap[ip] = true
		} else {
			log.Println("IP不正确: ", ip)
			return resIps, false
		}
	}

	DB.Find(&whiteList)
	// 提取 IP 地址到一个新的切片
	dbIPs := make([]string, len(whiteList))
	for i, wl := range whiteList {
		dbIPs[i] = wl.IP
	}

	// 根据 act 参数的值添加或删除 IP
	if act == "add" {
		resIps = append(dbIPs, inputIPs...)
	} else if act == "del" {
		resIps = subtractSlices(dbIPs, inputIPs)
	}

	resIps = removeDuplicationMap(resIps)

	fmt.Println(resIps)

	DB.Where("1=1").Delete(&WhiteList{})

	var newWhiteLists []WhiteList
	for _, ip := range resIps {
		newWhiteLists = append(newWhiteLists, WhiteList{IP: ip})
	}

	DB.CreateInBatches(newWhiteLists, len(resIps))

	return resIps, true
}

func subtractSlices(original, toRemove []string) []string {
	// 创建一个 map 标记需要移除的元素
	removeMap := make(map[string]bool)
	for _, item := range toRemove {
		removeMap[item] = true
	}

	// 构建一个新的切片，只包含未被标记移除的元素
	var result []string
	for _, item := range original {
		if !removeMap[item] {
			result = append(result, item)
		}
	}

	return result
}
func removeDuplicationMap(arr []string) []string {
	set := make(map[string]struct{}, len(arr))
	j := 0
	for _, v := range arr {
		_, ok := set[v]
		if ok {
			continue
		}
		set[v] = struct{}{}
		arr[j] = v
		j++
	}

	return arr[:j]
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
