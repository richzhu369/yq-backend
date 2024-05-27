package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/flowcontrol"
	"log"
	"os"
)

var kubeconfig *string
var DB *gorm.DB
var ERR error
var ClientSet *kubernetes.Clientset

var AwsAK string
var AwsSK string

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
	ERR = DB.AutoMigrate(&User{}, &WhiteList{}, &WhitelistLog{}, &SiteInfo{})
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

	// 获取aws配置
	awsConfigPath := ".config"
	data,err := os.ReadFile(awsConfigPath)
	if err!=nil{
		panic(err.Error())
	}

	jsonStr := string(data)
	AwsAK = gjson.Get(jsonStr,"AWS_ACCESS_KEY_ID").String()
	AwsSK = gjson.Get(jsonStr,"AWS_SECRET_ACCESS_KEY").String()
}

func main() {

	router := gin.Default()
	router.Use(CORSMiddleware())

	SetupRoutes(router)

	if err := router.Run(":8080"); err != nil {
		log.Fatal(err.Error())
	}

}
