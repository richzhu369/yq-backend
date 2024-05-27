package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func createSSL(c *gin.Context) {
	domainUID := c.PostForm("domainUID")
	site := getSiteInfoByUID(domainUID)

	// 加载 AWS 配置
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(AwsAK, AwsSK, "")),
	)
	if err != nil {
		log.Printf("Unable to load SDK config, %v", err)
	}

	// 创建 ACM 客户端
	client := acm.NewFromConfig(cfg)

	// 创建 SSL/TLS 证书请求
	input := &acm.RequestCertificateInput{
		DomainName:       aws.String("*." + site.Domain), // 设置通配符域名
		ValidationMethod: types.ValidationMethodDns,      // 使用 DNS 验证
	}

	// 请求 SSL/TLS 证书
	result, err := client.RequestCertificate(context.TODO(), input)
	if err != nil {
		log.Printf("Failed to request certificate, %v", err)
	}

	log.Printf("Certificate ARN: %sn", aws.ToString(result.CertificateArn))

	site.Process = "Certificate Requested"
	site.Status = "Pending verification by cname on aws"
	site.AwsSSLArn = aws.ToString(result.CertificateArn)
	updateSiteInfo(site)

	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"message": "证书创建成功",
	})
}

func GetSSLVerifyInfo(c *gin.Context) {
	domainUID := c.PostForm("domainUID")
	site := getSiteInfoByUID(domainUID)

	// 加载 AWS 配置
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("sa-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(AwsAK, AwsSK, "")),
	)
	if err != nil {
		log.Printf("Unable to load SDK config, %v", err)
	}

	// 创建 ACM 客户端
	client := acm.NewFromConfig(cfg)

	// 替换为您的证书 ARN
	certificateArn := site.AwsSSLArn

	// 调用 DescribeCertificate 获取证书详情
	describeCertInput := &acm.DescribeCertificateInput{
		CertificateArn: aws.String(certificateArn),
	}

	certDetail, err := client.DescribeCertificate(context.TODO(), describeCertInput)
	if err != nil {
		log.Printf("Failed to describe certificate, %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    20000,
			"message": "Cname验证信息获取失败",
			"res":     "failed",
			"error":   err.Error(),
		})
		return
	}

	site.Process = "等待Cname验证"
	site.Status = "waiting cname verification"

	// 输出 CNAME 验证信息
	for _, option := range certDetail.Certificate.DomainValidationOptions {
		log.Printf("DomainName: %sn", *option.DomainName)
		if option.ResourceRecord != nil {
			site.CnameKey = *option.ResourceRecord.Name
			site.CnameValue = *option.ResourceRecord.Value
			updateSiteInfo(site)
		}

	}

	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"message": "Cname验证信息获取成功",
		"res":     "success",
	})
}

func GetSSLStatus(c *gin.Context) {
	domainUID := c.PostForm("domainUID")
	site := getSiteInfoByUID(domainUID)

	// 加载 AWS 配置
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(AwsAK, AwsSK, "")),
	)
	if err != nil {
		log.Printf("Unable to load SDK config, %v", err)
	}

	// 创建 ACM 客户端
	client := acm.NewFromConfig(cfg)

	// 替换为您的证书 ARN
	certificateArn := site.AwsSSLArn

	// 调用 DescribeCertificate 获取证书详情
	describeCertInput := &acm.DescribeCertificateInput{
		CertificateArn: aws.String(certificateArn),
	}

	certDetail, err := client.DescribeCertificate(context.TODO(), describeCertInput)
	if err != nil {
		log.Printf("failed to describe certificate, %v", err)
	}

	// 输出每个域名的 CNAME 验证状态
	for _, option := range certDetail.Certificate.DomainValidationOptions {
		if option.ValidationStatus == "SUCCESS" {
			log.Println("CNAME validation successful for domain:", *option.DomainName)
			site.Process = "SSL cname 验证成功"
			site.Status = "SSL cname verify down"
			updateSiteInfo(site)

			c.JSON(http.StatusOK, gin.H{
				"code":    20000,
				"message": "SSL Cname验证成功",
				"res":     "success",
			})
		} else {
			log.Println("CNAME validation pending for domain:", *option.DomainName)
			c.JSON(http.StatusOK, gin.H{
				"code":    20000,
				"message": "SSL Cname 未生效，请等待",
				"res":     "failed",
			})
		}
	}

}
