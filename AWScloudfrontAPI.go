package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
	"github.com/gin-gonic/gin"
	"log"
)

func createCloudFront(c *gin.Context) {
	domainUID := c.PostForm("domainUID")
	site := getSiteInfoByUID(domainUID)

	// 加载 AWS 配置
	cfg, err := config.LoadDefaultConfig(context.TODO(),config.WithRegion("sa-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("AKIAUIUNFZ2WZZIQ6TNC", "1m41PBiFkxh5UloUY06sVvkQztUtnH+VgHyEcYMi", "")),
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
				TargetOriginId:       aws.String("myS3Origin"),
				ViewerProtocolPolicy: "redirect-to-https",
				ForwardedValues: &types.ForwardedValues{
					QueryString: aws.Bool(false),
					Cookies: &types.CookiePreference{
						Forward: "none",
					},
				},
				MinTTL: aws.Int64(0),
			},
			Origins: &types.Origins{
				Quantity: aws.Int32(1),
				Items: []types.Origin{
					{
						Id:         aws.String("myS3Origin"),
						DomainName: aws.String("af400c4b64edf4620810c92cd5dd5d82-499673083.sa-east-1.elb.amazonaws.com"),
						CustomOriginConfig: &types.CustomOriginConfig{
							HTTPPort:             aws.Int32(80),
							HTTPSPort:            aws.Int32(443),
							OriginProtocolPolicy: "http-only",
						},
					},
				},
			},
			Comment: aws.String(""),
			Enabled: aws.Bool(true),
		},
	}

	// 创建 CloudFront 分发
	result, err := client.CreateDistribution(context.TODO(), input)
	if err != nil {
		log.Fatalf("Failed to create distribution, %v", err)
	}

	fmt.Printf("CloudFront Distribution Created: %sn", aws.ToString(result.Distribution.Id))
}


