package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	v1 "k8s.io/api/networking/v1"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

func ingressWhitelist(c *gin.Context) {
	ips := c.PostForm("ips")
	act := c.PostForm("act")
	username, _ := c.Cookie("Username")

	fmt.Println(ips)
	fmt.Println(act)
	fmt.Println(username)

	resIps, resStatus := prepareWhitelist(ips, act, username)
	if !resStatus {
		c.JSON(http.StatusOK, gin.H{
			"code":    40003,
			"message": "IP地址不正确.",
		})

		return
	}

	// 这里判断加白操作 是否执行成功
	ChangeWhitelist(resIps)


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
func prepareWhitelist(ips, act, userName string) (resIps []string, runRes bool) {
	var whiteList []WhiteList
	var whiteListLog WhitelistLog

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

	log.Println("加白的ip为：", resIps)

	DB.Where("1=1").Delete(&WhiteList{})


	var newWhiteLists []WhiteList
	for _, ip := range resIps {
		newWhiteLists = append(newWhiteLists, WhiteList{IP: ip})
	}

	DB.CreateInBatches(newWhiteLists, len(resIps))

	logIps := strings.Join(resIps,",")
	whiteListLog.OpUser = userName
	whiteListLog.IP = logIps
	DB.Create(&whiteListLog)

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

func ChangeWhitelist(ips []string) {

	ingressList := GetAllCWIngress(ClientSet)

	// 使用并发控制
	var wg sync.WaitGroup
	updateChan := make(chan error)
	var failedUpdates []string

	for _, ingress := range ingressList.Items {
		wg.Add(1)
		//sleep减速，太快了 aws的k8s受不了
		time.Sleep(500 * time.Millisecond)
		go func(ingress v1.Ingress) {
			defer wg.Done()
			err := AddIPsToWhitelist(ClientSet, ingress.Namespace, ingress.Name, ips)
			if err != nil {
				updateChan <- err
				return
			}
			updateChan <- nil
		}(ingress)
	}

	go func() {
		wg.Wait()
		close(updateChan)
	}()

	// 处理并发操作结果
	for err := range updateChan {
		if err != nil {
			log.Printf("加白失败: %v", err)
			failedUpdates = append(failedUpdates, err.Error())
		}
	}
	if len(failedUpdates) > 0 {
		log.Println("部分更新未成功:", failedUpdates)
	} else {
		log.Println("所有更新已成功完成")
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

//func fetchWhitelistLogs(c *gin.Context) {
//	var whitelistLogs []WhitelistLog
//	DB.Find(&whitelistLogs)
//
//	c.JSON(http.StatusOK, gin.H{
//		"code": 20000,
//		"data": gin.H{
//			"total": len(whitelistLogs),
//			"items": whitelistLogs,
//		},
//	})
//}
func fetchWhitelistLogs(c *gin.Context) {
	// 分页参数
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	sort := c.DefaultQuery("sort", "+id")

	var whitelistLogs []WhitelistLog
	var total int64

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
	DB.Order(sortOrder).Offset(offset).Limit(limit).Find(&whitelistLogs)
	DB.Model(&WhitelistLog{}).Count(&total) // 计算总数

	// 返回 JSON 响应
	c.JSON(http.StatusOK, gin.H{
		"code": 20000,
		"data": gin.H{
			"total": total,
			"items": whitelistLogs,
		},
	})
}
