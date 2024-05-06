package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/flowcontrol"
	"log"
)

var kubeconfig *string
var DB *gorm.DB
var ERR error
var ClientSet *kubernetes.Clientset

func init() {
	// 从命令行参数中获取 kubeconfig 文件路径
	kubeconfig = flag.String("kubeconfig", "", "Path to the kubeconfig file")
	flag.Parse()

	// 初始化数据库
	DB, ERR = gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if ERR != nil {
		log.Fatal(ERR.Error())
	}

	// 自动迁移模式
	ERR = DB.AutoMigrate(&User{},WhitelistLog{})
	if ERR != nil {
		log.Fatal("failed to migrate database: ", ERR)
	}

	// 构建配置
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	config.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(1000, 1000)
	if err != nil {
		panic(err.Error())
	}

	// 创建 Kubernetes 客户端
	ClientSet, ERR = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

}

func main() {

	router := gin.Default()
	router.Use(CORSMiddleware())

	SetupRoutes(router)

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err.Error())
	}

}
