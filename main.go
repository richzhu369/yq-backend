package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/flowcontrol"
	"net/http"
)

var kubeconfig *string

func init() {
	// 从命令行参数中获取 kubeconfig 文件路径
	kubeconfig = flag.String("kubeconfig", "", "Path to the kubeconfig file")
	flag.Parse()
}

func main() {

	// 构建配置
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	config.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(1000, 1000)
	if err != nil {
		panic(err.Error())
	}

	// 创建 Kubernetes 客户端
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	router := gin.Default()
	router.Use(CORSMiddleware())

	router.POST("/api/ingress/whitelist", func(c *gin.Context) {
		ips := c.PostForm("ips")
		act := c.PostForm("act")

		ChangeWhitelist(ips, act, clientSet)
	})

	router.POST("/api/user/login", func(c *gin.Context) {
		var loginReq LoginRequest
		if err := c.ShouldBindJSON(&loginReq); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 现在可以使用 loginReq.Username 和 loginReq.Password 获取 JSON 中的值
		fmt.Printf("Username: %s, Password: %sn", loginReq.Username, loginReq.Password)

		// 进行登录验证逻辑...

		c.JSON(http.StatusOK, gin.H{
			"code": 20000,
			"data": gin.H{
				"token": "admin-token",
			},
		})
	})

	router.GET("/api/user/info", func(c *gin.Context) {
		token:=c.Query("token")
		fmt.Println(token)

		c.JSON(http.StatusOK, gin.H{
			"code":20000,
			"data": gin.H{
				"role": []string{"admin"},
				"introduction": "Hello iam a admin",
				"avatar": "https://wpimg.wallstcn.com/f778738c-e4f8-4870-b634-56703b4acafe.gif",
				"name": "Super Admin YQ",
			},
		})
	})

	router.Run(":8080")
}
