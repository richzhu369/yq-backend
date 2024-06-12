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
		site.Status = "Failed"
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
		site.Status = "Failed"
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
	site := getMerchantByName(merchantName)

	// 加载 AWS 配置
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("sa-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(publicProperty.AwsAK, publicProperty.AwsSK, "")),
	)
	if err != nil {
		log.Printf("Unable to load SDK config, %v", err)
		site.Status = "Failed"
		updateMerchantInfo(site)
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
		log.Printf("Failed to describe certificate, %v\n", err)
		site.Status = "Failed"
		updateMerchantInfo(site)
		upgradeProgress(8, merchantName, "el-icon-close", "primary")
		return false
	}

	site.Process = "等待Cname验证"

	// 输出 CNAME 验证信息
	for _, option := range certDetail.Certificate.DomainValidationOptions {
		log.Printf("DomainName: %sn", *option.DomainName)
		if option.ResourceRecord != nil {
			site.CnameKey = *option.ResourceRecord.Name
			site.CnameValue = *option.ResourceRecord.Value
			updateMerchantInfo(site)
		}

	}

	upgradeProgress(8, merchantName, "el-icon-check", "primary")
	upgradeProgress(9, merchantName, "el-icon-loading", "primary")
	return true
}

//func GetSSLStatus(merchantName string) bool {
//	site := getMerchantByName(merchantName)
//
//	// 加载 AWS 配置
//	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"),
//		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(publicProperty.AwsAK, publicProperty.AwsSK, "")),
//	)
//	if err != nil {
//		log.Printf("Unable to load SDK config, %v", err)
//	}
//
//	// 创建 ACM 客户端
//	client := acm.NewFromConfig(cfg)
//
//	// 替换为您的证书 ARN
//	certificateArn := site.AwsSSLArn
//
//	// 调用 DescribeCertificate 获取证书详情
//	describeCertInput := &acm.DescribeCertificateInput{
//		CertificateArn: aws.String(certificateArn),
//	}
//
//	certDetail, err := client.DescribeCertificate(context.TODO(), describeCertInput)
//	if err != nil {
//		log.Printf("failed to describe certificate, %v", err)
//	}
//
//	//fmt.Println(certDetail.Certificate.DomainValidationOptions[2])
//	//fmt.Println(certDetail.Certificate.DomainValidationOptions[5])
//	// 输出每个域名的 CNAME 验证状态
//	for _, option := range certDetail.Certificate.DomainValidationOptions {
//		if option.ValidationStatus == "SUCCESS" {
//			log.Println("CNAME validation successful for domain:", *option.DomainName)
//			site.Process = "SSL cname verify down"
//			updateMerchantInfo(site)
//
//			upgradeProgress(10, merchantName, "el-icon-check", "primary")
//			upgradeProgress(11, merchantName, "el-icon-loading", "primary")
//		} else {
//			log.Println("CNAME validation pending for domain:", *option.DomainName)
//			upgradeProgress(10, merchantName, "el-icon-close", "primary")
//		}
//	}
//
//	//todo : 这里要改，每检测一次，就返回一次结果，上面的for不对
//	return true
//}

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
		upgradeProgress(10, merchantName, "el-icon-close", "primary")
		site.Status = "failed"
		updateMerchantInfo(site)
		return false
	}
}
