package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/acm"
	"github.com/aws/aws-sdk-go-v2/service/acm/types"
	"log"
)

func createSSL(merchantName string) bool {

	site := getMerchantByName(merchantName)

	// 加载 AWS 配置
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(publicProperty.AwsAK, publicProperty.AwsSK, "")),
	)
	if err != nil {
		log.Printf("Unable to load SDK config, %v", err)
		site.Status = "failed"
		updateMerchantInfo(site)
		return false
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
		site.Status = "failed"
		updateMerchantInfo(site)
		return false
	}

	log.Printf("Certificate ARN: %s\n", aws.ToString(result.CertificateArn))

	site.Process = "Pending verification by cname on aws"
	site.AwsSSLArn = aws.ToString(result.CertificateArn)
	updateMerchantInfo(site)
	upgradeProgress(7, merchantName, "el-icon-check", "primary")
	upgradeProgress(8, merchantName, "el-icon-loading", "primary")
	return true
}

func GetSSLVerifyInfo(merchantName string) bool {
	//merchantName := c.PostForm("merchantName")

	site := getMerchantByName(merchantName)

	log.Println("传入商户名：", merchantName)
	log.Println("当前商户名：", site.MerchantName)
	// 加载 AWS 配置
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(publicProperty.AwsAK, publicProperty.AwsSK, "")),
	)
	if err != nil {
		log.Printf("Unable to load SDK config, %v", err)
		site.Status = "failed"
		updateMerchantInfo(site)
		return false
	}

	// 创建 ACM 客户端
	client := acm.NewFromConfig(cfg)
	// 指定查询的arn
	log.Println("site.AwsSSLArn")
	log.Println(site.AwsSSLArn)
	certificateArn := site.AwsSSLArn

	// 调用 DescribeCertificate 获取证书详情
	describeCertInput := &acm.DescribeCertificateInput{
		CertificateArn: aws.String(certificateArn),
	}

	log.Println("1")
	certDetail, err := client.DescribeCertificate(context.TODO(), describeCertInput)
	if err != nil {
		log.Printf("Failed to describe certificate, %v\n", err)
		site.Status = "failed"
		updateMerchantInfo(site)
		upgradeProgress(8, merchantName, "el-icon-close", "primary")
		return false
	}
	log.Println("2")

	site.Process = "等待Cname验证"

	if len(certDetail.Certificate.DomainValidationOptions) > 0 {

		option:= certDetail.Certificate.DomainValidationOptions[0]
		if option.ResourceRecord !=nil{
			log.Println("3")
			if option.ResourceRecord.Name !=nil{
				log.Println(*certDetail.Certificate.DomainValidationOptions[0].ResourceRecord.Name)
				site.CnameKey = *certDetail.Certificate.DomainValidationOptions[0].ResourceRecord.Name
				log.Println("3.1")
			}else {
				return false
			}
			if option.ResourceRecord.Value !=nil{
				log.Println(*certDetail.Certificate.DomainValidationOptions[0].ResourceRecord.Value)
				site.CnameValue = *certDetail.Certificate.DomainValidationOptions[0].ResourceRecord.Value
				log.Println("3.2")
			}else {
				return false
			}
		}else {
			return false
		}

		log.Println("4")
		updateMerchantInfo(site)
		log.Println("5")

		upgradeProgress(8, merchantName, "el-icon-check", "primary")
		upgradeProgress(9, merchantName, "el-icon-loading", "primary")
		log.Println("获取SSL证书验证的Cname和Value成功")
	} else {
		upgradeProgress(8, merchantName, "el-icon-close", "primary")
		log.Println("获取SSL证书验证的Cname和Value失败")
	}

	// 输出 CNAME 验证信息
	//for _, option := range certDetail.Certificate.DomainValidationOptions {
	//	log.Println("option的值：", *option.ResourceRecord.Name, *option.ResourceRecord.Value)
	//	if option.ResourceRecord != nil {
	//		log.Println("记录的商户cname为：", "CNAME:", *option.ResourceRecord.Name, "Value:", *option.ResourceRecord.Value, site.MerchantName)
	//		site.CnameKey = *option.ResourceRecord.Name
	//		site.CnameValue = *option.ResourceRecord.Value
	//
	//		updateMerchantInfo(site)
	//		log.Println("更新商户的Cname和Value: ", site.MerchantName, site.CnameKey, site.CnameValue)
	//	} else {
	//		log.Println("option.ResourceRecord 验证记录为空")
	//	}
	//}

	return true

	//c.JSON(http.StatusOK, gin.H{"message": "success"})
}

func GetSSLStatus(merchantName string) bool {
	site := getMerchantByName(merchantName)

	// 加载 AWS 配置
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(publicProperty.AwsAK, publicProperty.AwsSK, "")),
	)
	if err != nil {
		log.Printf("Unable to load SDK config, %v", err)
		return false
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
		return false
	}

	if certDetail.Certificate.DomainValidationOptions[0].ValidationStatus == "SUCCESS" {
		log.Println("CNAME validation successful for domain:", *certDetail.Certificate.DomainValidationOptions[0].DomainName)
		site.Process = "SSL cname verify down"
		updateMerchantInfo(site)
		upgradeProgress(10, merchantName, "el-icon-check", "primary")
		upgradeProgress(11, merchantName, "el-icon-loading", "primary")
		return true
	} else {
		log.Println("CNAME validation pending for domain:", *certDetail.Certificate.DomainValidationOptions[0].DomainName)
		site.Status = "failed"
		updateMerchantInfo(site)
		return false
	}
}
