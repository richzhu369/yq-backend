package main

import (
	"context"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2/admin"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func createTopic(c *gin.Context) {
	domainUID := c.PostForm("domainUID")
	site := getSiteInfoByUID(domainUID)

	replacedStr := strings.ReplaceAll(publicProperty.MQTopics, "NAME", site.SiteName)
	topicsArr := strings.Split(replacedStr, ",")

	log.Println("准备要创建的topics", topicsArr)

	// 设置 RocketMQ 服务器地址
	nameServerAddress := []string{publicProperty.MQServer}

	// 创建 RocketMQ 管理员客户端
	mqAdmin, err := admin.NewAdmin(
		admin.WithResolver(primitive.NewPassthroughResolver(nameServerAddress)),
	)
	if err != nil {
		log.Fatalf("创建管理员客户端失败: %s\n", err)
	}
	defer mqAdmin.Close()

	// 创建 Topic
	for _, v := range topicsArr {
		err := mqAdmin.CreateTopic(
			context.Background(),
			admin.WithTopicCreate(v),
			admin.WithBrokerAddrCreate(publicProperty.MQBroker),
		)
		if err != nil {
			fmt.Println(err.Error())
		}
	}

	// 检查创建topic是否成功
	// 获取 Topic 列表
	topicList, _ := mqAdmin.FetchAllTopicList(context.Background())
	if err != nil {
		log.Fatalf("获取 Topic 列表失败: %sn", err)
	}

	log.Println("当前所有Topic", topicList.TopicList)

	verifyRes := isASubsetOfB(topicsArr, topicList.TopicList)

	if verifyRes {
		c.JSON(http.StatusOK, gin.H{
			"code":    20000,
			"res":     "success",
			"message": "创建topics 成功",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"code":    20000,
			"res":     "failed",
			"message": "创建topics 失败",
		})
	}
}

func isASubsetOfB(listA, listB []string) bool {
	exists := make(map[string]bool)
	for _, b := range listB {
		exists[b] = true
	}
	for _, a := range listA {
		if !exists[a] {
			return false
		}
	}
	return true
}
