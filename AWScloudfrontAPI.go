package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

func createCloudFront(c *gin.Context) {
	domainUID := c.PostForm("domainUID")
	site := getSiteInfoByUID(domainUID)

	// 加载 AWS 配置
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("sa-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(AwsAK, AwsSK, "")),
	)
	if err != nil {
		log.Fatalf("Unable to load SDK config, %v", err)
	}

	// 创建 CloudFront 客户端
	client := cloudfront.NewFromConfig(cfg)

	// 创建 CloudFront 分发配置
	input := &cloudfront.CreateDistributionInput{
		DistributionConfig: &types.DistributionConfig{
			CallerReference: aws.String("unique-caller-reference"),
			Aliases: &types.Aliases{
				Quantity: aws.Int32(1),
				Items:    []string{site.AwsCdnDomain},
			},
			DefaultCacheBehavior: &types.DefaultCacheBehavior{
				TargetOriginId:        aws.String("af400c4b64edf4620810c92cd5dd5d82-499673083.sa-east-1.elb.amazonaws.com"),
				ViewerProtocolPolicy:  "redirect-to-https",
				CachePolicyId:         aws.String("4135ea2d-6df8-44a3-9df3-4b5a84be39ad"),
				OriginRequestPolicyId: aws.String("216adef6-5c7f-47e4-b989-5492eafa07d3"),
				AllowedMethods: &types.AllowedMethods{
					Items: []types.Method{
						types.MethodGet,
						types.MethodHead,
						types.MethodOptions,
						types.MethodPut,
						types.MethodPost,
						types.MethodPatch,
						types.MethodDelete,
					},
					Quantity: aws.Int32(7),
					CachedMethods: &types.CachedMethods{
						Items: []types.Method{
							types.MethodGet,
							types.MethodHead,
						},
						Quantity: aws.Int32(2),
					},
				},
			},
			Origins: &types.Origins{
				Quantity: aws.Int32(1),
				Items: []types.Origin{
					{
						Id:         aws.String("af400c4b64edf4620810c92cd5dd5d82-499673083.sa-east-1.elb.amazonaws.com"),
						DomainName: aws.String("af400c4b64edf4620810c92cd5dd5d82-499673083.sa-east-1.elb.amazonaws.com"),
						CustomOriginConfig: &types.CustomOriginConfig{
							HTTPPort:             aws.Int32(80),
							HTTPSPort:            aws.Int32(443),
							OriginProtocolPolicy: "http-only",
						},
					},
				},
			},
			ViewerCertificate: &types.ViewerCertificate{
				ACMCertificateArn:      aws.String(site.AwsSSLArn),
				SSLSupportMethod:       "sni-only",
				MinimumProtocolVersion: "TLSv1.2_2021",
			},
			Comment: aws.String("由yq-devops平台创建"),
			Enabled: aws.Bool(true),
		},
	}

	// 创建 CloudFront 分发
	result, err := client.CreateDistribution(context.TODO(), input)
	if err != nil {
		log.Printf("Failed to create distribution, %v", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    20000,
			"message": "创建cloudfront失败",
			"res":     "failed",
		})
		return
	}

	log.Printf("CloudFront Distribution Created: %s\n", aws.ToString(result.Distribution.Id))
	site.Process = "创建cloudfront 成功"
	site.Status = "create cloudfront success"
	site.CloudFrontID = *result.Distribution.Id
	updateSiteInfo(site)

	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"message": "创建cloudfront成功",
		"res":     "success",
	})
}

func GetCloudFrontDomain(c *gin.Context) {
	domainUID := c.PostForm("domainUID")
	site := getSiteInfoByUID(domainUID)

	// 加载 AWS 配置
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("sa-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(AwsAK, AwsSK, "")),
	)
	if err != nil {
		log.Printf("Unable to load SDK config, %v\n", err)
	}

	// 创建 CloudFront 客户端
	client := cloudfront.NewFromConfig(cfg)

	// CloudFront 分发的 ID
	distributionID := site.CloudFrontID

	// 调用 GetDistribution 获取分发详情
	getDistInput := &cloudfront.GetDistributionInput{
		Id: aws.String(distributionID),
	}

	resp, err := client.GetDistribution(context.TODO(), getDistInput)
	if err != nil {
		log.Printf("Failed to get distribution, %v\n", err)

		c.JSON(http.StatusOK, gin.H{
			"code":    20000,
			"message": "获取cloudfront失败",
		})
		return
	}

	// 输出 CloudFront 分发的域名
	log.Printf("CloudFront Distribution Domain Name: %s\n", *resp.Distribution.DomainName)
	site.CloudfrontRecord = *resp.Distribution.DomainName
	site.Process = "获取cloudfront 域名成功"
	site.Status = "get cloudfront domain success"
	updateSiteInfo(site)
	c.JSON(http.StatusOK, gin.H{
		"code":    20000,
		"message": "获取cloudfront域名成功",
	})
}
