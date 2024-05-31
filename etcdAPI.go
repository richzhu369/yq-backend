package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func createETCD(c *gin.Context) {
	domainUID := c.PostForm("domainUID")
	site := getSiteInfoByUID(domainUID)

	// 假设文件路径是 'config.toml'
	filePath := "./etcdConfigure"

	// 设置 ETCD 键值对
	key := "/bs/" + site.SiteName + ".toml"

	// 读取文件内容
	value, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("读取文件失败: %v", err)
	}

	replaceStr := strings.ReplaceAll(string(value), "SITE_NAME", site.SiteName)
	replaceStr = strings.ReplaceAll(replaceStr, "CW_DOMAIN", site.CwDomain)
	replaceStr = strings.ReplaceAll(replaceStr, "REDIS_NUM", strconv.Itoa(int(publicProperty.RedisDBNumber)))

	fmt.Println(replaceStr)

	// 创建 ETCD 客户端
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{publicProperty.ETCDServer}, // 替换为实际的 ETCD 服务器地址
		DialTimeout: 10 * time.Second,
	})
	if err != nil {
		log.Fatalf("创建 ETCD 客户端失败: %v", err)
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	_, err = cli.Put(ctx, key, replaceStr)
	cancel()
	if err != nil {
		log.Println("设置 ETCD 键值对失败：" + key)
		c.JSON(http.StatusOK, gin.H{
			"code":    20000,
			"message": "设置 ETCD 键值对失败：" + key,
			"error":   err.Error(),
			"res":     false,
		})
		return
	}

	publicProperty.RedisDBNumber += 1
	fmt.Println(publicProperty.RedisDBNumber)
	DB.Model(&PublicProperty{}).Where("id = 1").Update("redis_db_number", publicProperty.RedisDBNumber)
	//DB.Model(&PublicProperty{}).Update("RedisDBNumber", publicProperty.RedisDBNumber)

	log.Printf("键 %s 创建成功", key)
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"message": "设置 ETCD 键值对成功:" + key,
		"res":     true,
	})
}
